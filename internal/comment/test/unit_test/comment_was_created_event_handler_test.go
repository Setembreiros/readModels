package comment_test

import (
	"encoding/json"
	comment_handler "readmodels/internal/comment/handler"
	mock_comment_handler "readmodels/internal/comment/handler/test/mock"
	"readmodels/internal/model"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var commentWasCreatedEventHandler *comment_handler.CommentWasCreatedEventHandler
var commentWasCreatedEventService *mock_comment_handler.MockCommentWasCreatedEventService

func setUpHandler(t *testing.T) {
	SetUp(t)
	commentWasCreatedEventService = mock_comment_handler.NewMockCommentWasCreatedEventService(ctrl)
	commentWasCreatedEventHandler = comment_handler.NewCommentWasCreatedEventHandler(commentWasCreatedEventService)
}

func TestHandleCommentWasCreatedEvent(t *testing.T) {
	setUpHandler(t)
	timeNow := time.Now().UTC().Format(model.TimeLayout)
	data := &comment_handler.CommentWasCreatedEvent{
		CommentId: uint64(123456),
		Username:  "user123",
		PostId:    "post123",
		Content:   "Exemplo de content",
		CreatedAt: timeNow,
	}
	event, _ := json.Marshal(data)
	expectedTime, _ := time.Parse(model.TimeLayout, timeNow)
	expectedComment := &model.Comment{
		CommentId: uint64(123456),
		Username:  "user123",
		PostId:    "post123",
		Content:   "Exemplo de content",
		CreatedAt: expectedTime,
	}
	commentWasCreatedEventService.EXPECT().CreateComment(expectedComment)

	commentWasCreatedEventHandler.Handle(event)
}

func TestInvalidDataInCommentWasCreatedEventHandler(t *testing.T) {
	setUpHandler(t)
	invalidData := "invalid data"
	event, _ := json.Marshal(invalidData)

	commentWasCreatedEventHandler.Handle(event)

	assert.Contains(t, loggerOutput.String(), "Invalid event data")
}
