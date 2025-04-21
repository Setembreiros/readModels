package comment_test

import (
	"encoding/json"
	comment_handler "readmodels/internal/comment/handler"
	mock_comment_handler "readmodels/internal/comment/handler/test/mock"
	"testing"

	"github.com/stretchr/testify/assert"
)

var commentWasDeletedEventHandler *comment_handler.CommentWasDeletedEventHandler
var commentWasDeletedEventService *mock_comment_handler.MockCommentWasDeletedEventService

func setUpCommentWasDeletedEventHandler(t *testing.T) {
	SetUp(t)
	commentWasDeletedEventService = mock_comment_handler.NewMockCommentWasDeletedEventService(ctrl)
	commentWasDeletedEventHandler = comment_handler.NewCommentWasDeletedEventHandler(commentWasDeletedEventService)
}

func TestHandleCommentWasDeletedEvent(t *testing.T) {
	setUpCommentWasDeletedEventHandler(t)
	data := &comment_handler.CommentWasDeletedEvent{
		CommentId: uint64(123456),
		PostId:    "post1",
	}
	event, _ := json.Marshal(data)
	commentWasDeletedEventService.EXPECT().DeleteComment(data.PostId, data.CommentId)

	commentWasDeletedEventHandler.Handle(event)
}

func TestInvalidDataInCommentWasDeletedEventHandler(t *testing.T) {
	setUpHandler(t)
	invalidData := "invalid data"
	event, _ := json.Marshal(invalidData)

	commentWasDeletedEventHandler.Handle(event)

	assert.Contains(t, loggerOutput.String(), "Invalid event data")
}
