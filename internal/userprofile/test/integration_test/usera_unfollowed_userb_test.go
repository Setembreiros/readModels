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

var userAUnfollowedUserBEventDb *database.Database
var userAUnfollowedUserBEventLoggerOutput bytes.Buffer
var userAUnfollowedUserBEventHandler *userprofile_handler.UserAUnfollowedUserBEventHandler

func setUpUserAUnfollowedUserBEventHandler() {
	ctx, _ := context.WithCancel(context.Background())
	provider := provider.NewProvider("test")
	userAUnfollowedUserBEventDb, _ = provider.ProvideDb(ctx)
	userAUnfollowedUserBEventDb.ApplyMigrations(ctx)
	log.Logger = log.Output(&userAUnfollowedUserBEventLoggerOutput)
	userAUnfollowedUserBEventHandler = userprofile_handler.NewUserAUnfollowedUserBEventHandler(userprofile.UserProfileRepository(*userAUnfollowedUserBEventDb))
}

func tearDownUserAUnfollowedUserBEvent() {
	userAUnfollowedUserBEventDb.Client.Clean()
}

func TestHandlingUserAUnfollowedUserBEvent_WhenItReturnsSuccess(t *testing.T) {
	setUpUserAUnfollowedUserBEventHandler()
	defer tearDownUserAUnfollowedUserBEvent()
	userA := &model.UserProfile{
		Username:  "usernameA",
		Name:      "user name A",
		Bio:       "",
		Link:      "",
		Followers: 1,
		Followees: 1,
	}
	userB := &model.UserProfile{
		Username:  "usernameB",
		Name:      "user name B",
		Bio:       "",
		Link:      "",
		Followers: 1,
		Followees: 1,
	}
	AddUserProfileTouserAUnfollowedUserBEventDatabase(t, userA)
	AddUserProfileTouserAUnfollowedUserBEventDatabase(t, userB)
	event := createUnfollowedUserBEvent(userA.Username, userB.Username)

	userAUnfollowedUserBEventHandler.Handle(event)

	assertFollowersDecreased(t, userB.Username)
	assertFolloweesDecreased(t, userA.Username)
}

func AddUserProfileTouserAUnfollowedUserBEventDatabase(t *testing.T, user *model.UserProfile) {
	err := userAUnfollowedUserBEventDb.Client.InsertData("UserProfile", user)
	assert.Nil(t, err)
}

func createUnfollowedUserBEvent(followerId, followeeId string) []byte {
	eventData := &userprofile_handler.UserAUnfollowedUserBEvent{
		FollowerID: followerId,
		FolloweeID: followeeId,
	}

	dataEvent, err := json.Marshal(eventData)
	if err != nil {
		return nil
	}

	return dataEvent
}

func assertFollowersDecreased(t *testing.T, username string) {
	userProfileKey := &database.UserProfileKey{
		Username: username,
	}
	var userProfile model.UserProfile
	err := userAUnfollowedUserBEventDb.Client.GetData("UserProfile", userProfileKey, &userProfile)
	assert.Nil(t, err)
	assert.Equal(t, userProfile.Followees, 1)
	assert.Equal(t, userProfile.Followers, 0)
}

func assertFolloweesDecreased(t *testing.T, username string) {
	userProfileKey := &database.UserProfileKey{
		Username: username,
	}
	var userProfile model.UserProfile
	err := userAUnfollowedUserBEventDb.Client.GetData("UserProfile", userProfileKey, &userProfile)
	assert.Nil(t, err)
	assert.Equal(t, userProfile.Followers, 1)
	assert.Equal(t, userProfile.Followees, 0)
}
