package comment_test

import (
	"encoding/json"
	"readmodels/internal/comment"
	comment_handler "readmodels/internal/comment/handler"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var handler *comment_handler.CommentWasCreatedEventHandler

func setUpHandler(t *testing.T) {
	SetUp(t)
	handler = comment_handler.NewCommentWasCreatedEventHandler(repository)
}

func TestHandleCommentWasCreatedEvent(t *testing.T) {
	setUpHandler(t)
	timeLayout := "2006-01-02T15:04:05.000000000Z"
	timeNow := time.Now().UTC().Format(timeLayout)
	data := &comment_handler.CommentWasCreatedEvent{
		CommentId: "123456",
		Username:  "user123",
		PostId:    "post123",
		Content:   "Exemplo de content",
		CreatedAt: timeNow,
	}
	event, _ := json.Marshal(data)
	expectedTime, _ := time.Parse(timeLayout, timeNow)
	expectedComment := &comment.Comment{
		CommentId: "123456",
		Username:  "user123",
		PostId:    "post123",
		Content:   "Exemplo de content",
		CreatedAt: expectedTime,
	}
	repository.EXPECT().AddNewComment(expectedComment)

	handler.Handle(event)
}

func TestInvalidDataInCommentWasCreatedEventHandler(t *testing.T) {
	setUpHandler(t)
	invalidData := "invalid data"
	event, _ := json.Marshal(invalidData)

	handler.Handle(event)

	assert.Contains(t, loggerOutput.String(), "Invalid event data")
}
