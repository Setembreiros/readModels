package comment_test

import (
	"errors"
	"net/http"
	"net/url"
	"readmodels/internal/comment"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
)

var controller *comment.CommentController

func setUpController(t *testing.T) {
	SetUp(t)
	controller = comment.NewCommentController(repository)
}

func TestGetCommentsByPostIdWithController_WhenSuccess(t *testing.T) {
	setUpController(t)
	ginContext.Request, _ = http.NewRequest("GET", "/comments", nil)
	expectedPostId := "post1"
	expectedLastCommentId := uint64(4)
	expectedLimit := 4
	ginContext.Params = []gin.Param{{Key: "postId", Value: expectedPostId}}
	u := url.Values{}
	u.Add("lastCommentId", strconv.FormatUint(expectedLastCommentId, 10))
	u.Add("limit", strconv.Itoa(expectedLimit))
	ginContext.Request.URL.RawQuery = u.Encode()
	timeLayout := "2006-01-02T15:04:05.00Z"
	timeNowString := time.Now().UTC().Format(timeLayout)
	timeNow, _ := time.Parse(timeLayout, timeNowString)
	expectedComments := []*comment.Comment{
		{
			CommentId: uint64(5),
			Username:  "username1",
			PostId:    "post1",
			Content:   "o meu comentario 1",
			CreatedAt: timeNow,
		},
		{
			CommentId: uint64(6),
			Username:  "username2",
			PostId:    "post1",
			Content:   "o meu comentario 2",
			CreatedAt: timeNow,
		},
		{
			CommentId: uint64(7),
			Username:  "username1",
			PostId:    "post1",
			Content:   "o meu comentario 3",
			CreatedAt: timeNow,
		},
	}
	repository.EXPECT().GetCommentsByPostId(expectedPostId, expectedLastCommentId, expectedLimit).Return(expectedComments, uint64(7), nil)
	expectedBodyResponse := `{
		"error": false,
		"message": "200 OK",
		"content": {
			"comments":[	
			{
				"commentId": 5,
				"username":  "username1",
				"postId":    "post1",
				"content": "o meu comentario 1",
				"createdAt": "` + timeNowString + `"
			},
			{
				"commentId": 6,
				"username":  "username2",
				"postId":    "post1",
				"content": "o meu comentario 2",
				"createdAt": "` + timeNowString + `"
			},
			{
				"commentId": 7,
				"username":  "username1",
				"postId":    "post1",
				"content": "o meu comentario 3",
				"createdAt": "` + timeNowString + `"
			}
			],
			"lastCommentId":7
		}
	}`

	controller.GetCommentsByPostId(ginContext)

	assert.Equal(t, apiResponse.Code, 200)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}

func TestGetCommentsByPostIdWithController_WhenSuccessWithDefaultPaginationParameters(t *testing.T) {
	setUpController(t)
	ginContext.Request, _ = http.NewRequest("GET", "/comments", nil)
	expectedPostId := "post1"
	ginContext.Params = []gin.Param{{Key: "postId", Value: expectedPostId}}
	expectedDefaultLastCommentId := uint64(0)
	expectedDefaultLimit := 12
	timeLayout := "2006-01-02T15:04:05.00Z"
	timeNowString := time.Now().UTC().Format(timeLayout)
	timeNow, _ := time.Parse(timeLayout, timeNowString)
	expectedComments := []*comment.Comment{
		{
			CommentId: uint64(5),
			Username:  "username1",
			PostId:    "post1",
			Content:   "o meu comentario 1",
			CreatedAt: timeNow,
		},
		{
			CommentId: uint64(6),
			Username:  "username2",
			PostId:    "post1",
			Content:   "o meu comentario 2",
			CreatedAt: timeNow,
		},
		{
			CommentId: uint64(7),
			Username:  "username1",
			PostId:    "post1",
			Content:   "o meu comentario 3",
			CreatedAt: timeNow,
		},
	}
	repository.EXPECT().GetCommentsByPostId(expectedPostId, expectedDefaultLastCommentId, expectedDefaultLimit).Return(expectedComments, uint64(7), nil)
	expectedBodyResponse := `{
		"error": false,
		"message": "200 OK",
		"content": {
			"comments":[	
			{
				"commentId": 5,
				"username":  "username1",
				"postId":    "post1",
				"content": "o meu comentario 1",
				"createdAt": "` + timeNowString + `"
			},
			{
				"commentId": 6,
				"username":  "username2",
				"postId":    "post1",
				"content": "o meu comentario 2",
				"createdAt": "` + timeNowString + `"
			},
			{
				"commentId": 7,
				"username":  "username1",
				"postId":    "post1",
				"content": "o meu comentario 3",
				"createdAt": "` + timeNowString + `"
			}
			],
			"lastCommentId":7
		}
	}`

	controller.GetCommentsByPostId(ginContext)

	assert.Equal(t, apiResponse.Code, 200)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}

func TestInternalServerErrorOnGetCommentsByPostIdWithController_WhenServiceCallFails(t *testing.T) {
	setUpController(t)
	ginContext.Request, _ = http.NewRequest("GET", "/comments", nil)
	expectedPostId := "post1"
	ginContext.Params = []gin.Param{{Key: "postId", Value: expectedPostId}}
	expectedError := errors.New("some error")
	repository.EXPECT().GetCommentsByPostId(expectedPostId, uint64(0), 12).Return([]*comment.Comment{}, uint64(0), expectedError)
	expectedBodyResponse := `{
		"error": true,
		"message": "` + expectedError.Error() + `",
		"content": null
	}`

	controller.GetCommentsByPostId(ginContext)

	assert.Equal(t, apiResponse.Code, 500)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}

func TestBadRequestErrorOnGetUserPostsWithController_WhenNoPostId(t *testing.T) {
	setUpController(t)
	ginContext.Request, _ = http.NewRequest("GET", "/comments", nil)
	expectedError := "Missing parameter postId"
	expectedBodyResponse := `{
		"error": true,
		"message": "` + expectedError + `",
		"content":null
	}`

	controller.GetCommentsByPostId(ginContext)

	assert.Equal(t, apiResponse.Code, 400)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}

func TestBadRequestErrorOnGetUserPostsWithController_WhenLastCommentIdIsNotANumber(t *testing.T) {
	setUpController(t)
	ginContext.Request, _ = http.NewRequest("GET", "/comments", nil)
	expectedPostId := "post1"
	ginContext.Params = []gin.Param{{Key: "postId", Value: expectedPostId}}
	wrongLastCommentId := "wrongCommentId"
	limit := 6
	u := url.Values{}
	u.Add("lastCommentId", wrongLastCommentId)
	u.Add("limit", strconv.Itoa(limit))
	ginContext.Request.URL.RawQuery = u.Encode()
	expectedError := "Invalid pagination parameters, lastCommentId must be a positive number"
	expectedBodyResponse := `{
		"error": true,
		"message": "` + expectedError + `",
		"content":null
	}`

	controller.GetCommentsByPostId(ginContext)

	assert.Equal(t, apiResponse.Code, 400)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}

func TestBadRequestErrorOnGetUserPostsWithController_WhenLimitSmallerThanOne(t *testing.T) {
	setUpController(t)
	ginContext.Request, _ = http.NewRequest("GET", "/comments", nil)
	expectedPostId := "post1"
	ginContext.Params = []gin.Param{{Key: "postId", Value: expectedPostId}}
	lastCommentId := uint64(5)
	wrongLimit := 0
	u := url.Values{}
	u.Add("lastCommentId", strconv.FormatUint(lastCommentId, 10))
	u.Add("limit", strconv.Itoa(wrongLimit))
	ginContext.Request.URL.RawQuery = u.Encode()
	expectedError := "Invalid pagination parameters, limit must be greater than 0"
	expectedBodyResponse := `{
		"error": true,
		"message": "` + expectedError + `",
		"content":null
	}`

	controller.GetCommentsByPostId(ginContext)

	assert.Equal(t, apiResponse.Code, 400)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}
