package comment_handler

import (
	"readmodels/internal/comment"
	common_data "readmodels/internal/common/data"

	"github.com/rs/zerolog/log"
)

type CommentWasDeletedEvent struct {
	CommentId uint64 `json:"commentId"`
}

type CommentWasDeletedEventService interface {
	DeleteComment(commentId uint64)
}

type CommentWasDeletedEventHandler struct {
	service CommentWasDeletedEventService
}

func NewCommentWasDeletedEventHandler(repository comment.Repository) *CommentWasDeletedEventHandler {
	return &CommentWasDeletedEventHandler{
		service: comment.NewCommentService(repository),
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

	handler.service.DeleteComment(commentWasDeletedEvent.CommentId)
}
