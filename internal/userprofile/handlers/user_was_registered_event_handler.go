package userprofile_handler

import (
	"encoding/json"
	"log"
	userprofile "readmodels/internal/userprofile"
)

type UserWasRegisteredEvent struct {
	UserId   string `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	UserType string `json:"user_type"`
	Region   string `json:"region"`
	FullName string `json:"full_name"`
}

type UserWasRegisteredEvenService interface {
	CreateNewUserProfile(data *userprofile.UserProfile)
}

type UserWasRegisteredEventHandler struct {
	service  UserWasRegisteredEvenService
	infoLog  *log.Logger
	errorLog *log.Logger
}

func NewUserWasRegisteredEventHandler(infoLog, errorLog *log.Logger, repository userprofile.Repository) *UserWasRegisteredEventHandler {
	return &UserWasRegisteredEventHandler{
		service:  userprofile.NewUserProfileService(infoLog, errorLog, repository),
		infoLog:  infoLog,
		errorLog: errorLog,
	}
}

func (handler *UserWasRegisteredEventHandler) Handle(event []byte) {
	var userWasRegisteredEvent UserWasRegisteredEvent
	handler.infoLog.Printf("Handling UserWasRegisteredEvent\n")

	err := Decode(event, &userWasRegisteredEvent)
	if err != nil {
		handler.errorLog.Printf("Invalid event data, err: %s\n", err)
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
