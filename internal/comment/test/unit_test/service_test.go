package comment_test

import (
	"errors"
	"readmodels/internal/comment"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var commentService *comment.CommentService

func setUpService(t *testing.T) {
	SetUp(t)
	commentService = comment.NewCommentService(repository)
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
	repository.EXPECT().AddNewComment(data)

	commentService.CreateNewComment(data)

	assert.Contains(t, loggerOutput.String(), "Comment with id 123456 was added")
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
	repository.EXPECT().AddNewComment(data).Return(errors.New("some error"))

	commentService.CreateNewComment(data)

	assert.Contains(t, loggerOutput.String(), "Error adding comment with id 123456")
}
