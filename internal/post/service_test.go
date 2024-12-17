package post_test

import (
	"bytes"
	"errors"
	"fmt"
	"readmodels/internal/post"
	mock_post "readmodels/internal/post/mock"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

var serviceLoggerOutput bytes.Buffer
var serviceRepository *mock_post.MockRepository
var postService *post.PostService

func setUpService(t *testing.T) {
	ctrl := gomock.NewController(t)
	serviceRepository = mock_post.NewMockRepository(ctrl)
	log.Logger = log.Output(&serviceLoggerOutput)
	postService = post.NewPostService(serviceRepository)
}

func TestCreateNewPostMetadataWithService(t *testing.T) {
	setUpService(t)
	timeNow := time.Now().UTC()
	data := &post.PostMetadata{
		PostId:      "123456",
		Username:    "user123",
		Type:        "TEXT",
		Title:       "Exemplo de Título",
		Description: "Exemplo de Descrição",
		CreatedAt:   timeNow,
		LastUpdated: timeNow,
	}
	serviceRepository.EXPECT().AddNewPostMetadata(data)

	postService.CreateNewPostMetadata(data)

	assert.Contains(t, serviceLoggerOutput.String(), "Post metadata for id 123456 was added")
}

func TestErrorOnCreateNewPostMetadataWithService(t *testing.T) {
	setUpService(t)
	timeNow := time.Now().UTC()
	data := &post.PostMetadata{
		PostId:      "123456",
		Username:    "user123",
		Type:        "TEXT",
		Title:       "Exemplo de Título",
		Description: "Exemplo de Descrição",
		CreatedAt:   timeNow,
		LastUpdated: timeNow,
	}
	serviceRepository.EXPECT().AddNewPostMetadata(data).Return(errors.New("some error"))

	postService.CreateNewPostMetadata(data)

	assert.Contains(t, serviceLoggerOutput.String(), "Error adding post metadata for id 123456")
}

func TestGetPostMetadatasByUserWithService(t *testing.T) {
	setUpService(t)
	username := "username1"
	lastPostId := "post4"
	lastPostCreatedAt := "0001-01-03T00:00:00Z"
	limit := 3
	timeNow := time.Now().UTC()
	expectedData := []*post.PostMetadata{
		{
			PostId:      "123456",
			Username:    username,
			Type:        "TEXT",
			Title:       "Exemplo de Título",
			Description: "Exemplo de Descrição",
			CreatedAt:   timeNow,
			LastUpdated: timeNow,
		},
		{
			PostId:      "abcdef",
			Username:    username,
			Type:        "IMAGE",
			Title:       "Exemplo de Título 2",
			Description: "Exemplo de Descrição 2",
			CreatedAt:   timeNow,
			LastUpdated: timeNow,
		},
	}
	serviceRepository.EXPECT().GetPostMetadatasByUser(username, lastPostId, lastPostCreatedAt, limit).Return(expectedData, "post7", "0001-01-06T00:00:00Z", nil)

	postService.GetPostMetadatasByUser(username, lastPostId, lastPostCreatedAt, limit)
}

func TestErrorOnGetPostMetadatasByUserWithService(t *testing.T) {
	setUpService(t)
	username := "username1"
	lastPostId := "post4"
	lastPostCreatedAt := "0001-01-03T00:00:00Z"
	limit := 2
	expectedData := []*post.PostMetadata{}
	serviceRepository.EXPECT().GetPostMetadatasByUser(username, lastPostId, lastPostCreatedAt, limit).Return(expectedData, "", "", errors.New("some error"))

	postService.GetPostMetadatasByUser(username, lastPostId, lastPostCreatedAt, limit)

	assert.Contains(t, serviceLoggerOutput.String(), "Error getting post metadatas for username "+username)
}

func TestRemovePostMetadataWithService(t *testing.T) {
	setUpService(t)
	postIds := []string{"123456", "abcdef", "1a2b3e"}
	serviceRepository.EXPECT().RemovePostMetadata(postIds)

	postService.RemovePostMetadata(postIds)

	assert.Contains(t, serviceLoggerOutput.String(), fmt.Sprintf("Post metadatas for ids %v were removed", postIds))
}

func TestRemovePostMetadataWithService_Error(t *testing.T) {
	setUpService(t)
	postIds := []string{"123456", "abcdef", "1a2b3e"}
	serviceRepository.EXPECT().RemovePostMetadata(postIds).Return(errors.New("some error"))

	postService.RemovePostMetadata(postIds)

	assert.Contains(t, serviceLoggerOutput.String(), fmt.Sprintf("Error removing post metadatas for id %v", postIds))
}
