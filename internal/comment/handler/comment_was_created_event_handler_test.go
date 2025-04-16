package comment_handler_test

import (
	"bytes"
	"encoding/json"
	"readmodels/internal/comment"
	comment_handler "readmodels/internal/comment/handler"
	mock_comment "readmodels/internal/comment/mock"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

var loggerOutput bytes.Buffer
var repository *mock_comment.MockRepository
var handler *comment_handler.CommentWasCreatedEventHandler

func setUpHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	repository = mock_comment.NewMockRepository(ctrl)
	log.Logger = log.Output(&loggerOutput)
	handler = comment_handler.NewCommentWasCreatedEventHandler(repository)
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
	expectedComment := &comment.Comment{
		CommentId: uint64(123456),
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
