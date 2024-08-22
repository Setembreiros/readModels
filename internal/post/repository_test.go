package post_test

import (
	database "readmodels/internal/db"
	mock_database "readmodels/internal/db/mock"
	"readmodels/internal/post"
	"testing"
	"time"

	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
)

var client *mock_database.MockDatabaseClient
var postRepository post.PostRepository

func setUp(t *testing.T) {
	ctrl := gomock.NewController(t)
	client = mock_database.NewMockDatabaseClient(ctrl)
	postRepository = post.PostRepository(*database.NewDatabase(client))
}

func TestAddNewPostMetadataInRepository(t *testing.T) {
	setUp(t)
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
	client.EXPECT().InsertData("PostMetadata", data)

	postRepository.AddNewPostMetadata(data)
}

func TestGetPostMetadatasByUserInRepository(t *testing.T) {
	setUp(t)
	username := "username1"
	timeNow := time.Now().UTC()
	data := []*database.PostMetadata{
		{
			PostId:      "123456",
			Username:    username,
			Type:        "TEXT",
			FileType:    "txt",
			Title:       "Exemplo de Título",
			Description: "Exemplo de Descrição",
			CreatedAt:   timeNow,
			LastUpdated: timeNow,
		},
		{
			PostId:      "abcdef",
			Username:    username,
			Type:        "IMAGE",
			FileType:    "png",
			Title:       "Exemplo de Título 2",
			Description: "Exemplo de Descrição 2",
			CreatedAt:   timeNow,
			LastUpdated: timeNow,
		},
	}
	expectedResult := []*post.PostMetadata{
		{
			PostId:      "123456",
			Username:    username,
			Type:        "TEXT",
			FileType:    "txt",
			Title:       "Exemplo de Título",
			Description: "Exemplo de Descrição",
			CreatedAt:   timeNow,
			LastUpdated: timeNow,
		},
		{
			PostId:      "abcdef",
			Username:    username,
			Type:        "IMAGE",
			FileType:    "png",
			Title:       "Exemplo de Título 2",
			Description: "Exemplo de Descrição 2",
			CreatedAt:   timeNow,
			LastUpdated: timeNow,
		},
	}
	client.EXPECT().GetPostsByIndexUser(username).Return(data, nil)

	result, _ := postRepository.GetPostMetadatasByUser(username)

	assert.Equal(t, expectedResult, result)
}

func TestRemovePostMetadataInRepository(t *testing.T) {
	setUp(t)
	postIds := []string{"123456", "abcdef", "1a2b3e"}
	expectedSlice := []any{"123456", "abcdef", "1a2b3e"}
	client.EXPECT().RemoveMultipleData("Posts", expectedSlice)

	postRepository.RemovePostMetadata(postIds)
}
