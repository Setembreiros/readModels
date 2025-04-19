package integration_test_comments

import (
	"context"
	"readmodels/internal/comment"
	comment_handler "readmodels/internal/comment/handler"
	database "readmodels/internal/db"
	integration_test_arrange "readmodels/test/integration_test_common/arrange"
	integration_test_assert "readmodels/test/integration_test_common/assert"
	"readmodels/test/test_common"
	"testing"
	"time"
)

var db *database.Database
var handler *comment_handler.CommentWasCreatedEventHandler

func setUp(t *testing.T) {
	// Real infrastructure and services
	ctx := context.TODO()
	db = integration_test_arrange.CreateTestDatabase(t, ctx)
	repository := comment.CommentRepository(*db)
	handler = comment_handler.NewCommentWasCreatedEventHandler(repository)
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

	handler.Handle(event)

	integration_test_assert.AssertCommentExists(t, db, data.CommentId, expectedComment)
}
