package post_test

import (
	database "readmodels/internal/db"
	mock_database "readmodels/internal/db/test/mock"
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
	lastPostId := "post4"
	lastPostCreatedAt := "0001-01-03T00:00:00Z"
	limit := 3
	timeNow := time.Now().UTC()
	data := []*database.PostMetadata{
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
	expectedResult := []*post.PostMetadata{
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
	expectedLastPostId := "post7"
	expectedLastPostCreatedAt := "0001-01-06T00:00:00Z"
	client.EXPECT().GetPostsByIndexUser(username, lastPostId, lastPostCreatedAt, limit).Return(data, expectedLastPostId, expectedLastPostCreatedAt, nil)

	result, lastPostId, lastPostCreatedAt, _ := postRepository.GetPostMetadatasByUser(username, lastPostId, lastPostCreatedAt, limit)

	assert.Equal(t, expectedResult, result)
	assert.Equal(t, expectedLastPostId, lastPostId)
	assert.Equal(t, expectedLastPostCreatedAt, lastPostCreatedAt)
}

func TestRemovePostMetadataInRepository(t *testing.T) {
	setUp(t)
	postIds := []string{"123456", "abcdef", "1a2b3e"}
	expectedKeys := []any{
		&database.PostMetadataKey{
			PostId: "123456",
		},
		&database.PostMetadataKey{
			PostId: "abcdef",
		},
		&database.PostMetadataKey{
			PostId: "1a2b3e",
		},
	}
	client.EXPECT().RemoveMultipleData("PostMetadata", expectedKeys)

	postRepository.RemovePostMetadata(postIds)
}
