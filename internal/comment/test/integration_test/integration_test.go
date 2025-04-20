package integration_test_comments

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"readmodels/internal/comment"
	comment_handler "readmodels/internal/comment/handler"
	database "readmodels/internal/db"
	integration_test_arrange "readmodels/test/integration_test_common/arrange"
	integration_test_assert "readmodels/test/integration_test_common/assert"
	"readmodels/test/test_common"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

var db *database.Database
var controller *comment.CommentController
var commentWasCreatedEventHandler *comment_handler.CommentWasCreatedEventHandler
var commentWasDeletedEventHandler *comment_handler.CommentWasDeletedEventHandler
var apiResponse *httptest.ResponseRecorder
var ginContext *gin.Context

func setUp(t *testing.T) {
	apiResponse = httptest.NewRecorder()
	ginContext, _ = gin.CreateTestContext(apiResponse)

	// Real infrastructure and services
	db = integration_test_arrange.CreateTestDatabase(t, ginContext)
	repository := comment.CommentRepository(*db)
	controller = comment.NewCommentController(repository)
	commentWasCreatedEventHandler = comment_handler.NewCommentWasCreatedEventHandler(repository)
	commentWasDeletedEventHandler = comment_handler.NewCommentWasDeletedEventHandler(repository)
}

func tearDown() {
	db.Client.Clean()
}

func TestCreateNewComment_WhenDatabaseReturnsSuccess(t *testing.T) {
	setUp(t)
	defer tearDown()
	timeLayout := "2006-01-02T15:04:05.000000000Z"
	timeNow := time.Now().UTC().Format(timeLayout)
	data := &comment_handler.CommentWasCreatedEvent{
		CommentId: uint64(123456),
		Username:  "user123",
		PostId:    "post123",
		Content:   "Exemplo de content",
		CreatedAt: timeNow,
	}
	event, _ := test_common.SerializeData(data)
	expectedTime, _ := time.Parse(timeLayout, data.CreatedAt)
	expectedComment := &comment.Comment{
		CommentId: data.CommentId,
		Username:  data.Username,
		PostId:    data.PostId,
		Content:   data.Content,
		CreatedAt: expectedTime,
	}

	commentWasCreatedEventHandler.Handle(event)

	integration_test_assert.AssertCommentExists(t, db, data.CommentId, expectedComment)
}

func TestGetCommentsByPostId_WhenDatabaseReturnsSuccess(t *testing.T) {
	setUp(t)
	defer tearDown()
	timeLayout := "2006-01-02T15:04:05.00Z"
	timeNowString := time.Now().UTC().Format(timeLayout)
	timeNow, _ := time.Parse(timeLayout, timeNowString)
	populateDb(t, timeNow)
	postId := "post1"
	lastCommentId := uint64(2)
	limit := 4
	ginContext.Request, _ = http.NewRequest("GET", "/comments", nil)
	ginContext.Params = []gin.Param{{Key: "postId", Value: postId}}
	u := url.Values{}
	u.Add("lastCommentId", strconv.FormatUint(lastCommentId, 10))
	u.Add("limit", strconv.Itoa(limit))
	ginContext.Request.URL.RawQuery = u.Encode()
	expectedBodyResponse := `{
		"error": false,
		"message": "200 OK",
		"content": {
			"comments":[	
			{
				"commentId": 3,
				"username":  "username1",
				"postId":    "post1",
				"content": 	 "o meu comentario 3",
				"createdAt": "` + timeNowString + `"
			},	
			{
				"commentId": 6,
				"username":  "username2",
				"postId":    "post1",
				"content": 	 "o meu comentario 6",
				"createdAt": "` + timeNowString + `"
			},	
			{
				"commentId": 8,
				"username":  "username3",
				"postId":    "post1",
				"content": 	 "o meu comentario 8",
				"createdAt": "` + timeNowString + `"
			},
			{
				"commentId": 9,
				"username":  "username1",
				"postId":    "post1",
				"content": 	 "o meu comentario 9",
				"createdAt": "` + timeNowString + `"
			}
			],
			"lastCommentId":9
		}
	}`

	controller.GetCommentsByPostId(ginContext)

	integration_test_assert.AssertSuccessResult(t, apiResponse, expectedBodyResponse)
}

func TestDeleteComment_WhenDatabaseReturnsSuccess(t *testing.T) {
	setUp(t)
	defer tearDown()
	existingComment := &comment.Comment{
		CommentId: uint64(1234),
		Username:  "user123",
		PostId:    "post123",
		Content:   "Exemplo de content",
	}
	integration_test_arrange.AddCommentToDatabase(t, db, existingComment)
	data := &comment_handler.CommentWasDeletedEvent{
		CommentId: existingComment.CommentId,
	}
	event, _ := test_common.SerializeData(data)

	commentWasDeletedEventHandler.Handle(event)

	integration_test_assert.AssertCommentDoesNotExist(t, db, data.CommentId)
}

func populateDb(t *testing.T, time time.Time) {
	existingComments := []*comment.Comment{
		{
			CommentId: uint64(1),
			Username:  "username1",
			PostId:    "post1",
			Content:   "o meu comentario 1",
			CreatedAt: time,
		},
		{
			CommentId: uint64(2),
			Username:  "username2",
			PostId:    "post1",
			Content:   "o meu comentario 2",
			CreatedAt: time,
		},
		{
			CommentId: uint64(3),
			Username:  "username1",
			PostId:    "post1",
			Content:   "o meu comentario 3",
			CreatedAt: time,
		},
		{
			CommentId: uint64(4),
			Username:  "user123",
			PostId:    "post2",
			Content:   "o meu comentario 4",
			CreatedAt: time,
		},
		{
			CommentId: uint64(5),
			Username:  "user123",
			PostId:    "post2",
			Content:   "o meu comentario 5",
			CreatedAt: time,
		},
		{
			CommentId: uint64(6),
			Username:  "username2",
			PostId:    "post1",
			Content:   "o meu comentario 6",
			CreatedAt: time,
		},
		{
			CommentId: uint64(7),
			Username:  "username1",
			PostId:    "post2",
			Content:   "o meu comentario 7",
			CreatedAt: time,
		},
		{
			CommentId: uint64(8),
			Username:  "username3",
			PostId:    "post1",
			Content:   "o meu comentario 8",
			CreatedAt: time,
		},
		{
			CommentId: uint64(9),
			Username:  "username1",
			PostId:    "post1",
			Content:   "o meu comentario 9",
			CreatedAt: time,
		},
		{
			CommentId: uint64(10),
			Username:  "user123",
			PostId:    "post2",
			Content:   "o meu comentario 10",
			CreatedAt: time,
		},
		{
			CommentId: uint64(11),
			Username:  "user123",
			PostId:    "post1",
			Content:   "o meu comentario 11",
			CreatedAt: time,
		},
	}

	for _, existingComment := range existingComments {
		integration_test_arrange.AddCommentToDatabase(t, db, existingComment)
	}
}
