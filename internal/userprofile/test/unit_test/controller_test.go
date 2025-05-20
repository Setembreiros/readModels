package userprofile_test

import (
	"bytes"
	"errors"
	"net/http/httptest"
	database "readmodels/internal/db"
	"readmodels/internal/model"
	"readmodels/internal/userprofile"
	mock_userprofile "readmodels/internal/userprofile/test/mock"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog/log"
)

var controllerLoggerOutput bytes.Buffer
var controllerRepository *mock_userprofile.MockRepository
var controller *userprofile.UserProfileController
var apiResponse *httptest.ResponseRecorder
var ginContext *gin.Context

func setUpController(t *testing.T) {
	ctrl := gomock.NewController(t)
	controllerRepository = mock_userprofile.NewMockRepository(ctrl)
	log.Logger = log.Output(&controllerLoggerOutput)
	controller = userprofile.NewUserProfileController(controllerRepository)
	gin.SetMode(gin.TestMode)
	apiResponse = httptest.NewRecorder()
	ginContext, _ = gin.CreateTestContext(apiResponse)
}

func TestGetUserProfile(t *testing.T) {
	setUpController(t)
	username := "username1"
	ginContext.Params = []gin.Param{{Key: "username", Value: username}}
	data := &model.UserProfile{
		Username:        "username1",
		Name:            "user name",
		Bio:             "",
		Link:            "",
		FollowersAmount: 10,
		FolloweesAmount: 20,
		PostsAmount:     60,
	}
	controllerRepository.EXPECT().GetUserProfile(username).Return(data, nil)
	expectedBodyResponse := `{
		"error": false,
		"message": "200 OK",
		"content": {
			"username": "username1",
			"name": "user name",
			"bio": "",
			"link": "",
			"followersAmount": 10,
			"followeesAmount": 20,
			"postsAmount": 60
		}
	}`

	controller.GetUserProfile(ginContext)

	assert.Equal(t, apiResponse.Code, 200)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}

func TestUserNotFoundOnGetUserProfile(t *testing.T) {
	setUpController(t)
	noExistingUsername := "noExistingUsername"
	ginContext.Params = []gin.Param{{Key: "username", Value: noExistingUsername}}
	expectedData := &model.UserProfile{
		Username:        "",
		Name:            "",
		Bio:             "",
		Link:            "",
		FollowersAmount: 10,
		FolloweesAmount: 20,
	}
	expectedNotFoundError := &database.NotFoundError{}
	controllerRepository.EXPECT().GetUserProfile(noExistingUsername).Return(expectedData, expectedNotFoundError)
	expectedBodyResponse := `{
		"error": true,
		"message": "User Profile not found for username ` + noExistingUsername + `",
		"content":null
	}`

	controller.GetUserProfile(ginContext)

	assert.Equal(t, apiResponse.Code, 404)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}

func TestInternalServerErrorOnGetUserProfile(t *testing.T) {
	setUpController(t)
	username := "username1"
	ginContext.Params = []gin.Param{{Key: "username", Value: username}}
	expectedData := &model.UserProfile{
		Username:        "",
		Name:            "",
		Bio:             "",
		Link:            "",
		FollowersAmount: 10,
		FolloweesAmount: 20,
	}
	expectedError := errors.New("some error")
	controllerRepository.EXPECT().GetUserProfile(username).Return(expectedData, expectedError)
	expectedBodyResponse := `{
		"error": true,
		"message": "` + expectedError.Error() + `",
		"content":null
	}`

	controller.GetUserProfile(ginContext)

	assert.Equal(t, apiResponse.Code, 500)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}

func removeSpace(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(s, " ", ""), "\t", ""), "\n", "")
}
