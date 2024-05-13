package userprofile_test

import (
	"bytes"
	"errors"
	"log"
	"net/http/httptest"
	database "readmodels/internal/db"
	"readmodels/internal/userprofile"
	mock_userprofile "readmodels/internal/userprofile/mock"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
)

var controllerLoggerOutput bytes.Buffer
var controllerRepository *mock_userprofile.MockRepository
var controller *userprofile.UserProfileController
var apiResponse *httptest.ResponseRecorder
var ginContext *gin.Context

func setUpHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	controllerRepository = mock_userprofile.NewMockRepository(ctrl)
	mockLogger := log.New(&controllerLoggerOutput, "", log.LstdFlags)
	controller = userprofile.NewUserProfileController(mockLogger, mockLogger, controllerRepository)
	gin.SetMode(gin.TestMode)
	apiResponse = httptest.NewRecorder()
	ginContext, _ = gin.CreateTestContext(apiResponse)
}

func TestGetUserProfile(t *testing.T) {
	setUpHandler(t)
	username := "username1"
	ginContext.Params = []gin.Param{{Key: "username", Value: username}}
	data := &userprofile.UserProfile{
		UserId:   "user1",
		Username: "username1",
		Name:     "user name",
		Bio:      "",
		Link:     "",
	}
	controllerRepository.EXPECT().GetUserProfile(username).Return(data, nil)
	expectedBodyResponse := `{
		"error": false,
		"message": "200 OK",
		"content": {
			"user_id": "user1",
			"username": "username1",
			"name": "user name",
			"bio": "",
			"link": ""
		}
	}`

	controller.GetUserProfile(ginContext)

	assert.Equal(t, apiResponse.Code, 200)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}

func TestUserNotFoundOnGetUserProfile(t *testing.T) {
	setUpHandler(t)
	noExistingUsername := "noExistingUsername"
	ginContext.Params = []gin.Param{{Key: "username", Value: noExistingUsername}}
	expectedData := &userprofile.UserProfile{
		UserId:   "",
		Username: "",
		Name:     "",
		Bio:      "",
		Link:     "",
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

func TestInternalServerOnGetUserProfile(t *testing.T) {
	setUpHandler(t)
	username := "username1"
	ginContext.Params = []gin.Param{{Key: "username", Value: username}}
	expectedData := &userprofile.UserProfile{
		UserId:   "",
		Username: "",
		Name:     "",
		Bio:      "",
		Link:     "",
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
