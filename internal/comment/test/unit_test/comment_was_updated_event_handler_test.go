package comment_test

import (
	"encoding/json"
	comment_handler "readmodels/internal/comment/handler"
	"readmodels/internal/model"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var commentWasUpdatedEventHandler *comment_handler.CommentWasUpdatedEventHandler

func setUpCommentWasUpdatedEventHandler(t *testing.T) {
	SetUp(t)
	commentWasUpdatedEventHandler = comment_handler.NewCommentWasUpdatedEventHandler(repository)
}

func TestHandleCommentWasUpdatedEvent(t *testing.T) {
	setUpCommentWasUpdatedEventHandler(t)
	timeNow := time.Now().UTC().Format(model.TimeLayout)
	data := &comment_handler.CommentWasUpdatedEvent{
		CommentId: uint64(123456),
		Content:   "Exemplo de content",
		UpdatedAt: timeNow,
	}
	event, _ := json.Marshal(data)
	expectedTime, _ := time.Parse(model.TimeLayout, timeNow)
	expectedComment := &model.Comment{
		CommentId: data.CommentId,
		Content:   data.Content,
		UpdatedAt: expectedTime,
	}
	repository.EXPECT().UpdateComment(expectedComment).Return(nil)

	commentWasUpdatedEventHandler.Handle(event)
}

func TestInvalidDataInCommentWasUpdatedEventHandler(t *testing.T) {
	setUpCommentWasUpdatedEventHandler(t)
	invalidData := "invalid data"
	event, _ := json.Marshal(invalidData)

	commentWasUpdatedEventHandler.Handle(event)

	assert.Contains(t, loggerOutput.String(), "Invalid event data")
}
