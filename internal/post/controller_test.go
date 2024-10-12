package post_test

import (
	"bytes"
	"errors"
	"net/http/httptest"
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
	controllerRepository.EXPECT().GetPostMetadatasByUser(username).Return(data, nil)
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
			"created_at":   "` + timeNow.Format("2006-01-02T15:04:05.0000000Z") + `",
			"last_updated": "` + timeNow.Format("2006-01-02T15:04:05.0000000Z") + `"
		},
		{
			"post_id":      "abcdef",
			"username":    "username1",
			"type":        "IMAGE",
			"title":       "Exemplo de Título 2",
			"description": "Exemplo de Descrição 2",
			"created_at":   "` + timeNow.Format("2006-01-02T15:04:05.0000000Z") + `",
			"last_updated": "` + timeNow.Format("2006-01-02T15:04:05.0000000Z") + `"
		}
		]}
	}`

	controller.GetPostMetadatasByUser(ginContext)

	assert.Equal(t, apiResponse.Code, 200)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}

func TestInternalServerErrorOnGetPostMetadatasByUser(t *testing.T) {
	setUpHandler(t)
	username := "username1"
	ginContext.Params = []gin.Param{{Key: "username", Value: username}}
	expectedData := []*post.PostMetadata{}
	expectedError := errors.New("some error")
	controllerRepository.EXPECT().GetPostMetadatasByUser(username).Return(expectedData, expectedError)
	expectedBodyResponse := `{
		"error": true,
		"message": "` + expectedError.Error() + `",
		"content":null
	}`

	controller.GetPostMetadatasByUser(ginContext)

	assert.Equal(t, apiResponse.Code, 500)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}

func removeSpace(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(s, " ", ""), "\t", ""), "\n", "")
}
