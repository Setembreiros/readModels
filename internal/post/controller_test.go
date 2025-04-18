package post_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
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
var controllerRepository *mock_post.MockRepository
var controller *post.PostController
var apiResponse *httptest.ResponseRecorder
var ginContext *gin.Context
var timeLayout string = "2006-01-02T15:04:05.0000000Z"

func setUpHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	controllerRepository = mock_post.NewMockRepository(ctrl)
	log.Logger = log.Output(&controllerLoggerOutput)
	controller = post.NewPostController(controllerRepository)
	gin.SetMode(gin.TestMode)
	apiResponse = httptest.NewRecorder()
	ginContext, _ = gin.CreateTestContext(apiResponse)
}

func TestGetPostMetadatasByUser(t *testing.T) {
	setUpHandler(t)
	username := "username1"
	lastPostId := "post4"
	lastPostCreatedAt := "0001-01-03T00:00:00Z"
	limit := "4"
	ginContext.Request = &http.Request{
		Header: make(http.Header),
		URL:    &url.URL{},
	}
	ginContext.Request.Method = "GET"
	ginContext.Request.Header.Set("Content-Type", "application/json")
	ginContext.Params = []gin.Param{{Key: "username", Value: username}}
	u := url.Values{}
	u.Add("lastPostId", lastPostId)
	u.Add("lastPostCreatedAt", lastPostCreatedAt)
	u.Add("limit", limit)
	ginContext.Request.URL.RawQuery = u.Encode()
	timeNow := time.Now().UTC()
	data := []*post.PostMetadata{
		{
			PostId:      "123456",
			Username:    username,
			Type:        "TEXT",
			Title:       "Exemplo de Título",
			Description: "Exemplo de Descrição",
			CreatedAt:   timeNow,
			LastUpdated: timeNow,
		},
		{
			PostId:      "abcdef",
			Username:    username,
			Type:        "IMAGE",
			Title:       "Exemplo de Título 2",
			Description: "Exemplo de Descrição 2",
			CreatedAt:   timeNow,
			LastUpdated: timeNow,
		},
	}
	controllerRepository.EXPECT().GetPostMetadatasByUser(username, lastPostId, lastPostCreatedAt, 4).Return(data, "post7", "0001-01-06T00:00:00Z", nil)
	expectedBodyResponse := `{
		"error": false,
		"message": "200 OK",
		"content": {"posts":[
		{
			"post_id":      "123456",
			"username":   "username1",
			"type":        "TEXT",
			"title":       "Exemplo de Título",
			"description": "Exemplo de Descrição",
			"created_at":   "` + timeNow.Format(timeLayout) + `",
			"last_updated": "` + timeNow.Format(timeLayout) + `"
		},
		{
			"post_id":      "abcdef",
			"username":    "username1",
			"type":        "IMAGE",
			"title":       "Exemplo de Título 2",
			"description": "Exemplo de Descrição 2",
			"created_at":   "` + timeNow.Format(timeLayout) + `",
			"last_updated": "` + timeNow.Format(timeLayout) + `"
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
	ginContext.Params = []gin.Param{{Key: "username", Value: username}}
	timeNow := time.Now().UTC()
	data := []*post.PostMetadata{
		{
			PostId:      "123456",
			Username:    username,
			Type:        "TEXT",
			Title:       "Exemplo de Título",
			Description: "Exemplo de Descrição",
			CreatedAt:   timeNow,
			LastUpdated: timeNow,
		},
		{
			PostId:      "abcdef",
			Username:    username,
			Type:        "IMAGE",
			Title:       "Exemplo de Título 2",
			Description: "Exemplo de Descrição 2",
			CreatedAt:   timeNow,
			LastUpdated: timeNow,
		},
	}
	expectedDefaultLastPostId := ""
	expectedDefaultLastPostCreatedAt := ""
	expectedDefaultLimit := 6
	controllerRepository.EXPECT().GetPostMetadatasByUser(username, expectedDefaultLastPostId, expectedDefaultLastPostCreatedAt, expectedDefaultLimit).Return(data, "post7", "0001-01-06T00:00:00Z", nil)
	expectedBodyResponse := `{
		"error": false,
		"message": "200 OK",
		"content": {"posts":[
		{
			"post_id":      "123456",
			"username":   "username1",
			"type":        "TEXT",
			"title":       "Exemplo de Título",
			"description": "Exemplo de Descrição",
			"created_at":   "` + timeNow.Format(timeLayout) + `",
			"last_updated": "` + timeNow.Format(timeLayout) + `"
		},
		{
			"post_id":      "abcdef",
			"username":    "username1",
			"type":        "IMAGE",
			"title":       "Exemplo de Título 2",
			"description": "Exemplo de Descrição 2",
			"created_at":   "` + timeNow.Format(timeLayout) + `",
			"last_updated": "` + timeNow.Format(timeLayout) + `"
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
	lastPostId := "post4"
	lastPostCreatedAt := "0001-01-03T00:00:00Z"
	limit := "4"
	ginContext.Request = &http.Request{
		Header: make(http.Header),
		URL:    &url.URL{},
	}
	ginContext.Request.Method = "GET"
	ginContext.Request.Header.Set("Content-Type", "application/json")
	ginContext.Params = []gin.Param{{Key: "username", Value: username}}
	u := url.Values{}
	u.Add("lastPostId", lastPostId)
	u.Add("lastPostCreatedAt", lastPostCreatedAt)
	u.Add("limit", limit)
	ginContext.Request.URL.RawQuery = u.Encode()
	expectedData := []*post.PostMetadata{}
	expectedError := errors.New("some error")
	controllerRepository.EXPECT().GetPostMetadatasByUser(username, lastPostId, lastPostCreatedAt, 4).Return(expectedData, "", "", expectedError)
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
	lastPostId := "post4"
	lastPostCreatedAt := "0001-01-03T00:00:00Z"
	wrongLimit := "0"
	ginContext.Request = &http.Request{
		Header: make(http.Header),
		URL:    &url.URL{},
	}
	ginContext.Request.Method = "GET"
	ginContext.Request.Header.Set("Content-Type", "application/json")
	ginContext.Params = []gin.Param{{Key: "username", Value: username}}
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
	lastPostId := "post4"
	lastPostCreatedAt := ""
	wrongLimit := "2"
	ginContext.Request = &http.Request{
		Header: make(http.Header),
		URL:    &url.URL{},
	}
	ginContext.Request.Method = "GET"
	ginContext.Request.Header.Set("Content-Type", "application/json")
	ginContext.Params = []gin.Param{{Key: "username", Value: username}}
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
