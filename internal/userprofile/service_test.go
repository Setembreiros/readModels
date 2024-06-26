package userprofile_test

import (
	"bytes"
	"errors"
	"readmodels/internal/userprofile"
	mock_userprofile "readmodels/internal/userprofile/mock"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

var serviceLoggerOutput bytes.Buffer
var serviceRepository *mock_userprofile.MockRepository
var userProfileService *userprofile.UserProfileService

func setUpService(t *testing.T) {
	ctrl := gomock.NewController(t)
	serviceRepository = mock_userprofile.NewMockRepository(ctrl)
	log.Logger = log.Output(&serviceLoggerOutput)
	userProfileService = userprofile.NewUserProfileService(serviceRepository)
}

func TestCreateNewUserProfileWithService(t *testing.T) {
	setUpService(t)
	data := &userprofile.UserProfile{
		Username: "username1",
		Name:     "user name",
		Bio:      "",
		Link:     "",
	}
	serviceRepository.EXPECT().AddNewUserProfile(data)

	userProfileService.CreateNewUserProfile(data)

	assert.Contains(t, serviceLoggerOutput.String(), "User Profile for user username1 was added")
}

func TestErrorOnCreateNewUserProfileWithService(t *testing.T) {
	setUpService(t)
	data := &userprofile.UserProfile{
		Username: "username1",
		Name:     "user name",
		Bio:      "",
		Link:     "",
	}
	serviceRepository.EXPECT().AddNewUserProfile(data).Return(errors.New("some error"))

	userProfileService.CreateNewUserProfile(data)

	assert.Contains(t, serviceLoggerOutput.String(), "Error adding user")
}

func TestUpdateUserProfileWithService(t *testing.T) {
	setUpService(t)
	data := &userprofile.UserProfile{
		Username: "username1",
		Name:     "user name",
		Bio:      "O mellor usuario do mundo",
		Link:     "www.exemplo.com",
	}
	serviceRepository.EXPECT().UpdateUserProfile(data)

	userProfileService.UpdateUserProfile(data)

	assert.Contains(t, serviceLoggerOutput.String(), "User Profile for user username1 was updated")
}

func TestErrorOnUpdateUserProfileWithService(t *testing.T) {
	setUpService(t)
	data := &userprofile.UserProfile{
		Username: "username1",
		Name:     "user name",
		Bio:      "O mellor usuario do mundo",
		Link:     "www.exemplo.com",
	}
	serviceRepository.EXPECT().UpdateUserProfile(data).Return(errors.New("some error"))

	userProfileService.UpdateUserProfile(data)

	assert.Contains(t, serviceLoggerOutput.String(), "Error updating user")
}

func TestGetUserProfileWithService(t *testing.T) {
	setUpService(t)
	username := "username1"
	expectedData := &userprofile.UserProfile{
		Username: "username1",
		Name:     "user name",
		Bio:      "",
		Link:     "",
	}
	serviceRepository.EXPECT().GetUserProfile(username).Return(expectedData, nil)

	userProfileService.GetUserProfile(username)
}

func TestErrorOnGetUserProfileWithService(t *testing.T) {
	setUpService(t)
	username := "username1"
	expectedData := &userprofile.UserProfile{}
	serviceRepository.EXPECT().GetUserProfile(username).Return(expectedData, errors.New("some error"))

	userProfileService.GetUserProfile(username)

	assert.Contains(t, serviceLoggerOutput.String(), "Error getting userprofile for username "+username)
}
