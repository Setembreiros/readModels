package comment_test

import (
	"encoding/json"
	comment_handler "readmodels/internal/comment/handler"
	"readmodels/internal/model"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var commentWasCreatedEventHandler *comment_handler.CommentWasCreatedEventHandler

func setUpHandler(t *testing.T) {
	SetUp(t)
	commentWasCreatedEventHandler = comment_handler.NewCommentWasCreatedEventHandler(repository)
}

func TestHandleCommentWasCreatedEvent(t *testing.T) {
	setUpHandler(t)
	timeLayout := "2006-01-02T15:04:05.000000000Z"
	timeNow := time.Now().UTC().Format(timeLayout)
	data := &comment_handler.CommentWasCreatedEvent{
		CommentId: uint64(123456),
		Username:  "user123",
		PostId:    "post123",
		Content:   "Exemplo de content",
		CreatedAt: timeNow,
	}
	event, _ := json.Marshal(data)
	expectedTime, _ := time.Parse(timeLayout, timeNow)
	expectedComment := &model.Comment{
		CommentId: uint64(123456),
		Username:  "user123",
		PostId:    "post123",
		Content:   "Exemplo de content",
		CreatedAt: expectedTime,
	}
	repository.EXPECT().AddNewComment(expectedComment)

	commentWasCreatedEventHandler.Handle(event)
}

func TestInvalidDataInCommentWasCreatedEventHandler(t *testing.T) {
	setUpHandler(t)
	invalidData := "invalid data"
	event, _ := json.Marshal(invalidData)

	commentWasCreatedEventHandler.Handle(event)

	assert.Contains(t, loggerOutput.String(), "Invalid event data")
}
