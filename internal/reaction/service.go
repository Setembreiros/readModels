package reaction

import (
	"readmodels/internal/model"

	"github.com/rs/zerolog/log"
)

//go:generate mockgen -source=service.go -destination=test/mock/service.go

type Repository interface {
	CreatePostLike(data *model.PostLike) error
	CreatePostSuperlike(data *model.PostSuperlike) error
	CreateReview(data *model.Review) error
	GetPostLikesMetadata(postId, lastUsername string, limit int) ([]*model.UserMetadata, string, error)
	GetPostSuperlikesMetadata(postId, lastUsername string, limit int) ([]*model.UserMetadata, string, error)
	DeletePostLike(data *model.PostLike) error
	DeletePostSuperlike(data *model.PostSuperlike) error
}

type ReactionService struct {
	repository Repository
}

func NewReactionService(repository Repository) *ReactionService {
	return &ReactionService{
		repository: repository,
	}
}

func (s *ReactionService) CreatePostLike(data *model.PostLike) {
	err := s.repository.CreatePostLike(data)
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Error creating postLike, username: %s -> postId: %s", data.User.Username, data.PostId)
		return
	}

	log.Info().Msgf("PostLike was created, username: %s -> postId: %s", data.User.Username, data.PostId)
}

func (s *ReactionService) CreatePostSuperlike(data *model.PostSuperlike) {
	err := s.repository.CreatePostSuperlike(data)
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Error creating postSuperlike, username: %s -> postId: %s", data.User.Username, data.PostId)
		return
	}

	log.Info().Msgf("PostSuperlike was created, username: %s -> postId: %s", data.User.Username, data.PostId)
}

func (s *ReactionService) CreateReview(data *model.Review) {
	err := s.repository.CreateReview(data)
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Error creating review with id %d in post %s", data.ReviewId, data.PostId)
		return
	}

	log.Info().Msgf("Review with id %d in post %s was created", data.ReviewId, data.PostId)
}

func (s *ReactionService) GetPostLikesMetadata(postId, lastUsername string, limit int) ([]*model.UserMetadata, string, error) {
	users, lastUsername, err := s.repository.GetPostLikesMetadata(postId, lastUsername, limit)
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Error getting post %s's likes", postId)
		return users, lastUsername, err
	}

	return users, lastUsername, nil
}

func (s *ReactionService) GetPostSuperlikesMetadata(postId, lastUsername string, limit int) ([]*model.UserMetadata, string, error) {
	users, lastUsername, err := s.repository.GetPostSuperlikesMetadata(postId, lastUsername, limit)
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Error getting post %s's superlikes", postId)
		return users, lastUsername, err
	}

	return users, lastUsername, nil
}

func (s *ReactionService) DeletePostLike(data *model.PostLike) {
	err := s.repository.DeletePostLike(data)
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Error deleting postLike, username: %s -> postId: %s", data.User.Username, data.PostId)
		return
	}

	log.Info().Msgf("PostLike was deleted, username: %s -> postId: %s", data.User.Username, data.PostId)
}

func (s *ReactionService) DeletePostSuperlike(data *model.PostSuperlike) {
	err := s.repository.DeletePostSuperlike(data)
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Error deleting postSuperlike, username: %s -> postId: %s", data.User.Username, data.PostId)
		return
	}

	log.Info().Msgf("PostSuperlike was deleted, username: %s -> postId: %s", data.User.Username, data.PostId)
}
