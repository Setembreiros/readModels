package reactio_handler

import (
	common_data "readmodels/internal/common/data"
	"readmodels/internal/model"

	"github.com/rs/zerolog/log"
)

//go:generate mockgen -source=user_liked_post_event_handler.go -destination=test/mock/user_liked_post_event_handler.go

type UserLikedPostEvent struct {
	Username string `json:"username"`
	PostId   string `json:"postId"`
}

type UserLikedPostEventService interface {
	CreateLikePost(data *model.LikePost)
}

type UserLikedPostEventHandler struct {
	service UserLikedPostEventService
}

func NewUserLikedPostEventHandler(service UserLikedPostEventService) *UserLikedPostEventHandler {
	return &UserLikedPostEventHandler{
		service: service,
	}
}

func (handler *UserLikedPostEventHandler) Handle(event []byte) {
	var userLikedPostEvent UserLikedPostEvent
	log.Info().Msg("Handling UserLikedPostEvent")

	err := common_data.DeserializeData(event, &userLikedPostEvent)
	if err != nil {
		log.Error().Stack().Err(err).Msg("Invalid event data")
		return
	}

	data, err := mapData(userLikedPostEvent)
	if err != nil {
		return
	}

	handler.service.CreateLikePost(data)
}

func mapData(event UserLikedPostEvent) (*model.LikePost, error) {
	return &model.LikePost{
		Username: event.Username,
		PostId:   event.PostId,
	}, nil
}
