package userprofile_test

import (
	"bytes"
	"errors"
	"readmodels/internal/model"
	"readmodels/internal/userprofile"
	mock_userprofile "readmodels/internal/userprofile/test/mock"
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
	data := &model.UserProfile{
		Username:        "username1",
		Name:            "user name",
		Bio:             "",
		Link:            "",
		FollowersAmount: 0,
		FolloweesAmount: 0,
	}
	serviceRepository.EXPECT().AddNewUserProfile(data)

	userProfileService.CreateNewUserProfile(data)

	assert.Contains(t, serviceLoggerOutput.String(), "User Profile for user username1 was added")
}

func TestErrorOnCreateNewUserProfileWithService(t *testing.T) {
	setUpService(t)
	data := &model.UserProfile{
		Username:        "username1",
		Name:            "user name",
		Bio:             "",
		Link:            "",
		FollowersAmount: 0,
		FolloweesAmount: 0,
	}
	serviceRepository.EXPECT().AddNewUserProfile(data).Return(errors.New("some error"))

	userProfileService.CreateNewUserProfile(data)

	assert.Contains(t, serviceLoggerOutput.String(), "Error adding user")
}

func TestUpdateUserProfileWithService(t *testing.T) {
	setUpService(t)
	data := &model.UserProfile{
		Username:        "username1",
		Name:            "user name",
		Bio:             "O mellor usuario do mundo",
		Link:            "www.exemplo.com",
		FollowersAmount: 10,
		FolloweesAmount: 20,
	}
	serviceRepository.EXPECT().UpdateUserProfile(data)

	userProfileService.UpdateUserProfile(data)

	assert.Contains(t, serviceLoggerOutput.String(), "User Profile for user username1 was updated")
}

func TestErrorOnUpdateUserProfileWithService(t *testing.T) {
	setUpService(t)
	data := &model.UserProfile{
		Username:        "username1",
		Name:            "user name",
		Bio:             "O mellor usuario do mundo",
		Link:            "www.exemplo.com",
		FollowersAmount: 10,
		FolloweesAmount: 20,
	}
	serviceRepository.EXPECT().UpdateUserProfile(data).Return(errors.New("some error"))

	userProfileService.UpdateUserProfile(data)

	assert.Contains(t, serviceLoggerOutput.String(), "Error updating user")
}

func TestGetUserProfileWithService(t *testing.T) {
	setUpService(t)
	username := "username1"
	expectedData := &model.UserProfile{
		Username:        "username1",
		Name:            "user name",
		Bio:             "",
		Link:            "",
		FollowersAmount: 10,
		FolloweesAmount: 20,
	}
	serviceRepository.EXPECT().GetUserProfile(username).Return(expectedData, nil)

	userProfileService.GetUserProfile(username)
}

func TestErrorOnGetUserProfileWithService(t *testing.T) {
	setUpService(t)
	username := "username1"
	expectedData := &model.UserProfile{}
	serviceRepository.EXPECT().GetUserProfile(username).Return(expectedData, errors.New("some error"))

	userProfileService.GetUserProfile(username)

	assert.Contains(t, serviceLoggerOutput.String(), "Error getting userprofile for username "+username)
}

func TestIncreaseFollowersWithService(t *testing.T) {
	setUpService(t)
	username := "username1"
	serviceRepository.EXPECT().IncreaseFollowers(username).Return(nil)

	userProfileService.IncreaseFollowers(username)
}

func TestErrorOnIncreaseFollowersWithService(t *testing.T) {
	setUpService(t)
	username := "username1"
	serviceRepository.EXPECT().IncreaseFollowers(username).Return(errors.New("some error"))

	userProfileService.IncreaseFollowers(username)

	assert.Contains(t, serviceLoggerOutput.String(), "Error increasing "+username+"'s followers")
}

func TestIncreaseFolloweesWithService(t *testing.T) {
	setUpService(t)
	username := "username1"
	serviceRepository.EXPECT().IncreaseFollowees(username).Return(nil)

	userProfileService.IncreaseFollowees(username)
}

func TestErrorOnIncreaseFolloweesWithService(t *testing.T) {
	setUpService(t)
	username := "username1"
	serviceRepository.EXPECT().IncreaseFollowees(username).Return(errors.New("some error"))

	userProfileService.IncreaseFollowees(username)

	assert.Contains(t, serviceLoggerOutput.String(), "Error increasing "+username+"'s followees")
}

func TestDecreaseFollowersWithService(t *testing.T) {
	setUpService(t)
	username := "username1"
	serviceRepository.EXPECT().DecreaseFollowers(username).Return(nil)

	userProfileService.DecreaseFollowers(username)
}

func TestErrorOnDecreaseFollowersWithService(t *testing.T) {
	setUpService(t)
	username := "username1"
	serviceRepository.EXPECT().DecreaseFollowers(username).Return(errors.New("some error"))

	userProfileService.DecreaseFollowers(username)

	assert.Contains(t, serviceLoggerOutput.String(), "Error decreasing "+username+"'s followers")
}

func TestDecreaseFolloweesWithService(t *testing.T) {
	setUpService(t)
	username := "username1"
	serviceRepository.EXPECT().DecreaseFollowees(username).Return(nil)

	userProfileService.DecreaseFollowees(username)
}

func TestErrorOnDecreaseFolloweesWithService(t *testing.T) {
	setUpService(t)
	username := "username1"
	serviceRepository.EXPECT().DecreaseFollowees(username).Return(errors.New("some error"))

	userProfileService.DecreaseFollowees(username)

	assert.Contains(t, serviceLoggerOutput.String(), "Error decreasing "+username+"'s followees")
}
