package reaction_handler

import (
	common_data "readmodels/internal/common/data"
	"readmodels/internal/model"

	"github.com/rs/zerolog/log"
)

//go:generate mockgen -source=user_superliked_post_event_handler.go -destination=test/mock/user_superliked_post_event_handler.go

type UserSuperlikedPostEvent struct {
	Username string `json:"username"`
	PostId   string `json:"postId"`
}

type UserSuperlikedPostEventService interface {
	CreatePostSuperlike(data *model.PostSuperlike)
}

type UserSuperlikedPostEventHandler struct {
	service UserSuperlikedPostEventService
}

func NewUserSuperlikedPostEventHandler(service UserSuperlikedPostEventService) *UserSuperlikedPostEventHandler {
	return &UserSuperlikedPostEventHandler{
		service: service,
	}
}

func (handler *UserSuperlikedPostEventHandler) Handle(event []byte) {
	var userSuperlikedPostEvent UserSuperlikedPostEvent
	log.Info().Msg("Handling UserSuperlikedPostEvent")

	err := common_data.DeserializeData(event, &userSuperlikedPostEvent)
	if err != nil {
		log.Error().Stack().Err(err).Msg("Invalid event data")
		return
	}

	data, err := mapUserSuperlikedPostEvent(userSuperlikedPostEvent)
	if err != nil {
		return
	}

	handler.service.CreatePostSuperlike(data)
}

func mapUserSuperlikedPostEvent(event UserSuperlikedPostEvent) (*model.PostSuperlike, error) {
	return &model.PostSuperlike{
		User: &model.UserMetadata{
			Username: event.Username,
		},
		PostId: event.PostId,
	}, nil
}
