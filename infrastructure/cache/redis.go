package database

import (
	"context"
	"encoding/json"
	"fmt"
	"readmodels/internal/model"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type RedisCacheClient struct {
	ctx    context.Context
	client *redis.Client
}

type CommentsData struct {
	Comments      []*model.Comment `json:"comments"`
	LastCommentId uint64           `json:"lastCommentId"`
}

type PostLikesData struct {
	PostLikes    []*model.UserMetadata `json:"postLikes"`
	LastUsername string                `json:"lastUsername"`
}

type PostSuperlikesData struct {
	PostSuperlikes []*model.UserMetadata `json:"postSuperlikes"`
	LastUsername   string                `json:"lastUsername"`
}

func NewRedisClient(cacheUri, cachePassword string, ctx context.Context) *RedisCacheClient {
	redisConfig := &redis.Options{
		Addr:     cacheUri,
		Password: cachePassword,
		DB:       0, // Use default DB
		Protocol: 2, // Connection protocol
	}
	client := &RedisCacheClient{
		ctx:    ctx,
		client: redis.NewClient(redisConfig),
	}

	client.verifyConnection()

	return client
}

func (c *RedisCacheClient) verifyConnection() {
	err := c.client.Set(c.ctx, "foo", "bar", 10*time.Second).Err()
	if err != nil {
		log.Error().Stack().Err(err).Msg("Conection to Redis not stablished")
		panic(err)
	}

	_, err = c.client.Get(c.ctx, "foo").Result()
	if err != nil {
		log.Error().Stack().Err(err).Msg("Conection to Redis not stablished")
		panic(err)
	}
	log.Info().Msgf("Connection to Redis established.")
}

func (c *RedisCacheClient) Clean() {
	err := c.client.FlushDB(c.ctx).Err()
	if err != nil {
		log.Warn().Stack().Err(err).Msg("Failed to clean entire Redis cache")
		return
	}

	log.Info().Msg("Entire Redis cache cleaned successfully")
}

func (c *RedisCacheClient) SetPostComments(postId string, lastCommentId uint64, limit int, comments []*model.Comment) {
	cacheKey := generateCommentsCacheKey(postId, lastCommentId, limit)

	newLastCommentId := uint64(0)
	if len(comments) > 0 {
		newLastCommentId = comments[len(comments)-1].CommentId
	}

	data := CommentsData{
		Comments:      comments,
		LastCommentId: newLastCommentId,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Warn().Stack().Err(err).Msg("Failed to serialize comments data")
		return
	}

	err = c.client.Set(c.ctx, cacheKey, jsonData, 5*time.Minute).Err()
	if err != nil {
		log.Warn().Stack().Err(err).Msg("Failed to set comments in cache")
	}
}

func (c *RedisCacheClient) SetPostLikes(postId string, lastUsername string, limit int, postLikes []*model.UserMetadata) {
	cacheKey := generatePostLikesCacheKey(postId, lastUsername, limit)

	newLastUsername := ""
	if len(postLikes) > 0 {
		newLastUsername = postLikes[len(postLikes)-1].Username
	}

	data := PostLikesData{
		PostLikes:    postLikes,
		LastUsername: newLastUsername,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Warn().Stack().Err(err).Msg("Failed to serialize postLikes data")
		return
	}

	err = c.client.Set(c.ctx, cacheKey, jsonData, 5*time.Minute).Err()
	if err != nil {
		log.Warn().Stack().Err(err).Msg("Failed to set postLikes in cache")
	}
}

func (c *RedisCacheClient) SetPostSuperlikes(postId string, lastUsername string, limit int, postSuperlikes []*model.UserMetadata) {
	cacheKey := generatePostSuperlikesCacheKey(postId, lastUsername, limit)

	newLastUsername := ""
	if len(postSuperlikes) > 0 {
		newLastUsername = postSuperlikes[len(postSuperlikes)-1].Username
	}

	data := PostSuperlikesData{
		PostSuperlikes: postSuperlikes,
		LastUsername:   newLastUsername,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Warn().Stack().Err(err).Msg("Failed to serialize postSuperlikes data")
		return
	}

	err = c.client.Set(c.ctx, cacheKey, jsonData, 5*time.Minute).Err()
	if err != nil {
		log.Warn().Stack().Err(err).Msg("Failed to set postSuperlikes in cache")
	}
}

func (c *RedisCacheClient) GetPostComments(postId string, lastCommentId uint64, limit int) ([]*model.Comment, uint64, bool) {
	cacheKey := generateCommentsCacheKey(postId, lastCommentId, limit)

	jsonData, err := c.client.Get(c.ctx, cacheKey).Bytes()
	if err != nil {
		if err == redis.Nil {
			return []*model.Comment{}, uint64(0), false
		}

		log.Warn().Stack().Err(err).Msg("Failed to retrieve comments from cache")
		return []*model.Comment{}, uint64(0), false
	}

	var data CommentsData
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		log.Warn().Stack().Err(err).Msg("Failed to deserialize comments data")
		return []*model.Comment{}, uint64(0), false
	}

	log.Info().Msgf("Data retrieve from cache for key %s", cacheKey)

	return data.Comments, data.LastCommentId, true
}

func (c *RedisCacheClient) GetPostLikes(postId string, lastUsername string, limit int) ([]*model.UserMetadata, string, bool) {
	cacheKey := generatePostLikesCacheKey(postId, lastUsername, limit)

	jsonData, err := c.client.Get(c.ctx, cacheKey).Bytes()
	if err != nil {
		if err == redis.Nil {
			return []*model.UserMetadata{}, "", false
		}

		log.Warn().Stack().Err(err).Msg("Failed to retrieve postLikes from cache")
		return []*model.UserMetadata{}, "", false
	}

	var data PostLikesData
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		log.Warn().Stack().Err(err).Msg("Failed to deserialize postLikes data")
		return []*model.UserMetadata{}, "", false
	}

	log.Info().Msgf("Data retrieve from cache for key %s", cacheKey)

	return data.PostLikes, data.LastUsername, true
}

func (c *RedisCacheClient) GetPostSuperlikes(postId string, lastUsername string, limit int) ([]*model.UserMetadata, string, bool) {
	cacheKey := generatePostSuperlikesCacheKey(postId, lastUsername, limit)

	jsonData, err := c.client.Get(c.ctx, cacheKey).Bytes()
	if err != nil {
		if err == redis.Nil {
			return []*model.UserMetadata{}, "", false
		}

		log.Warn().Stack().Err(err).Msg("Failed to retrieve postSuperlikes from cache")
		return []*model.UserMetadata{}, "", false
	}

	var data PostSuperlikesData
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		log.Warn().Stack().Err(err).Msg("Failed to deserialize postSuperlikes data")
		return []*model.UserMetadata{}, "", false
	}

	log.Info().Msgf("Data retrieve from cache for key %s", cacheKey)

	return data.PostSuperlikes, data.LastUsername, true
}

func generateCommentsCacheKey(postId string, lastCommentId uint64, limit int) string {
	return fmt.Sprintf("comments:%s:%d:%d", postId, lastCommentId, limit)
}

func generatePostLikesCacheKey(postId string, lastUsername string, limit int) string {
	return fmt.Sprintf("postLikes:%s:%s:%d", postId, lastUsername, limit)
}

func generatePostSuperlikesCacheKey(postId string, lastUsername string, limit int) string {
	return fmt.Sprintf("postSuperlikes:%s:%s:%d", postId, lastUsername, limit)
}
