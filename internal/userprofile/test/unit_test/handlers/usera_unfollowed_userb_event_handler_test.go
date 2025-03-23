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

var userAUnfollowedUserBEventLoggerOutput bytes.Buffer
var userAUnfollowedUserBEventRepository *mock_userprofile.MockRepository
var userAUnfollowedUserBEventHandler *userprofile_handler.UserAUnfollowedUserBEventHandler

func setUpUserAUnfollowedUserBEventHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	userAUnfollowedUserBEventRepository = mock_userprofile.NewMockRepository(ctrl)
	log.Logger = log.Output(&userAUnfollowedUserBEventLoggerOutput)
	userAUnfollowedUserBEventHandler = userprofile_handler.NewUserAUnfollowedUserBEventHandler(userAUnfollowedUserBEventRepository)
}

func TestHandleUserAUnfollowedUserBEventHandler(t *testing.T) {
	setUpUserAUnfollowedUserBEventHandler(t)
	data := &userprofile_handler.UserAUnfollowedUserBEvent{
		FollowerID: "usernameA",
		FolloweeID: "usernameB",
	}
	event, _ := json.Marshal(data)
	userAUnfollowedUserBEventRepository.EXPECT().DecreaseFollowers(data.FolloweeID)
	userAUnfollowedUserBEventRepository.EXPECT().DecreaseFollowees(data.FollowerID)

	userAUnfollowedUserBEventHandler.Handle(event)
}

func TestInvalidDataInUserAUnfollowedUserBEventHandler(t *testing.T) {
	setUpUserAUnfollowedUserBEventHandler(t)
	invalidData := "invalid data"
	event, _ := json.Marshal(invalidData)

	userAUnfollowedUserBEventHandler.Handle(event)

	assert.Contains(t, userAUnfollowedUserBEventLoggerOutput.String(), "Invalid event data")
}
