package post_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"readmodels/internal/model"
	"readmodels/internal/post"
	mock_post "readmodels/internal/post/mock"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog/log"
)

var controllerLoggerOutput bytes.Buffer
var controllerService *mock_post.MockService
var controller *post.PostController
var apiResponse *httptest.ResponseRecorder
var ginContext *gin.Context

func setUpHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	controllerService = mock_post.NewMockService(ctrl)
	log.Logger = log.Output(&controllerLoggerOutput)
	controller = post.NewPostController(controllerService)
	gin.SetMode(gin.TestMode)
	apiResponse = httptest.NewRecorder()
	ginContext, _ = gin.CreateTestContext(apiResponse)
}

func TestGetPostMetadatasByUser(t *testing.T) {
	setUpHandler(t)
	username := "username1"
	currentUsername := "username1"
	lastPostId := "post4"
	lastPostCreatedAt := "0001-01-03T00:00:00Z"
	limit := "4"
	ginContext.Request = &http.Request{
		Header: make(http.Header),
		URL:    &url.URL{},
	}
	ginContext.Request.Method = "GET"
	ginContext.Request.Header.Set("Content-Type", "application/json")
	ginContext.Params = []gin.Param{{Key: "username", Value: username}, {Key: "currentUsername", Value: currentUsername}}
	u := url.Values{}
	u.Add("lastPostId", lastPostId)
	u.Add("lastPostCreatedAt", lastPostCreatedAt)
	u.Add("limit", limit)
	ginContext.Request.URL.RawQuery = u.Encode()
	timeNow, _ := time.Parse(model.TimeLayout, time.Now().UTC().Format(model.TimeLayout))
	data := []*post.PostMetadata{
		{
			PostId:                    "123456",
			Username:                  username,
			Type:                      "TEXT",
			Title:                     "Exemplo de Título",
			Description:               "Exemplo de Descrición",
			Reviews:                   1,
			Comments:                  1,
			Likes:                     2,
			IsLikedByCurrentUser:      true,
			Superlikes:                3,
			IsSuperlikedByCurrentUser: true,
			CreatedAt:                 timeNow,
			LastUpdated:               timeNow,
		},
		{
			PostId:                    "abcdef",
			Username:                  username,
			Type:                      "IMAGE",
			Title:                     "Exemplo de Título 2",
			Description:               "Exemplo de Descrición 2",
			Reviews:                   1,
			Comments:                  1,
			Likes:                     2,
			IsLikedByCurrentUser:      true,
			Superlikes:                3,
			IsSuperlikedByCurrentUser: true,
			CreatedAt:                 timeNow,
			LastUpdated:               timeNow,
		},
	}
	controllerService.EXPECT().GetPostMetadatasByUser(username, currentUsername, lastPostId, lastPostCreatedAt, 4).Return(data, "post7", "0001-01-06T00:00:00Z", nil)
	expectedBodyResponse := `{
		"error": false,
		"message": "200 OK",
		"content": {"posts":[
		{
			"post_id":      "123456",
			"username":   "username1",
			"type":        "TEXT",
			"title":       "Exemplo de Título",
			"description": "Exemplo de Descrición",
			"reviews": 1,
			"comments": 1,
			"likes": 2,
			"isLikedByCurrentUser": true,
			"superlikes": 3,
			"isSuperlikedByCurrentUser": true,
			"created_at":   "` + timeNow.Format(model.TimeLayout) + `",
			"last_updated": "` + timeNow.Format(model.TimeLayout) + `"
		},
		{
			"post_id":      "abcdef",
			"username":    "username1",
			"type":        "IMAGE",
			"title":       "Exemplo de Título 2",
			"description": "Exemplo de Descrición 2",
			"reviews": 1,
			"comments": 1,
			"likes": 2,
			"isLikedByCurrentUser": true,
			"superlikes": 3,
			"isSuperlikedByCurrentUser": true,
			"created_at":   "` + timeNow.Format(model.TimeLayout) + `",
			"last_updated": "` + timeNow.Format(model.TimeLayout) + `"
		}
		],"limit":4,"lastPostId":"post7","lastPostCreatedAt":"0001-01-06T00:00:00Z"}
	}`

	controller.GetPostMetadatasByUser(ginContext)

	assert.Equal(t, apiResponse.Code, 200)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}

func TestGetPostMetadatasByUserWithDefaultPaginationParameters(t *testing.T) {
	setUpHandler(t)
	username := "username1"
	currentUsername := "username2"
	ginContext.Params = []gin.Param{{Key: "username", Value: username}, {Key: "currentUsername", Value: currentUsername}}
	timeNow, _ := time.Parse(model.TimeLayout, time.Now().UTC().Format(model.TimeLayout))
	data := []*post.PostMetadata{
		{
			PostId:      "123456",
			Username:    username,
			Type:        "TEXT",
			Title:       "Exemplo de Título",
			Description: "Exemplo de Descrición",
			CreatedAt:   timeNow,
			LastUpdated: timeNow,
		},
		{
			PostId:      "abcdef",
			Username:    username,
			Type:        "IMAGE",
			Title:       "Exemplo de Título 2",
			Description: "Exemplo de Descrición 2",
			CreatedAt:   timeNow,
			LastUpdated: timeNow,
		},
	}
	expectedDefaultLastPostId := ""
	expectedDefaultLastPostCreatedAt := ""
	expectedDefaultLimit := 6
	controllerService.EXPECT().GetPostMetadatasByUser(username, currentUsername, expectedDefaultLastPostId, expectedDefaultLastPostCreatedAt, expectedDefaultLimit).Return(data, "post7", "0001-01-06T00:00:00Z", nil)
	expectedBodyResponse := `{
		"error": false,
		"message": "200 OK",
		"content": {"posts":[
		{
			"post_id":      "123456",
			"username":   "username1",
			"type":        "TEXT",
			"title":       "Exemplo de Título",
			"description": "Exemplo de Descrición",
			"reviews": 0,
			"comments": 0,
			"likes": 0,
			"isLikedByCurrentUser": false,
			"superlikes": 0,
			"isSuperlikedByCurrentUser": false,
			"created_at":   "` + timeNow.Format(model.TimeLayout) + `",
			"last_updated": "` + timeNow.Format(model.TimeLayout) + `"
		},
		{
			"post_id":      "abcdef",
			"username":    "username1",
			"type":        "IMAGE",
			"title":       "Exemplo de Título 2",
			"description": "Exemplo de Descrición 2",
			"reviews": 0,
			"comments": 0,
			"likes": 0,
			"isLikedByCurrentUser": false,
			"superlikes": 0,
			"isSuperlikedByCurrentUser": false,
			"created_at":   "` + timeNow.Format(model.TimeLayout) + `",
			"last_updated": "` + timeNow.Format(model.TimeLayout) + `"
		}
		],"limit":6,"lastPostId":"post7","lastPostCreatedAt":"0001-01-06T00:00:00Z"}
	}`

	controller.GetPostMetadatasByUser(ginContext)

	assert.Equal(t, apiResponse.Code, 200)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}

func TestInternalServerErrorOnGetPostMetadatasByUser(t *testing.T) {
	setUpHandler(t)
	username := "username1"
	currentUsername := "username1"
	lastPostId := "post4"
	lastPostCreatedAt := "0001-01-03T00:00:00Z"
	limit := "4"
	ginContext.Request = &http.Request{
		Header: make(http.Header),
		URL:    &url.URL{},
	}
	ginContext.Request.Method = "GET"
	ginContext.Request.Header.Set("Content-Type", "application/json")
	ginContext.Params = []gin.Param{{Key: "username", Value: username}, {Key: "currentUsername", Value: currentUsername}}
	u := url.Values{}
	u.Add("lastPostId", lastPostId)
	u.Add("lastPostCreatedAt", lastPostCreatedAt)
	u.Add("limit", limit)
	ginContext.Request.URL.RawQuery = u.Encode()
	expectedData := []*post.PostMetadata{}
	expectedError := errors.New("some error")
	controllerService.EXPECT().GetPostMetadatasByUser(username, currentUsername, lastPostId, lastPostCreatedAt, 4).Return(expectedData, "", "", expectedError)
	expectedBodyResponse := `{
		"error": true,
		"message": "` + expectedError.Error() + `",
		"content":null
	}`

	controller.GetPostMetadatasByUser(ginContext)

	assert.Equal(t, apiResponse.Code, 500)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}

func TestBadRequestErrorOnGetUserPostsWhenLimitSmallerThanZero(t *testing.T) {
	setUpHandler(t)
	username := "username1"
	currentUsername := "username1"
	lastPostId := "post4"
	lastPostCreatedAt := "0001-01-03T00:00:00Z"
	wrongLimit := "0"
	ginContext.Request = &http.Request{
		Header: make(http.Header),
		URL:    &url.URL{},
	}
	ginContext.Request.Method = "GET"
	ginContext.Request.Header.Set("Content-Type", "application/json")
	ginContext.Params = []gin.Param{{Key: "username", Value: username}, {Key: "currentUsername", Value: currentUsername}}
	u := url.Values{}
	u.Add("lastPostId", lastPostId)
	u.Add("lastPostCreatedAt", lastPostCreatedAt)
	u.Add("limit", wrongLimit)
	ginContext.Request.URL.RawQuery = u.Encode()
	expectedError := "Invalid pagination parameters, limit has to be greater than 0"
	expectedBodyResponse := `{
		"error": true,
		"message": "` + expectedError + `",
		"content":null
	}`

	controller.GetPostMetadatasByUser(ginContext)

	assert.Equal(t, apiResponse.Code, 400)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}

func TestBadRequestErrorOnGetUserPostsWhenLastPostIdIsNotEmptyButLastPostCreatedAtIsEmpty(t *testing.T) {
	setUpHandler(t)
	username := "username1"
	currentUsername := "username1"
	lastPostId := "post4"
	lastPostCreatedAt := ""
	wrongLimit := "2"
	ginContext.Request = &http.Request{
		Header: make(http.Header),
		URL:    &url.URL{},
	}
	ginContext.Request.Method = "GET"
	ginContext.Request.Header.Set("Content-Type", "application/json")
	ginContext.Params = []gin.Param{{Key: "username", Value: username}, {Key: "currentUsername", Value: currentUsername}}
	u := url.Values{}
	u.Add("lastPostId", lastPostId)
	u.Add("lastPostCreatedAt", lastPostCreatedAt)
	u.Add("limit", wrongLimit)
	ginContext.Request.URL.RawQuery = u.Encode()
	expectedError := "Invalid pagination parameters, lastPostId and lastPostCreatedAt both have to have value or both have to be empty"
	expectedBodyResponse := `{
		"error": true,
		"message": "` + expectedError + `",
		"content":null
	}`

	controller.GetPostMetadatasByUser(ginContext)

	assert.Equal(t, apiResponse.Code, 400)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}

func removeSpace(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(s, " ", ""), "\t", ""), "\n", "")
}
