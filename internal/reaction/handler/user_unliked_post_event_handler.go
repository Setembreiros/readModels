package reaction_handler

import (
	common_data "readmodels/internal/common/data"
	"readmodels/internal/model"

	"github.com/rs/zerolog/log"
)

//go:generate mockgen -source=user_unliked_post_event_handler.go -destination=test/mock/user_unliked_post_event_handler.go

type UserUnlikedPostEvent struct {
	Username string `json:"username"`
	PostId   string `json:"postId"`
}

type UserUnlikedPostEventService interface {
	DeleteLikePost(data *model.LikePost)
}

type UserUnlikedPostEventHandler struct {
	service UserUnlikedPostEventService
}

func NewUserUnlikedPostEventHandler(service UserUnlikedPostEventService) *UserUnlikedPostEventHandler {
	return &UserUnlikedPostEventHandler{
		service: service,
	}
}

func (handler *UserUnlikedPostEventHandler) Handle(event []byte) {
	var userUnlikedPostEvent UserUnlikedPostEvent
	log.Info().Msg("Handling UserUnlikedPostEvent")

	err := common_data.DeserializeData(event, &userUnlikedPostEvent)
	if err != nil {
		log.Error().Stack().Err(err).Msg("Invalid event data")
		return
	}

	data, err := mapUserUnlikedPostEvent(userUnlikedPostEvent)
	if err != nil {
		return
	}

	handler.service.DeleteLikePost(data)
}

func mapUserUnlikedPostEvent(event UserUnlikedPostEvent) (*model.LikePost, error) {
	return &model.LikePost{
		Username: event.Username,
		PostId:   event.PostId,
	}, nil
}
