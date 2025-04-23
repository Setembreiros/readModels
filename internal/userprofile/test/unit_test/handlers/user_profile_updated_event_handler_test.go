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

var userProfileUpdatedEventLoggerOutput bytes.Buffer
var userProfileUpdatedEventRepository *mock_userprofile.MockRepository
var userProfileUpdatedEventHandler *userprofile_handler.UserProfileUpdatedEventHandler

func setUpuserProfileUpdatedEventHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	userProfileUpdatedEventRepository = mock_userprofile.NewMockRepository(ctrl)
	log.Logger = log.Output(&userProfileUpdatedEventLoggerOutput)
	userProfileUpdatedEventHandler = userprofile_handler.NewUserProfileUpdatedEventHandler(userProfileUpdatedEventRepository)
}

func TestHandleUserProfileUpdatedEventHandler(t *testing.T) {
	setUpuserProfileUpdatedEventHandler(t)
	data := &userprofile_handler.UserProfileUpdatedEvent{
		Username: "username1",
		Bio:      "o mellor usuario do mundo",
		Link:     "www.exemplo.com",
		FullName: "user lastname",
	}
	event, _ := json.Marshal(data)
	expectedUserprofile := &model.UserProfile{
		Username: "username1",
		Name:     "user lastname",
		Bio:      "o mellor usuario do mundo",
		Link:     "www.exemplo.com",
	}
	userProfileUpdatedEventRepository.EXPECT().UpdateUserProfile(expectedUserprofile)

	userProfileUpdatedEventHandler.Handle(event)
}

func TestInvalidDataInUserProfileUpdatedEventHandler(t *testing.T) {
	setUpuserProfileUpdatedEventHandler(t)
	invalidData := "invalid data"
	event, _ := json.Marshal(invalidData)

	userProfileUpdatedEventHandler.Handle(event)

	assert.Contains(t, userProfileUpdatedEventLoggerOutput.String(), "Invalid event data")
}
