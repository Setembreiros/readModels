package comment_test

import (
	"encoding/json"
	comment_handler "readmodels/internal/comment/handler"
	"testing"

	"github.com/stretchr/testify/assert"
)

var commentWasDeletedEventHandler *comment_handler.CommentWasDeletedEventHandler

func setUpCommentWasDeletedEventHandler(t *testing.T) {
	SetUp(t)
	commentWasDeletedEventHandler = comment_handler.NewCommentWasDeletedEventHandler(repository)
}

func TestHandleCommentWasDeletedEvent(t *testing.T) {
	setUpCommentWasDeletedEventHandler(t)
	data := &comment_handler.CommentWasDeletedEvent{
		CommentId: uint64(123456),
	}
	event, _ := json.Marshal(data)
	repository.EXPECT().DeleteComment(data.CommentId)

	commentWasDeletedEventHandler.Handle(event)
}

func TestInvalidDataInCommentWasDeletedEventHandler(t *testing.T) {
	setUpHandler(t)
	invalidData := "invalid data"
	event, _ := json.Marshal(invalidData)

	commentWasDeletedEventHandler.Handle(event)

	assert.Contains(t, loggerOutput.String(), "Invalid event data")
}
