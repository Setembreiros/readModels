package userprofile_handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	userprofile "readmodels/internal/userprofile"
	userprofile_handler "readmodels/internal/userprofile/handlers"
	mock_userprofile "readmodels/internal/userprofile/mock"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

var loggerOutput bytes.Buffer
var repository *mock_userprofile.MockRepository
var handler *userprofile_handler.UserWasRegisteredEventHandler

func setUpHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	repository = mock_userprofile.NewMockRepository(ctrl)
	log.Logger = log.Output(&loggerOutput)
	handler = userprofile_handler.NewUserWasRegisteredEventHandler(repository)
}

func TestHandleUserWasRegisteredEventHandler(t *testing.T) {
	setUpHandler(t)
	data := &userprofile_handler.UserWasRegisteredEvent{
		UserId:   "user1",
		Username: "username1",
		Email:    "email1",
		UserType: "UA",
		Region:   "Vigo",
		FullName: "user lastname",
	}
	event, _ := json.Marshal(data)
	expectedUserprofile := &userprofile.UserProfile{
		UserId:   "user1",
		Username: "username1",
		Name:     "user lastname",
		Bio:      "",
		Link:     "",
	}
	repository.EXPECT().AddNewUserProfile(expectedUserprofile)

	handler.Handle(event)
}

func TestInvalidDataInUserWasRegisteredEventHandler(t *testing.T) {
	setUpHandler(t)
	invalidData := "invalid data"
	event, _ := json.Marshal(invalidData)

	handler.Handle(event)

	fmt.Println(loggerOutput.String())
	assert.Contains(t, loggerOutput.String(), "Invalid event data")
}
