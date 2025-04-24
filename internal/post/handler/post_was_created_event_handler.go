package post_handler

import (
	common_data "readmodels/internal/common/data"
	"readmodels/internal/model"
	"readmodels/internal/post"
	"time"

	"github.com/rs/zerolog/log"
)

//go:generate mockgen -source=post_was_created_event_handler.go -destination=mock/post_was_created_event_handler.go

type Metadata struct {
	Username    string `json:"username"`
	Type        string `json:"type"`
	Title       string `json:"title"`
	Description string `json:"description"`
	CreatedAt   string `json:"createdAt"`
	LastUpdated string `json:"lastUpdated"`
}

type PostWasCreatedEvent struct {
	PostId   string   `json:"post_id"`
	Metadata Metadata `json:"metadata"`
}

type PostWasCreatedEventService interface {
	CreateNewPostMetadata(data *post.PostMetadata)
}

type PostWasCreatedEventHandler struct {
	service PostWasCreatedEventService
}

func NewPostWasCreatedEventHandler(service PostWasCreatedEventService) *PostWasCreatedEventHandler {
	return &PostWasCreatedEventHandler{
		service: service,
	}
}

func (handler *PostWasCreatedEventHandler) Handle(event []byte) {
	var postWasCreatedEvent PostWasCreatedEvent
	log.Info().Msg("Handling PostWasCreatedEvent")

	err := common_data.DeserializeData(event, &postWasCreatedEvent)
	if err != nil {
		log.Error().Stack().Err(err).Msg("Invalid event data")
		return
	}

	data, err := mapData(postWasCreatedEvent)
	if err != nil {
		return
	}

	handler.service.CreateNewPostMetadata(data)
}

func mapData(event PostWasCreatedEvent) (*post.PostMetadata, error) {
	parsedCreatedAt, err := time.Parse(model.TimeLayout, event.Metadata.CreatedAt)
	if err != nil {
		log.Error().Stack().Err(err).Msg("Error parsing time CreatedAt")
		return nil, err
	}
	parsedLastUpdatedAt, err := time.Parse(model.TimeLayout, event.Metadata.LastUpdated)

	if err != nil {
		log.Error().Stack().Err(err).Msg("Error parsing time LastUpdated")
		return nil, err
	}
	return &post.PostMetadata{
		PostId:      event.PostId,
		Username:    event.Metadata.Username,
		Type:        event.Metadata.Type,
		Title:       event.Metadata.Title,
		Description: event.Metadata.Description,
		CreatedAt:   parsedCreatedAt,
		LastUpdated: parsedLastUpdatedAt,
	}, nil
}
