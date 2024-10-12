package post_handler

import (
	"encoding/json"
	"readmodels/internal/post"
	"time"

	"github.com/rs/zerolog/log"
)

type Metadata struct {
	Username    string    `json:"username"`
	Type        string    `json:"type"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	LastUpdated time.Time `json:"last_updated"`
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

func NewPostWasCreatedEventHandler(repository post.Repository) *PostWasCreatedEventHandler {
	return &PostWasCreatedEventHandler{
		service: post.NewPostService(repository),
	}
}

func (handler *PostWasCreatedEventHandler) Handle(event []byte) {
	var postWasCreatedEvent PostWasCreatedEvent
	log.Info().Msg("Handling PostWasCreatedEvent")

	err := Decode(event, &postWasCreatedEvent)
	if err != nil {
		log.Error().Stack().Err(err).Msg("Invalid event data")
		return
	}

	data := mapData(postWasCreatedEvent)
	handler.service.CreateNewPostMetadata(data)
}

func Decode(datab []byte, data any) error {
	return json.Unmarshal(datab, &data)
}

func mapData(event PostWasCreatedEvent) *post.PostMetadata {
	return &post.PostMetadata{
		PostId:      event.PostId,
		Username:    event.Metadata.Username,
		Type:        event.Metadata.Type,
		Title:       event.Metadata.Title,
		Description: event.Metadata.Description,
		CreatedAt:   event.Metadata.CreatedAt,
		LastUpdated: event.Metadata.LastUpdated,
	}
}
