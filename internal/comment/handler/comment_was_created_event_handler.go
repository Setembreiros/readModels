package comment_handler

import (
	common_data "readmodels/internal/common/data"
	"readmodels/internal/model"
	"time"

	"github.com/rs/zerolog/log"
)

//go:generate mockgen -source=comment_was_created_event_handler.go -destination=test/mock/comment_was_created_event_handler.go

type CommentWasCreatedEvent struct {
	CommentId uint64 `json:"commentId"`
	Username  string `json:"username"`
	PostId    string `json:"postId"`
	Content   string `json:"content"`
	CreatedAt string `json:"createdAt"`
}

type CommentWasCreatedEventService interface {
	CreateComment(data *model.Comment)
}

type CommentWasCreatedEventHandler struct {
	service CommentWasCreatedEventService
}

func NewCommentWasCreatedEventHandler(service CommentWasCreatedEventService) *CommentWasCreatedEventHandler {
	return &CommentWasCreatedEventHandler{
		service: service,
	}
}

func (handler *CommentWasCreatedEventHandler) Handle(event []byte) {
	var commentWasCreatedEvent CommentWasCreatedEvent
	log.Info().Msg("Handling CommentWasCreatedEvent")

	err := common_data.DeserializeData(event, &commentWasCreatedEvent)
	if err != nil {
		log.Error().Stack().Err(err).Msg("Invalid event data")
		return
	}

	data, err := mapData(commentWasCreatedEvent)
	if err != nil {
		return
	}

	handler.service.CreateComment(data)
}

func mapData(event CommentWasCreatedEvent) (*model.Comment, error) {
	parsedCreatedAt, err := time.Parse(model.TimeLayout, event.CreatedAt)
	if err != nil {
		log.Error().Stack().Err(err).Msg("Error parsing time CreatedAt")
		return nil, err
	}

	return &model.Comment{
		CommentId: event.CommentId,
		Username:  event.Username,
		PostId:    event.PostId,
		Content:   event.Content,
		CreatedAt: parsedCreatedAt,
	}, nil
}
