package reaction_handler

import (
	common_data "readmodels/internal/common/data"
	"readmodels/internal/model"

	"github.com/rs/zerolog/log"
)

//go:generate mockgen -source=user_unsuperliked_post_event_handler.go -destination=test/mock/user_unsuperliked_post_event_handler.go

type UserUnsuperlikedPostEvent struct {
	Username string `json:"username"`
	PostId   string `json:"postId"`
}

type UserUnsuperlikedPostEventService interface {
	DeletePostSuperlike(data *model.PostSuperlike)
}

type UserUnsuperlikedPostEventHandler struct {
	service UserUnsuperlikedPostEventService
}

func NewUserUnsuperlikedPostEventHandler(service UserUnsuperlikedPostEventService) *UserUnsuperlikedPostEventHandler {
	return &UserUnsuperlikedPostEventHandler{
		service: service,
	}
}

func (handler *UserUnsuperlikedPostEventHandler) Handle(event []byte) {
	var userUnsuperlikedPostEvent UserUnsuperlikedPostEvent
	log.Info().Msg("Handling UserUnsuperlikedPostEvent")

	err := common_data.DeserializeData(event, &userUnsuperlikedPostEvent)
	if err != nil {
		log.Error().Stack().Err(err).Msg("Invalid event data")
		return
	}

	data, err := mapUserUnsuperlikedPostEvent(userUnsuperlikedPostEvent)
	if err != nil {
		return
	}

	handler.service.DeletePostSuperlike(data)
}

func mapUserUnsuperlikedPostEvent(event UserUnsuperlikedPostEvent) (*model.PostSuperlike, error) {
	return &model.PostSuperlike{
		User: &model.UserMetadata{
			Username: event.Username,
		},
		PostId: event.PostId,
	}, nil
}
