package post_handler_test

import (
	"bytes"
	"encoding/json"
	"readmodels/internal/model"
	"readmodels/internal/post"
	post_handler "readmodels/internal/post/handler"
	mock_post "readmodels/internal/post/handler/mock"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

var loggerOutput bytes.Buffer
var service *mock_post.MockPostWasCreatedEventService
var handler *post_handler.PostWasCreatedEventHandler

func setUpHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	service = mock_post.NewMockPostWasCreatedEventService(ctrl)
	log.Logger = log.Output(&loggerOutput)
	handler = post_handler.NewPostWasCreatedEventHandler(service)
}

func TestHandlePostWasCreatedEvent(t *testing.T) {
	setUpHandler(t)
	timeNow := time.Now().UTC().Format(model.TimeLayout)
	data := &post_handler.PostWasCreatedEvent{
		PostId: "123456",
		Metadata: post_handler.Metadata{
			Username:    "user123",
			Type:        "TEXT",
			Title:       "Exemplo de Título",
			Description: "Exemplo de Descrição",
			CreatedAt:   timeNow,
			LastUpdated: timeNow,
		},
	}
	event, _ := json.Marshal(data)
	expectedTime, _ := time.Parse(model.TimeLayout, timeNow)
	expectedPostMetadata := &post.PostMetadata{
		PostId:      "123456",
		Username:    "user123",
		Type:        "TEXT",
		Title:       "Exemplo de Título",
		Description: "Exemplo de Descrição",
		CreatedAt:   expectedTime,
		LastUpdated: expectedTime,
	}
	service.EXPECT().CreateNewPostMetadata(expectedPostMetadata)

	handler.Handle(event)
}

func TestInvalidDataInPostWasCreatedEventHandler(t *testing.T) {
	setUpHandler(t)
	invalidData := "invalid data"
	event, _ := json.Marshal(invalidData)

	handler.Handle(event)

	assert.Contains(t, loggerOutput.String(), "Invalid event data")
}
