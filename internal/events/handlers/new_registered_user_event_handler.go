package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	userprofile "readmodels/internal/user_profile"
)

type UserWasRegisteredEvent struct {
	UserId   string `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	UserType string `json:"user_type"`
	Region   string `json:"region"`
	FullName string `json:"full_name"`
}

type UserWasRegisteredEventHandler struct {
	service *userprofile.UserProfileService
}

func NewUserWasRegisteredEventHandler(infoLog, errorLog *log.Logger, repository userprofile.UserProfileRepository) *UserWasRegisteredEventHandler {
	return &UserWasRegisteredEventHandler{
		service: userprofile.NewUserProfileService(infoLog, errorLog, repository),
	}
}

func (handler *UserWasRegisteredEventHandler) Handle(event []byte) {
	var userWasRegisteredEvent UserWasRegisteredEvent

	err := Decode(event, &userWasRegisteredEvent)
	if err != nil {
		fmt.Printf("Invalid event data, err: %s\n", err)
		return
	}

	data := mapData(userWasRegisteredEvent)
	handler.service.CreateNewUserProfile(data)
}

func mapData(event UserWasRegisteredEvent) *userprofile.UserProfile {
	return &userprofile.UserProfile{
		UserId:   event.UserId,
		Username: event.Username,
		Name:     event.FullName,
		Bio:      "",
		Link:     "",
	}
}

func Decode(datab []byte, data any) error {
	return json.Unmarshal(datab, &data)
}
