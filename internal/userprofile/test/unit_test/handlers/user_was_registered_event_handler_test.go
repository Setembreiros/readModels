package userprofile_handler_test

import (
	"bytes"
	"encoding/json"
	"readmodels/internal/model"
	userprofile_handler "readmodels/internal/userprofile/handlers"
	mock_userprofile "readmodels/internal/userprofile/test/mock"
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
		Username: "username1",
		Email:    "email1",
		UserType: "UA",
		Region:   "Vigo",
		FullName: "user lastname",
	}
	event, _ := json.Marshal(data)
	expectedUserprofile := &model.UserProfile{
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

	assert.Contains(t, loggerOutput.String(), "Invalid event data")
}
