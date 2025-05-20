package post_handler_test

import (
	"bytes"
	"encoding/json"
	post_handler "readmodels/internal/post/handler"
	mock_post "readmodels/internal/post/mock"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

var postsWereDeletedEventHandlerLoggerOutput bytes.Buffer
var postsWereDeletedEventHandlerRepository *mock_post.MockRepository
var postsWereDeletedEventHandler *post_handler.PostsWereDeletedEventHandler

func setUpPostsWereDeletedEventHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	postsWereDeletedEventHandlerRepository = mock_post.NewMockRepository(ctrl)
	log.Logger = log.Output(&postsWereDeletedEventHandlerLoggerOutput)
	postsWereDeletedEventHandler = post_handler.NewPostsWereDeletedEventHandler(postsWereDeletedEventHandlerRepository)
}

func TestHandlePostsWereDeletedEvent(t *testing.T) {
	setUpPostsWereDeletedEventHandler(t)
	username := "username1"
	postIds := []string{"123456", "abcdef", "1a2b3e"}
	data := &post_handler.PostsWereDeletedEvent{
		Username: username,
		PostIds:  postIds,
	}
	event, _ := json.Marshal(data)
	postsWereDeletedEventHandlerRepository.EXPECT().RemovePostMetadata(username, postIds).Return(nil)

	postsWereDeletedEventHandler.Handle(event)
}

func TestHandlePostsWereDeletedEvent_ErrorInvalidData(t *testing.T) {
	setUpPostsWereDeletedEventHandler(t)
	invalidData := "invalid data"
	event, _ := json.Marshal(invalidData)

	postsWereDeletedEventHandler.Handle(event)

	assert.Contains(t, postsWereDeletedEventHandlerLoggerOutput.String(), "Invalid event data")
}
