package reaction_test

import (
	"errors"
	"net/http"
	"net/url"
	"readmodels/internal/model"
	"readmodels/internal/reaction"
	mock_reaction "readmodels/internal/reaction/test/mock"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
)

var controllerService *mock_reaction.MockControllerService
var controller *reaction.ReactionController

func setUpController(t *testing.T) {
	SetUp(t)
	controllerService = mock_reaction.NewMockControllerService(ctrl)
	controller = reaction.NewReactionController(controllerService)
}

func TestGetPostLikesMetadataWithController_WhenSuccess(t *testing.T) {
	setUpController(t)
	ginContext.Request, _ = http.NewRequest("GET", "/postLikes", nil)
	expectedPostId := "post1"
	expectedLastUsername := "username0"
	expectedLimit := 4
	ginContext.Params = []gin.Param{{Key: "postId", Value: expectedPostId}}
	u := url.Values{}
	u.Add("lastUsername", expectedLastUsername)
	u.Add("limit", strconv.Itoa(expectedLimit))
	ginContext.Request.URL.RawQuery = u.Encode()
	expectedPostLikes := []*model.UserMetadata{
		{
			Username: "username1",
			Name:     "fullname1",
		},
		{
			Username: "username2",
			Name:     "fullname2",
		},
		{
			Username: "username3",
			Name:     "fullname3",
		},
	}
	controllerService.EXPECT().GetPostLikesMetadata(expectedPostId, expectedLastUsername, expectedLimit).Return(expectedPostLikes, "username3", nil)
	expectedBodyResponse := `{
		"error": false,
		"message": "200 OK",
		"content": {
			"users":[	
			{
					"username":  "username1",
					"name":  "fullname1"
			},
			{
					"username":  "username2",
					"name":  "fullname2"
			},
			{
					"username":  "username3",
					"name":  "fullname3"
			}
			],
			"lastUsername":"username3"
		}
	}`

	controller.GetPostLikesMetadata(ginContext)

	assert.Equal(t, apiResponse.Code, 200)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}

func TestGetGetPostLikesMetadataWithController_WhenSuccessWithDefaultPaginationParameters(t *testing.T) {
	setUpController(t)
	ginContext.Request, _ = http.NewRequest("GET", "/postLikes", nil)
	expectedPostId := "post1"
	ginContext.Params = []gin.Param{{Key: "postId", Value: expectedPostId}}
	expectedDefaultLastUsername := ""
	expectedDefaultLimit := 12
	expectedPostLikes := []*model.UserMetadata{
		{
			Username: "username1",
			Name:     "fullname1",
		},
		{
			Username: "username2",
			Name:     "fullname2",
		},
		{
			Username: "username3",
			Name:     "fullname3",
		},
	}
	controllerService.EXPECT().GetPostLikesMetadata(expectedPostId, expectedDefaultLastUsername, expectedDefaultLimit).Return(expectedPostLikes, "username3", nil)
	expectedBodyResponse := `{
		"error": false,
		"message": "200 OK",
		"content": {
			"users":[	
			{
					"username":  "username1",
					"name":  "fullname1"
			},
			{
					"username":  "username2",
					"name":  "fullname2"
			},
			{
					"username":  "username3",
					"name":  "fullname3"
			}
			],
			"lastUsername":"username3"
		}
	}`

	controller.GetPostLikesMetadata(ginContext)

	assert.Equal(t, apiResponse.Code, 200)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}

func TestInternalServerErrorOnGetGetPostLikesMetadataWithController_WhenServiceCallFails(t *testing.T) {
	setUpController(t)
	ginContext.Request, _ = http.NewRequest("GET", "/postLikes", nil)
	expectedPostId := "post1"
	ginContext.Params = []gin.Param{{Key: "postId", Value: expectedPostId}}
	expectedError := errors.New("some error")
	controllerService.EXPECT().GetPostLikesMetadata(expectedPostId, "", 12).Return([]*model.UserMetadata{}, "", expectedError)
	expectedBodyResponse := `{
		"error": true,
		"message": "` + expectedError.Error() + `",
		"content": null
	}`

	controller.GetPostLikesMetadata(ginContext)

	assert.Equal(t, apiResponse.Code, 500)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}

func TestBadRequestErrorOnGetUserPostsWithController_WhenNoPostId(t *testing.T) {
	setUpController(t)
	ginContext.Request, _ = http.NewRequest("GET", "/postLikes", nil)
	expectedError := "Missing parameter postId"
	expectedBodyResponse := `{
		"error": true,
		"message": "` + expectedError + `",
		"content":null
	}`

	controller.GetPostLikesMetadata(ginContext)

	assert.Equal(t, apiResponse.Code, 400)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}

func TestBadRequestErrorOnGetUserPostsWithController_WhenLimitSmallerThanOne(t *testing.T) {
	setUpController(t)
	ginContext.Request, _ = http.NewRequest("GET", "/postLikes", nil)
	expectedPostId := "post1"
	ginContext.Params = []gin.Param{{Key: "postId", Value: expectedPostId}}
	lastUsername := "username0"
	wrongLimit := 0
	u := url.Values{}
	u.Add("lastUsername", lastUsername)
	u.Add("limit", strconv.Itoa(wrongLimit))
	ginContext.Request.URL.RawQuery = u.Encode()
	expectedError := "Invalid pagination parameters, limit must be greater than 0"
	expectedBodyResponse := `{
		"error": true,
		"message": "` + expectedError + `",
		"content":null
	}`

	controller.GetPostLikesMetadata(ginContext)

	assert.Equal(t, apiResponse.Code, 400)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}
