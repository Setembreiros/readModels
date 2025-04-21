package comment_handler

import (
	common_data "readmodels/internal/common/data"

	"github.com/rs/zerolog/log"
)

//go:generate mockgen -source=comment_was_deleted_event_handler.go -destination=test/mock/comment_was_deleted_event_handler.go

type CommentWasDeletedEvent struct {
	PostId    string `json:"postId"`
	CommentId uint64 `json:"commentId"`
}

type CommentWasDeletedEventService interface {
	DeleteComment(postId string, commentId uint64)
}

type CommentWasDeletedEventHandler struct {
	service CommentWasDeletedEventService
}

func NewCommentWasDeletedEventHandler(service CommentWasDeletedEventService) *CommentWasDeletedEventHandler {
	return &CommentWasDeletedEventHandler{
		service: service,
	}
}

func (handler *CommentWasDeletedEventHandler) Handle(event []byte) {
	var commentWasDeletedEvent CommentWasDeletedEvent
	log.Info().Msg("Handling CommentWasDeletedEvent")

	err := common_data.DeserializeData(event, &commentWasDeletedEvent)
	if err != nil {
		log.Error().Stack().Err(err).Msg("Invalid event data")
		return
	}

	handler.service.DeleteComment(commentWasDeletedEvent.PostId, commentWasDeletedEvent.CommentId)
}
