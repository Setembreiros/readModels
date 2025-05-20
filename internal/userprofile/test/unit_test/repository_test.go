package userprofile_test

import (
	database "readmodels/internal/db"
	mock_database "readmodels/internal/db/test/mock"
	"readmodels/internal/model"
	"readmodels/internal/userprofile"
	"testing"

	"github.com/golang/mock/gomock"
)

var client *mock_database.MockDatabaseClient
var userProfileRepository userprofile.UserProfileRepository

func setUp(t *testing.T) {
	ctrl := gomock.NewController(t)
	client = mock_database.NewMockDatabaseClient(ctrl)
	userProfileRepository = userprofile.UserProfileRepository(*database.NewDatabase(client))
}

func TestAddNewUserProfileInRepository(t *testing.T) {
	setUp(t)
	data := &model.UserProfile{
		Username: "username1",
		Name:     "user name",
		Bio:      "",
		Link:     "",
	}
	client.EXPECT().InsertData("UserProfile", data)

	userProfileRepository.AddNewUserProfile(data)
}

func TestUpdateUserProfileInRepository(t *testing.T) {
	setUp(t)
	username := "username1"
	expectedUserProfileKey := &database.UserProfileKey{
		Username: username,
	}
	data := &model.UserProfile{
		Username: username,
		Name:     "user name",
		Bio:      "O mellor usuario do mundo",
		Link:     "www.exemplo.com",
	}
	updateAttributes := map[string]interface{}{
		"Name": data.Name,
		"Bio":  data.Bio,
		"Link": data.Link,
	}
	client.EXPECT().UpdateData("UserProfile", expectedUserProfileKey, updateAttributes)

	userProfileRepository.UpdateUserProfile(data)
}

func TestGetUserProfileFromRepository(t *testing.T) {
	setUp(t)
	username := "username1"
	var userProfile model.UserProfile
	expectedUserProfileKey := &database.UserProfileKey{
		Username: username,
	}
	client.EXPECT().GetData("UserProfile", expectedUserProfileKey, &userProfile)

	userProfileRepository.GetUserProfile(username)
}

func TestIncreaseFollowersFromRepository(t *testing.T) {
	setUp(t)
	username := "username1"
	expectedUserProfileKey := &database.UserProfileKey{
		Username: username,
	}
	client.EXPECT().IncrementCounter("UserProfile", expectedUserProfileKey, "FollowersAmount", 1)

	userProfileRepository.IncreaseFollowers(username)
}

func TestIncreaseFolloweesFromRepository(t *testing.T) {
	setUp(t)
	username := "username1"
	expectedUserProfileKey := &database.UserProfileKey{
		Username: username,
	}
	client.EXPECT().IncrementCounter("UserProfile", expectedUserProfileKey, "FolloweesAmount", 1)

	userProfileRepository.IncreaseFollowees(username)
}

func TestDecreaseFollowersFromRepository(t *testing.T) {
	setUp(t)
	username := "username1"
	expectedUserProfileKey := &database.UserProfileKey{
		Username: username,
	}
	client.EXPECT().IncrementCounter("UserProfile", expectedUserProfileKey, "FollowersAmount", -1)

	userProfileRepository.DecreaseFollowers(username)
}

func TestDecreaseFolloweesFromRepository(t *testing.T) {
	setUp(t)
	username := "username1"
	expectedUserProfileKey := &database.UserProfileKey{
		Username: username,
	}
	client.EXPECT().IncrementCounter("UserProfile", expectedUserProfileKey, "FolloweesAmount", -1)

	userProfileRepository.DecreaseFollowees(username)
}
