package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"readmodels/cmd/provider"
	database "readmodels/internal/db"
	"readmodels/internal/model"
	"readmodels/internal/userprofile"
	userprofile_handler "readmodels/internal/userprofile/handlers"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

var userAFollowedUserBEventDb *database.Database
var userAFollowedUserBEventLoggerOutput bytes.Buffer
var userAFollowedUserBEventRepository *userprofile.UserProfileRepository
var userAFollowedUserBEventHandler *userprofile_handler.UserAFollowedUserBEventHandler

func setUpUserAFollowedUserBEventHandler(t *testing.T) {
	ctx, _ := context.WithCancel(context.Background())
	provider := provider.NewProvider("test")
	userAFollowedUserBEventDb, _ = provider.ProvideDb(ctx)
	userAFollowedUserBEventDb.ApplyMigrations(ctx)
	log.Logger = log.Output(&userAFollowedUserBEventLoggerOutput)
	userAFollowedUserBEventHandler = userprofile_handler.NewUserAFollowedUserBEventHandler(userprofile.UserProfileRepository(*userAFollowedUserBEventDb))
}

func tearDownUserFollowedUserBEvent() {
	userAFollowedUserBEventDb.Client.Truncate()
}

func TestHandlingUserAFollowedUserBEvent_WhenItReturnsSuccess(t *testing.T) {
	setUpUserAFollowedUserBEventHandler(t)
	defer tearDownUserFollowedUserBEvent()
	userA := &model.UserProfile{
		Username: "usernameA",
		Name:     "user name A",
		Bio:      "",
		Link:     "",
	}
	userB := &model.UserProfile{
		Username: "usernameB",
		Name:     "user name B",
		Bio:      "",
		Link:     "",
	}
	AddUserProfileToDatabase(t, userA)
	AddUserProfileToDatabase(t, userB)
	expectedUserAUnfollowedUserBEvent := &userprofile_handler.UserAFollowedUserBEvent{
		FollowerID: userA.Username,
		FolloweeID: userB.Username,
	}
	event := createEvent(expectedUserAUnfollowedUserBEvent)

	userAFollowedUserBEventHandler.Handle(event)

	assertFollowersIncreased(t, userB.Username)
	assertFolloweesIncreased(t, userA.Username)
}

func AddUserProfileToDatabase(t *testing.T, user *model.UserProfile) {
	err := userAFollowedUserBEventDb.Client.InsertData("UserProfile", user)
	assert.Nil(t, err)
}

func createEvent(eventData any) []byte {
	dataEvent, err := serializeData(eventData)
	if err != nil {
		return nil
	}

	return dataEvent
}

func serializeData(data any) ([]byte, error) {
	return json.Marshal(data)
}

func assertFollowersIncreased(t *testing.T, username string) {
	userProfileKey := &database.UserProfileKey{
		Username: username,
	}
	var userProfile model.UserProfile
	err := userAFollowedUserBEventDb.Client.GetData("UserProfile", userProfileKey, &userProfile)
	assert.Nil(t, err)
	assert.Equal(t, userProfile.FolloweesAmount, 0)
	assert.Equal(t, userProfile.FollowersAmount, 1)
}

func assertFolloweesIncreased(t *testing.T, username string) {
	userProfileKey := &database.UserProfileKey{
		Username: username,
	}
	var userProfile model.UserProfile
	err := userAFollowedUserBEventDb.Client.GetData("UserProfile", userProfileKey, &userProfile)
	assert.Nil(t, err)
	assert.Equal(t, userProfile.FollowersAmount, 0)
	assert.Equal(t, userProfile.FolloweesAmount, 1)
}
