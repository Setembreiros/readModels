package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"readmodels/cmd/provider"
	database "readmodels/internal/db"
	"readmodels/internal/userprofile"
	userprofile_handler "readmodels/internal/userprofile/handlers"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

var db *database.Database
var userAFollowedUserBEventLoggerOutput bytes.Buffer
var userAFollowedUserBEventRepository *userprofile.UserProfileRepository
var userAFollowedUserBEventHandler *userprofile_handler.UserAFollowedUserBEventHandler

func setUpUserAFollowedUserBEventHandler(t *testing.T) {
	ctx, _ := context.WithCancel(context.Background())
	provider := provider.NewProvider("test")
	db, _ = provider.ProvideDb(ctx)
	db.ApplyMigrations(ctx)
	log.Logger = log.Output(&userAFollowedUserBEventLoggerOutput)
	userAFollowedUserBEventHandler = userprofile_handler.NewUserAFollowedUserBEventHandler(userprofile.UserProfileRepository(*db))
}

func tearDown() {
	db.Client.Clean()
}

func TestHandlingUserAFollowedUserBEvent_WhenItReturnsSuccess(t *testing.T) {
	setUpUserAFollowedUserBEventHandler(t)
	defer tearDown()
	userA := &userprofile.UserProfile{
		Username: "usernameA",
		Name:     "user name A",
		Bio:      "",
		Link:     "",
	}
	userB := &userprofile.UserProfile{
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

func AddUserProfileToDatabase(t *testing.T, user *userprofile.UserProfile) {
	err := db.Client.InsertData("UserProfile", user)
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
	var userProfile userprofile.UserProfile
	err := db.Client.GetData("UserProfile", userProfileKey, &userProfile)
	assert.Nil(t, err)
	assert.Equal(t, userProfile.Followees, 0)
	assert.Equal(t, userProfile.Followers, 1)
}

func assertFolloweesIncreased(t *testing.T, username string) {
	userProfileKey := &database.UserProfileKey{
		Username: username,
	}
	var userProfile userprofile.UserProfile
	err := db.Client.GetData("UserProfile", userProfileKey, &userProfile)
	assert.Nil(t, err)
	assert.Equal(t, userProfile.Followers, 0)
	assert.Equal(t, userProfile.Followees, 1)
}
