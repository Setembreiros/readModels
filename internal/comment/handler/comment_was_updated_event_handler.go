package comment_handler

import (
	"readmodels/internal/comment"
	common_data "readmodels/internal/common/data"
	"readmodels/internal/model"
	"time"

	"github.com/rs/zerolog/log"
)

type CommentWasUpdatedEvent struct {
	CommentId uint64 `json:"commentId"`
	Content   string `json:"content"`
	UpdatedAt string `json:"updatedAt"`
}

type CommentWasUpdatedEventService interface {
	UpdateComment(data *model.Comment)
}

type CommentWasUpdatedEventHandler struct {
	service CommentWasUpdatedEventService
}

func NewCommentWasUpdatedEventHandler(repository comment.Repository) *CommentWasUpdatedEventHandler {
	return &CommentWasUpdatedEventHandler{
		service: comment.NewCommentService(repository),
	}
}

func (handler *CommentWasUpdatedEventHandler) Handle(event []byte) {
	var commentWasUpdatedEvent CommentWasUpdatedEvent
	log.Info().Msg("Handling CommentWasUpdatedEvent")

	err := common_data.DeserializeData(event, &commentWasUpdatedEvent)
	if err != nil {
		log.Error().Stack().Err(err).Msg("Invalid event data")
		return
	}

	data, err := mapUpdateEventData(commentWasUpdatedEvent)
	if err != nil {
		return
	}

	handler.service.UpdateComment(data)
}

func mapUpdateEventData(event CommentWasUpdatedEvent) (*model.Comment, error) {
	parsedUpdatedAt, err := time.Parse(model.TimeLayout, event.UpdatedAt)
	if err != nil {
		log.Error().Stack().Err(err).Msg("Error parsing time UpdatedAt")
		return nil, err
	}

	return &model.Comment{
		CommentId: event.CommentId,
		Content:   event.Content,
		UpdatedAt: parsedUpdatedAt,
	}, nil
}
