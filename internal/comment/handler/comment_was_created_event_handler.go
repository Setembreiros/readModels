package comment_handler

import (
	"readmodels/internal/comment"
	common_data "readmodels/internal/common/data"
	"readmodels/internal/model"
	"time"

	"github.com/rs/zerolog/log"
)

type CommentWasCreatedEvent struct {
	CommentId uint64 `json:"commentId"`
	Username  string `json:"username"`
	PostId    string `json:"postId"`
	Content   string `json:"content"`
	CreatedAt string `json:"createdAt"`
}

type CommentWasCreatedEventService interface {
	CreateNewComment(data *model.Comment)
}

type CommentWasCreatedEventHandler struct {
	service CommentWasCreatedEventService
}

func NewCommentWasCreatedEventHandler(repository comment.Repository) *CommentWasCreatedEventHandler {
	return &CommentWasCreatedEventHandler{
		service: comment.NewCommentService(repository),
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

	handler.service.CreateNewComment(data)
}

func mapData(event CommentWasCreatedEvent) (*model.Comment, error) {
	timeLayout := "2006-01-02T15:04:05.000000000Z"
	parsedCreatedAt, err := time.Parse(timeLayout, event.CreatedAt)
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
