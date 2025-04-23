package userprofile_handler

import (
	common_data "readmodels/internal/common/data"
	"readmodels/internal/model"
	userprofile "readmodels/internal/userprofile"

	"github.com/rs/zerolog/log"
)

type UserProfileUpdatedEvent struct {
	Username string `json:"username"`
	Bio      string `json:"bio"`
	Link     string `json:"link"`
	FullName string `json:"full_name"`
}

type UserProfileUpdatedEventService interface {
	UpdateUserProfile(data *model.UserProfile)
}

type UserProfileUpdatedEventHandler struct {
	service UserProfileUpdatedEventService
}

func NewUserProfileUpdatedEventHandler(repository userprofile.Repository) *UserProfileUpdatedEventHandler {
	return &UserProfileUpdatedEventHandler{
		service: userprofile.NewUserProfileService(repository),
	}
}

func (handler *UserProfileUpdatedEventHandler) Handle(event []byte) {
	var userProfileUpdatedEvent UserProfileUpdatedEvent
	log.Info().Msg("Handling UserProfileUpdatedEvent")

	err := common_data.DeserializeData(event, &userProfileUpdatedEvent)
	if err != nil {
		log.Error().Stack().Err(err).Msg("Invalid event data")
		return
	}

	data := mapToUserProfile(userProfileUpdatedEvent)
	handler.service.UpdateUserProfile(data)
}

func mapToUserProfile(event UserProfileUpdatedEvent) *model.UserProfile {
	return &model.UserProfile{
		Username: event.Username,
		Name:     event.FullName,
		Bio:      event.Bio,
		Link:     event.Link,
	}
}
