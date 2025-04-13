package userprofile_handler_test

import (
	"bytes"
	"encoding/json"
	"testing"

	userprofile_handler "readmodels/internal/userprofile/handlers"
	mock_userprofile "readmodels/internal/userprofile/test/mock"

	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

var userAFollowedUserBEventLoggerOutput bytes.Buffer
var userAFollowedUserBEventRepository *mock_userprofile.MockRepository
var userAFollowedUserBEventHandler *userprofile_handler.UserAFollowedUserBEventHandler

func setUpUserAFollowedUserBEventHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	userAFollowedUserBEventRepository = mock_userprofile.NewMockRepository(ctrl)
	log.Logger = log.Output(&userAFollowedUserBEventLoggerOutput)
	userAFollowedUserBEventHandler = userprofile_handler.NewUserAFollowedUserBEventHandler(userAFollowedUserBEventRepository)
}

func TestHandleUserAFollowedUserBEventHandler(t *testing.T) {
	setUpUserAFollowedUserBEventHandler(t)
	data := &userprofile_handler.UserAFollowedUserBEvent{
		FollowerID: "usernameA",
		FolloweeID: "usernameB",
	}
	event, _ := json.Marshal(data)
	userAFollowedUserBEventRepository.EXPECT().IncreaseFollowers(data.FolloweeID)
	userAFollowedUserBEventRepository.EXPECT().IncreaseFollowees(data.FollowerID)

	userAFollowedUserBEventHandler.Handle(event)
}

func TestInvalidDataInUserAFollowedUserBEventHandler(t *testing.T) {
	setUpUserAFollowedUserBEventHandler(t)
	invalidData := "invalid data"
	event, _ := json.Marshal(invalidData)

	userAFollowedUserBEventHandler.Handle(event)

	assert.Contains(t, userAFollowedUserBEventLoggerOutput.String(), "Invalid event data")
}
