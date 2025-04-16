package comment_test

import (
	"bytes"
	"errors"
	"readmodels/internal/comment"
	mock_comment "readmodels/internal/comment/mock"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

var serviceLoggerOutput bytes.Buffer
var serviceRepository *mock_comment.MockRepository
var commentService *comment.CommentService

func setUpService(t *testing.T) {
	ctrl := gomock.NewController(t)
	serviceRepository = mock_comment.NewMockRepository(ctrl)
	log.Logger = log.Output(&serviceLoggerOutput)
	commentService = comment.NewCommentService(serviceRepository)
}

func TestCreateNewCommentWithService(t *testing.T) {
	setUpService(t)
	timeNow := time.Now().UTC()
	data := &comment.Comment{
		CommentId: uint64(123456),
		Username:  "user123",
		PostId:    "post123",
		Content:   "Exemplo de content",
		CreatedAt: timeNow,
	}
	serviceRepository.EXPECT().AddNewComment(data)

	commentService.CreateNewComment(data)

	assert.Contains(t, serviceLoggerOutput.String(), "Comment with id 123456 was added")
}

func TestErrorOnCreateNewCommentWithService(t *testing.T) {
	setUpService(t)
	timeNow := time.Now().UTC()
	data := &comment.Comment{
		CommentId: uint64(123456),
		Username:  "user123",
		PostId:    "post123",
		Content:   "Exemplo de content",
		CreatedAt: timeNow,
	}
	serviceRepository.EXPECT().AddNewComment(data).Return(errors.New("some error"))

	commentService.CreateNewComment(data)

	assert.Contains(t, serviceLoggerOutput.String(), "Error adding comment with id 123456")
}
