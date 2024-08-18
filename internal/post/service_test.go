package post_test

import (
	"bytes"
	"errors"
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
		FileType:    "txt",
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
		FileType:    "txt",
		Title:       "Exemplo de Título",
		Description: "Exemplo de Descrição",
		CreatedAt:   timeNow,
		LastUpdated: timeNow,
	}
	serviceRepository.EXPECT().AddNewPostMetadata(data).Return(errors.New("some error"))

	postService.CreateNewPostMetadata(data)

	assert.Contains(t, serviceLoggerOutput.String(), "Error adding post metadata for id 123456")
}
