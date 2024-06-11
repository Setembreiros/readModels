package userprofile_handler

import (
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
	UpdateUserProfile(data *userprofile.UserProfile)
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

	err := Decode(event, &userProfileUpdatedEvent)
	if err != nil {
		log.Error().Stack().Err(err).Msg("Invalid event data")
		return
	}

	data := mapToUserProfile(userProfileUpdatedEvent)
	handler.service.UpdateUserProfile(data)
}

func mapToUserProfile(event UserProfileUpdatedEvent) *userprofile.UserProfile {
	return &userprofile.UserProfile{
		Username: event.Username,
		Name:     event.FullName,
		Bio:      event.Bio,
		Link:     event.Link,
	}
}
