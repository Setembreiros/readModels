package integration_test_comments

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"readmodels/internal/comment"
	comment_handler "readmodels/internal/comment/handler"
	database "readmodels/internal/db"
	"readmodels/internal/model"
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
var commentWasUpdatedEventHandler *comment_handler.CommentWasUpdatedEventHandler
var commentWasDeletedEventHandler *comment_handler.CommentWasDeletedEventHandler
var apiResponse *httptest.ResponseRecorder
var ginContext *gin.Context

func setUp(t *testing.T) {
	apiResponse = httptest.NewRecorder()
	ginContext, _ = gin.CreateTestContext(apiResponse)

	// Real infrastructure and services
	db = integration_test_arrange.CreateTestDatabase(t, ginContext)
	repository := comment.NewCommentRepository(db)
	service := comment.NewCommentService(repository)
	controller = comment.NewCommentController(repository)
	commentWasCreatedEventHandler = comment_handler.NewCommentWasCreatedEventHandler(service)
	commentWasUpdatedEventHandler = comment_handler.NewCommentWasUpdatedEventHandler(repository)
	commentWasDeletedEventHandler = comment_handler.NewCommentWasDeletedEventHandler(service)
}

func tearDown() {
	db.Client.Truncate()
}

func TestCreateNewComment_WhenDatabaseReturnsSuccess(t *testing.T) {
	setUp(t)
	defer tearDown()
	existingPost := &database.PostMetadata{
		PostId:   "post123",
		Username: "username1",
		Type:     "TEXT",
	}
	integration_test_arrange.AddPostToDatabase(t, db, existingPost)
	timeNow := time.Now().UTC().Format(model.TimeLayout)
	data := &comment_handler.CommentWasCreatedEvent{
		CommentId: uint64(123456),
		Username:  "user123",
		PostId:    "post123",
		Content:   "Exemplo de content",
		CreatedAt: timeNow,
	}
	event, _ := test_common.SerializeData(data)
	expectedTime, _ := time.Parse(model.TimeLayout, data.CreatedAt)
	expectedComment := &model.Comment{
		CommentId: data.CommentId,
		Username:  data.Username,
		PostId:    data.PostId,
		Content:   data.Content,
		CreatedAt: expectedTime,
	}

	commentWasCreatedEventHandler.Handle(event)

	integration_test_assert.AssertCommentExists(t, db, data.CommentId, expectedComment)
	integration_test_assert.AssertPostCommentsIncreased(t, db, existingPost.PostId)
}

func TestGetCommentsByPostId_WhenDatabaseReturnsSuccess(t *testing.T) {
	setUp(t)
	defer tearDown()
	timeNowString := time.Now().UTC().Format(model.TimeLayout)
	timeNow, _ := time.Parse(model.TimeLayout, timeNowString)
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
				"postId":    "post1",
				"username":  "username1",
				"content": 	 "o meu comentario 3",
				"createdAt": "` + timeNowString + `",
				"updatedAt": "` + timeNowString + `"
			},	
			{
				"commentId": 6,
				"postId":    "post1",
				"username":  "username2",
				"content": 	 "o meu comentario 6",
				"createdAt": "` + timeNowString + `",
				"updatedAt": "` + timeNowString + `"
			},	
			{
				"commentId": 8,
				"postId":    "post1",
				"username":  "username3",
				"content": 	 "o meu comentario 8",
				"createdAt": "` + timeNowString + `",
				"updatedAt": "` + timeNowString + `"
			},
			{
				"commentId": 9,
				"postId":    "post1",
				"username":  "username1",
				"content": 	 "o meu comentario 9",
				"createdAt": "` + timeNowString + `",
				"updatedAt": "` + timeNowString + `"
			}
			],
			"lastCommentId":9
		}
	}`

	controller.GetCommentsByPostId(ginContext)

	integration_test_assert.AssertSuccessResult(t, apiResponse, expectedBodyResponse)
}

func TestUpdateComment_WhenDatabaseReturnsSuccess(t *testing.T) {
	setUp(t)
	defer tearDown()
	existingComment := &model.Comment{
		CommentId: uint64(1234),
		Username:  "user123",
		PostId:    "post123",
		Content:   "Exemplo de content",
		CreatedAt: time.Now(),
	}
	integration_test_arrange.AddCommentToDatabase(t, db, existingComment)
	timeNow := time.Now().UTC().Format(model.TimeLayout)
	data := &comment_handler.CommentWasUpdatedEvent{
		CommentId: existingComment.CommentId,
		Content:   "Exemplo de content actualizado",
		UpdatedAt: timeNow,
	}
	event, _ := test_common.SerializeData(data)
	expectedTime, _ := time.Parse(model.TimeLayout, data.UpdatedAt)
	expectedComment := &model.Comment{
		CommentId: data.CommentId,
		Username:  existingComment.Username,
		PostId:    existingComment.PostId,
		Content:   data.Content,
		CreatedAt: existingComment.CreatedAt,
		UpdatedAt: expectedTime,
	}

	commentWasUpdatedEventHandler.Handle(event)

	integration_test_assert.AssertCommentExists(t, db, data.CommentId, expectedComment)
}

func TestDeleteComment_WhenDatabaseReturnsSuccess(t *testing.T) {
	setUp(t)
	defer tearDown()
	existingComment := &model.Comment{
		CommentId: uint64(1234),
		Username:  "user123",
		PostId:    "post123",
		Content:   "Exemplo de content",
	}
	integration_test_arrange.AddCommentToDatabase(t, db, existingComment)
	data := &comment_handler.CommentWasDeletedEvent{
		CommentId: existingComment.CommentId,
		PostId:    existingComment.PostId,
	}
	event, _ := test_common.SerializeData(data)

	commentWasDeletedEventHandler.Handle(event)

	integration_test_assert.AssertCommentDoesNotExist(t, db, data.CommentId)
	integration_test_assert.AssertPostCommentsDecreased(t, db, existingComment.PostId)
}

func populateDb(t *testing.T, time time.Time) {
	existingComments := []*model.Comment{
		{
			CommentId: uint64(1),
			Username:  "username1",
			PostId:    "post1",
			Content:   "o meu comentario 1",
			CreatedAt: time,
			UpdatedAt: time,
		},
		{
			CommentId: uint64(2),
			Username:  "username2",
			PostId:    "post1",
			Content:   "o meu comentario 2",
			CreatedAt: time,
			UpdatedAt: time,
		},
		{
			CommentId: uint64(3),
			Username:  "username1",
			PostId:    "post1",
			Content:   "o meu comentario 3",
			CreatedAt: time,
			UpdatedAt: time,
		},
		{
			CommentId: uint64(4),
			Username:  "user123",
			PostId:    "post2",
			Content:   "o meu comentario 4",
			CreatedAt: time,
			UpdatedAt: time,
		},
		{
			CommentId: uint64(5),
			Username:  "user123",
			PostId:    "post2",
			Content:   "o meu comentario 5",
			CreatedAt: time,
			UpdatedAt: time,
		},
		{
			CommentId: uint64(6),
			Username:  "username2",
			PostId:    "post1",
			Content:   "o meu comentario 6",
			CreatedAt: time,
			UpdatedAt: time,
		},
		{
			CommentId: uint64(7),
			Username:  "username1",
			PostId:    "post2",
			Content:   "o meu comentario 7",
			CreatedAt: time,
			UpdatedAt: time,
		},
		{
			CommentId: uint64(8),
			Username:  "username3",
			PostId:    "post1",
			Content:   "o meu comentario 8",
			CreatedAt: time,
			UpdatedAt: time,
		},
		{
			CommentId: uint64(9),
			Username:  "username1",
			PostId:    "post1",
			Content:   "o meu comentario 9",
			CreatedAt: time,
			UpdatedAt: time,
		},
		{
			CommentId: uint64(10),
			Username:  "user123",
			PostId:    "post2",
			Content:   "o meu comentario 10",
			CreatedAt: time,
			UpdatedAt: time,
		},
		{
			CommentId: uint64(11),
			Username:  "user123",
			PostId:    "post1",
			Content:   "o meu comentario 11",
			CreatedAt: time,
			UpdatedAt: time,
		},
	}

	for _, existingComment := range existingComments {
		integration_test_arrange.AddCommentToDatabase(t, db, existingComment)
	}
}
