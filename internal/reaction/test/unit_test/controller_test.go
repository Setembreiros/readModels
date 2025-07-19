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
	"time"

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
	controllerService.EXPECT().GetLikesMetadataByPostId(expectedPostId, expectedLastUsername, expectedLimit).Return(expectedPostLikes, "username3", nil)
	expectedBodyResponse := `{
		"error": false,
		"message": "200 OK",
		"content": {
			"postLikes":[	
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

func TestGetPostLikesMetadataWithController_WhenSuccessWithDefaultPaginationParameters(t *testing.T) {
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
	controllerService.EXPECT().GetLikesMetadataByPostId(expectedPostId, expectedDefaultLastUsername, expectedDefaultLimit).Return(expectedPostLikes, "username3", nil)
	expectedBodyResponse := `{
		"error": false,
		"message": "200 OK",
		"content": {
			"postLikes":[	
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

func TestInternalServerErrorOnGetPostLikesMetadataWithController_WhenServiceCallFails(t *testing.T) {
	setUpController(t)
	ginContext.Request, _ = http.NewRequest("GET", "/postLikes", nil)
	expectedPostId := "post1"
	ginContext.Params = []gin.Param{{Key: "postId", Value: expectedPostId}}
	expectedError := errors.New("some error")
	controllerService.EXPECT().GetLikesMetadataByPostId(expectedPostId, "", 12).Return([]*model.UserMetadata{}, "", expectedError)
	expectedBodyResponse := `{
		"error": true,
		"message": "` + expectedError.Error() + `",
		"content": null
	}`

	controller.GetPostLikesMetadata(ginContext)

	assert.Equal(t, apiResponse.Code, 500)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}

func TestBadRequestErrorOnGetUserPostLikesWithController_WhenNoPostId(t *testing.T) {
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

func TestBadRequestErrorOnGetPostLikesWithController_WhenLimitSmallerThanOne(t *testing.T) {
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

func TestGetPostSuperlikesMetadataWithController_WhenSuccess(t *testing.T) {
	setUpController(t)
	ginContext.Request, _ = http.NewRequest("GET", "/postSuperlikes", nil)
	expectedPostId := "post1"
	expectedLastUsername := "username0"
	expectedLimit := 4
	ginContext.Params = []gin.Param{{Key: "postId", Value: expectedPostId}}
	u := url.Values{}
	u.Add("lastUsername", expectedLastUsername)
	u.Add("limit", strconv.Itoa(expectedLimit))
	ginContext.Request.URL.RawQuery = u.Encode()
	expectedPostSuperlikes := []*model.UserMetadata{
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
	controllerService.EXPECT().GetSuperlikesMetadataByPostId(expectedPostId, expectedLastUsername, expectedLimit).Return(expectedPostSuperlikes, "username3", nil)
	expectedBodyResponse := `{
		"error": false,
		"message": "200 OK",
		"content": {
			"postSuperlikes":[	
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

	controller.GetPostSuperlikesMetadata(ginContext)

	assert.Equal(t, apiResponse.Code, 200)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}

func TestGetPostSuperlikesMetadataWithController_WhenSuccessWithDefaultPaginationParameters(t *testing.T) {
	setUpController(t)
	ginContext.Request, _ = http.NewRequest("GET", "/postSuperlikes", nil)
	expectedPostId := "post1"
	ginContext.Params = []gin.Param{{Key: "postId", Value: expectedPostId}}
	expectedDefaultLastUsername := ""
	expectedDefaultLimit := 12
	expectedPostSuperlikes := []*model.UserMetadata{
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
	controllerService.EXPECT().GetSuperlikesMetadataByPostId(expectedPostId, expectedDefaultLastUsername, expectedDefaultLimit).Return(expectedPostSuperlikes, "username3", nil)
	expectedBodyResponse := `{
		"error": false,
		"message": "200 OK",
		"content": {
			"postSuperlikes":[	
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

	controller.GetPostSuperlikesMetadata(ginContext)

	assert.Equal(t, apiResponse.Code, 200)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}

func TestInternalServerErrorOnGetPostSuperlikesMetadataWithController_WhenServiceCallFails(t *testing.T) {
	setUpController(t)
	ginContext.Request, _ = http.NewRequest("GET", "/postSuperlikes", nil)
	expectedPostId := "post1"
	ginContext.Params = []gin.Param{{Key: "postId", Value: expectedPostId}}
	expectedError := errors.New("some error")
	controllerService.EXPECT().GetSuperlikesMetadataByPostId(expectedPostId, "", 12).Return([]*model.UserMetadata{}, "", expectedError)
	expectedBodyResponse := `{
		"error": true,
		"message": "` + expectedError.Error() + `",
		"content": null
	}`

	controller.GetPostSuperlikesMetadata(ginContext)

	assert.Equal(t, apiResponse.Code, 500)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}

func TestBadRequestErrorOnGetPostsSuperlikesWithController_WhenNoPostId(t *testing.T) {
	setUpController(t)
	ginContext.Request, _ = http.NewRequest("GET", "/postSuperlikes", nil)
	expectedError := "Missing parameter postId"
	expectedBodyResponse := `{
		"error": true,
		"message": "` + expectedError + `",
		"content":null
	}`

	controller.GetPostSuperlikesMetadata(ginContext)

	assert.Equal(t, apiResponse.Code, 400)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}

func TestBadRequestErrorOnGetPostSuperlikessWithController_WhenLimitSmallerThanOne(t *testing.T) {
	setUpController(t)
	ginContext.Request, _ = http.NewRequest("GET", "/postSuperlikes", nil)
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

	controller.GetPostSuperlikesMetadata(ginContext)

	assert.Equal(t, apiResponse.Code, 400)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}

func TestGetReviewsByPostIdWithController_WhenSuccess(t *testing.T) {
	setUpController(t)
	ginContext.Request, _ = http.NewRequest("GET", "/reviews", nil)
	expectedPostId := "post1"
	expectedLastReviewId := uint64(4)
	expectedLimit := 4
	ginContext.Params = []gin.Param{{Key: "postId", Value: expectedPostId}}
	u := url.Values{}
	u.Add("lastReviewId", strconv.FormatUint(expectedLastReviewId, 10))
	u.Add("limit", strconv.Itoa(expectedLimit))
	ginContext.Request.URL.RawQuery = u.Encode()
	timeNowString := time.Now().UTC().Format(model.TimeLayout)
	timeNow, _ := time.Parse(model.TimeLayout, timeNowString)
	expectedReviews := []*model.Review{
		{
			ReviewId:  uint64(5),
			Username:  "username1",
			PostId:    "post1",
			Content:   "a miña review 1",
			Rating:    3,
			CreatedAt: timeNow,
			UpdatedAt: timeNow,
		},
		{
			ReviewId:  uint64(6),
			Username:  "username2",
			PostId:    "post1",
			Content:   "a miña review 2",
			Rating:    3,
			CreatedAt: timeNow,
			UpdatedAt: timeNow,
		},
		{
			ReviewId:  uint64(7),
			Username:  "username1",
			PostId:    "post1",
			Content:   "a miña review 3",
			Rating:    3,
			CreatedAt: timeNow,
			UpdatedAt: timeNow,
		},
	}
	controllerService.EXPECT().GetReviewsByPostId(expectedPostId, expectedLastReviewId, expectedLimit).Return(expectedReviews, uint64(7), nil)
	expectedBodyResponse := `{
		"error": false,
		"message": "200 OK",
		"content": {
			"reviews":[	
			{
				"reviewId": 5,
				"postId":    "post1",
				"username":  "username1",
				"content": "a miña review 1",
				"rating": 3,
				"createdAt": "` + timeNowString + `",
				"updatedAt": "` + timeNowString + `"
			},
			{
				"reviewId": 6,
				"postId":    "post1",
				"username":  "username2",
				"content": "a miña review 2",
				"rating": 3,
				"createdAt": "` + timeNowString + `",
				"updatedAt": "` + timeNowString + `"
			},
			{
				"reviewId": 7,
				"postId":    "post1",
				"username":  "username1",
				"content": "a miña review 3",
				"rating": 3,
				"createdAt": "` + timeNowString + `",
				"updatedAt": "` + timeNowString + `"
			}
			],
			"lastReviewId":7
		}
	}`

	controller.GetReviewsByPostId(ginContext)

	assert.Equal(t, apiResponse.Code, 200)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}

func TestGetReviewsByPostIdWithController_WhenSuccessWithDefaultPaginationParameters(t *testing.T) {
	setUpController(t)
	ginContext.Request, _ = http.NewRequest("GET", "/reviews", nil)
	expectedPostId := "post1"
	ginContext.Params = []gin.Param{{Key: "postId", Value: expectedPostId}}
	expectedDefaultLastReviewId := uint64(0)
	expectedDefaultLimit := 12
	timeNowString := time.Now().UTC().Format(model.TimeLayout)
	timeNow, _ := time.Parse(model.TimeLayout, timeNowString)
	expectedReviews := []*model.Review{
		{
			ReviewId:  uint64(5),
			Username:  "username1",
			PostId:    "post1",
			Content:   "a miña review 1",
			Rating:    3,
			CreatedAt: timeNow,
			UpdatedAt: timeNow,
		},
		{
			ReviewId:  uint64(6),
			Username:  "username2",
			PostId:    "post1",
			Content:   "a miña review 2",
			Rating:    3,
			CreatedAt: timeNow,
			UpdatedAt: timeNow,
		},
		{
			ReviewId:  uint64(7),
			Username:  "username1",
			PostId:    "post1",
			Content:   "a miña review 3",
			Rating:    3,
			CreatedAt: timeNow,
			UpdatedAt: timeNow,
		},
	}
	controllerService.EXPECT().GetReviewsByPostId(expectedPostId, expectedDefaultLastReviewId, expectedDefaultLimit).Return(expectedReviews, uint64(7), nil)
	expectedBodyResponse := `{
		"error": false,
		"message": "200 OK",
		"content": {
			"reviews":[	
			{
				"reviewId": 5,
				"postId":    "post1",
				"username":  "username1",
				"content": "a miña review 1",
				"rating": 3,
				"createdAt": "` + timeNowString + `",
				"updatedAt": "` + timeNowString + `"
			},
			{
				"reviewId": 6,
				"postId":    "post1",
				"username":  "username2",
				"content": "a miña review 2",
				"rating": 3,
				"createdAt": "` + timeNowString + `",
				"updatedAt": "` + timeNowString + `"
			},
			{
				"reviewId": 7,
				"postId":    "post1",
				"username":  "username1",
				"content": "a miña review 3",
				"rating": 3,
				"createdAt": "` + timeNowString + `",
				"updatedAt": "` + timeNowString + `"
			}
			],
			"lastReviewId":7
		}
	}`

	controller.GetReviewsByPostId(ginContext)

	assert.Equal(t, apiResponse.Code, 200)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}

func TestInternalServerErrorOnGetReviewsByPostIdWithController_WhenServiceCallFails(t *testing.T) {
	setUpController(t)
	ginContext.Request, _ = http.NewRequest("GET", "/reviews", nil)
	expectedPostId := "post1"
	ginContext.Params = []gin.Param{{Key: "postId", Value: expectedPostId}}
	expectedError := errors.New("some error")
	controllerService.EXPECT().GetReviewsByPostId(expectedPostId, uint64(0), 12).Return([]*model.Review{}, uint64(0), expectedError)
	expectedBodyResponse := `{
		"error": true,
		"message": "` + expectedError.Error() + `",
		"content": null
	}`

	controller.GetReviewsByPostId(ginContext)

	assert.Equal(t, apiResponse.Code, 500)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}

func TestBadRequestErrorOnGetUserPostsWithController_WhenNoPostId(t *testing.T) {
	setUpController(t)
	ginContext.Request, _ = http.NewRequest("GET", "/reviews", nil)
	expectedError := "Missing parameter postId"
	expectedBodyResponse := `{
		"error": true,
		"message": "` + expectedError + `",
		"content":null
	}`

	controller.GetReviewsByPostId(ginContext)

	assert.Equal(t, apiResponse.Code, 400)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}

func TestBadRequestErrorOnGetUserPostsWithController_WhenLastReviewIdIsNotANumber(t *testing.T) {
	setUpController(t)
	ginContext.Request, _ = http.NewRequest("GET", "/reviews", nil)
	expectedPostId := "post1"
	ginContext.Params = []gin.Param{{Key: "postId", Value: expectedPostId}}
	wrongLastReviewId := "wrongReviewId"
	limit := 6
	u := url.Values{}
	u.Add("lastReviewId", wrongLastReviewId)
	u.Add("limit", strconv.Itoa(limit))
	ginContext.Request.URL.RawQuery = u.Encode()
	expectedError := "Invalid pagination parameters, lastReviewId must be a positive number"
	expectedBodyResponse := `{
		"error": true,
		"message": "` + expectedError + `",
		"content":null
	}`

	controller.GetReviewsByPostId(ginContext)

	assert.Equal(t, apiResponse.Code, 400)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}

func TestBadRequestErrorOnGetUserPostsWithController_WhenLimitSmallerThanOne(t *testing.T) {
	setUpController(t)
	ginContext.Request, _ = http.NewRequest("GET", "/reviews", nil)
	expectedPostId := "post1"
	ginContext.Params = []gin.Param{{Key: "postId", Value: expectedPostId}}
	lastReviewId := uint64(5)
	wrongLimit := 0
	u := url.Values{}
	u.Add("lastReviewId", strconv.FormatUint(lastReviewId, 10))
	u.Add("limit", strconv.Itoa(wrongLimit))
	ginContext.Request.URL.RawQuery = u.Encode()
	expectedError := "Invalid pagination parameters, limit must be greater than 0"
	expectedBodyResponse := `{
		"error": true,
		"message": "` + expectedError + `",
		"content":null
	}`

	controller.GetReviewsByPostId(ginContext)

	assert.Equal(t, apiResponse.Code, 400)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}
