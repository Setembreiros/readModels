package userprofile_test

import (
	database "readmodels/internal/db"
	mock_database "readmodels/internal/db/test/mock"
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
	data := &userprofile.UserProfile{
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
	expectedUserProfileKey := &userprofile.UserProfileKey{
		Username: username,
	}
	data := &userprofile.UserProfile{
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
	var userProfile userprofile.UserProfile
	expectedUserProfileKey := &userprofile.UserProfileKey{
		Username: username,
	}
	client.EXPECT().GetData("UserProfile", expectedUserProfileKey, &userProfile)

	userProfileRepository.GetUserProfile(username)
}

func TestIncreaseFollowersFromRepository(t *testing.T) {
	setUp(t)
	username := "username1"
	expectedUserProfileKey := &userprofile.UserProfileKey{
		Username: username,
	}
	client.EXPECT().IncrementCounter("UserProfile", expectedUserProfileKey, "Followers", 1)

	userProfileRepository.IncreaseFollowers(username)
}

func TestIncreaseFolloweesFromRepository(t *testing.T) {
	setUp(t)
	username := "username1"
	expectedUserProfileKey := &userprofile.UserProfileKey{
		Username: username,
	}
	client.EXPECT().IncrementCounter("UserProfile", expectedUserProfileKey, "Followees", 1)

	userProfileRepository.IncreaseFollowees(username)
}
