package post_test

import (
	database "readmodels/internal/db"
	mock_database "readmodels/internal/db/mock"
	"readmodels/internal/post"
	"testing"
	"time"

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
