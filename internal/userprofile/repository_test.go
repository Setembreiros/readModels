package userprofile_test

import (
	"log"
	"os"
	database "readmodels/internal/db"
	mock_database "readmodels/internal/db/mock"
	"readmodels/internal/userprofile"
	"testing"

	"github.com/golang/mock/gomock"
)

var client *mock_database.MockDatabaseClient
var userProfileRepository userprofile.UserProfileRepository

func setUp(t *testing.T) {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	ctrl := gomock.NewController(t)
	client = mock_database.NewMockDatabaseClient(ctrl)
	userProfileRepository = userprofile.UserProfileRepository(*database.NewDatabase(client, infoLog))
}

func TestAddNewUserProfileInRepository(t *testing.T) {
	setUp(t)
	data := &userprofile.UserProfile{
		UserId:   "user1",
		Username: "username1",
		Name:     "user name",
		Bio:      "",
		Link:     "",
	}
	client.EXPECT().InsertData("UserProfile", data)

	userProfileRepository.AddNewUserProfile(data)
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