package unit_test_follow

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"readmodels/internal/follow"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
)

var controller *follow.FollowController
var apiResponse *httptest.ResponseRecorder
var ginContext *gin.Context

func setUpController(t *testing.T) {
	setUp(t)
	controller = follow.NewFollowController(repository)
	gin.SetMode(gin.TestMode)
	apiResponse = httptest.NewRecorder()
	ginContext, _ = gin.CreateTestContext(apiResponse)
}

func TestGetFollowersMetadata(t *testing.T) {
	setUpController(t)
	ginContext.Request, _ = http.NewRequest("GET", "/followers", nil)
	followerId1 := "USERA"
	followerId2 := "USERB"
	followerId3 := "USERC"
	u := url.Values{}
	u.Add("followerId", followerId1)
	u.Add("followerId", followerId2)
	u.Add("followerId", followerId3)
	ginContext.Request.URL.RawQuery = u.Encode()
	expectedData := &[]follow.FollowerMetadata{
		{
			Username: followerId1,
			Name:     "fullname1",
		},
		{
			Username: followerId2,
			Name:     "fullname2",
		},
		{
			Username: followerId3,
			Name:     "fullname3",
		},
	}
	repository.EXPECT().GetFollowersMetadata([]string{followerId1, followerId2, followerId3}).Return(expectedData, nil)
	expectedBodyResponse := `{
		"error": false,
		"message": "200 OK",
		"content": {"followers":[
		{
			"username":      "` + followerId1 + `",
			"fullname":   "` + (*expectedData)[0].Name + `"
		},
		{
			"username":      "` + followerId2 + `",
			"fullname":   "` + (*expectedData)[1].Name + `"
		},
		{
			"username":      "` + followerId3 + `",
			"fullname":   "` + (*expectedData)[2].Name + `"
		}
		]}
	}`

	controller.GetFollowersMetadata(ginContext)

	assert.Equal(t, apiResponse.Code, 200)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}

func TestInternalServerErrorOnGetFollowersMetadata(t *testing.T) {
	setUpController(t)
	ginContext.Request, _ = http.NewRequest("GET", "/followers", nil)
	followerId1 := "USERA"
	followerId2 := "USERB"
	followerId3 := "USERC"
	u := url.Values{}
	u.Add("followerId", followerId1)
	u.Add("followerId", followerId2)
	u.Add("followerId", followerId3)
	ginContext.Request.URL.RawQuery = u.Encode()
	expectedData := &[]follow.FollowerMetadata{}
	expectedError := errors.New("some error")
	repository.EXPECT().GetFollowersMetadata([]string{followerId1, followerId2, followerId3}).Return(expectedData, expectedError)
	expectedBodyResponse := `{
		"error": true,
		"message": "` + expectedError.Error() + `",
		"content":null
	}`

	controller.GetFollowersMetadata(ginContext)

	assert.Equal(t, apiResponse.Code, 500)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}
