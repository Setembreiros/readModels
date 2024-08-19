package post_handler_test

import (
	"bytes"
	"encoding/json"
	"readmodels/internal/post"
	post_handler "readmodels/internal/post/handler"
	mock_post "readmodels/internal/post/mock"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

var loggerOutput bytes.Buffer
var repository *mock_post.MockRepository
var handler *post_handler.PostWasCreatedEventHandler

func setUpHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	repository = mock_post.NewMockRepository(ctrl)
	log.Logger = log.Output(&loggerOutput)
	handler = post_handler.NewPostWasCreatedEventHandler(repository)
}

func TestHandlePostWasCreatedEvent(t *testing.T) {
	setUpHandler(t)
	timeNow := time.Now().UTC()
	data := &post_handler.PostWasCreatedEvent{
		PostId: "123456",
		Metadata: post_handler.Metadata{
			Username:    "user123",
			Type:        "TEXT",
			FileType:    "txt",
			Title:       "Exemplo de Título",
			Description: "Exemplo de Descrição",
			CreatedAt:   timeNow,
			LastUpdated: timeNow,
		},
	}
	event, _ := json.Marshal(data)
	expectedPostMetadata := &post.PostMetadata{
		PostId:      "123456",
		Username:    "user123",
		Type:        "TEXT",
		FileType:    "txt",
		Title:       "Exemplo de Título",
		Description: "Exemplo de Descrição",
		CreatedAt:   timeNow,
		LastUpdated: timeNow,
	}
	repository.EXPECT().AddNewPostMetadata(expectedPostMetadata)

	handler.Handle(event)
}

func TestInvalidDataInPostWasCreatedEventHandler(t *testing.T) {
	setUpHandler(t)
	invalidData := "invalid data"
	event, _ := json.Marshal(invalidData)

	handler.Handle(event)

	assert.Contains(t, loggerOutput.String(), "Invalid event data")
}
