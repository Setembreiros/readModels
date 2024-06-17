package userprofile_handler

import (
	"encoding/json"
	userprofile "readmodels/internal/userprofile"

	"github.com/rs/zerolog/log"
)

type UserWasRegisteredEvent struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	UserType string `json:"user_type"`
	Region   string `json:"region"`
	FullName string `json:"full_name"`
}

type UserWasRegisteredEventService interface {
	CreateNewUserProfile(data *userprofile.UserProfile)
}

type UserWasRegisteredEventHandler struct {
	service UserWasRegisteredEventService
}

func NewUserWasRegisteredEventHandler(repository userprofile.Repository) *UserWasRegisteredEventHandler {
	return &UserWasRegisteredEventHandler{
		service: userprofile.NewUserProfileService(repository),
	}
}

func (handler *UserWasRegisteredEventHandler) Handle(event []byte) {
	var userWasRegisteredEvent UserWasRegisteredEvent
	log.Info().Msg("Handling UserWasRegisteredEvent")

	err := Decode(event, &userWasRegisteredEvent)
	if err != nil {
		log.Error().Stack().Err(err).Msg("Invalid event data")
		return
	}

	data := mapData(userWasRegisteredEvent)
	handler.service.CreateNewUserProfile(data)
}

func mapData(event UserWasRegisteredEvent) *userprofile.UserProfile {
	return &userprofile.UserProfile{
		Username: event.Username,
		Name:     event.FullName,
		Bio:      "",
		Link:     "",
	}
}

func Decode(datab []byte, data any) error {
	return json.Unmarshal(datab, &data)
}
