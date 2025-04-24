package comment_test

import (
	"errors"
	"fmt"
	"readmodels/internal/comment"
	"readmodels/internal/model"
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
	data := &model.Comment{
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
	data := &model.Comment{
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

func TestGetCommentsByPostIdWithService(t *testing.T) {
	setUpService(t)
	postId := "post1"
	expectedComments := []*model.Comment{
		{
			CommentId: uint64(5),
			Username:  "username1",
			PostId:    "post1",
			Content:   "o meu comentario 1",
			CreatedAt: time.Now(),
		},
		{
			CommentId: uint64(6),
			Username:  "username2",
			PostId:    "post1",
			Content:   "o meu comentario 2",
			CreatedAt: time.Now(),
		},
		{
			CommentId: uint64(7),
			Username:  "username1",
			PostId:    "post1",
			Content:   "o meu comentario 3",
			CreatedAt: time.Now(),
		},
	}
	expectedLastCommentId := uint64(7)
	repository.EXPECT().GetCommentsByPostId(postId, uint64(0), 12).Return(expectedComments, expectedLastCommentId, nil)

	commets, lastCommentId, err := commentService.GetCommentsByPostId(postId, uint64(0), 12)
	assert.Nil(t, err)
	assert.ElementsMatch(t, expectedComments, commets)
	assert.Equal(t, expectedLastCommentId, lastCommentId)
}

func TestErrorOnGetCommentsByPostIdWithService(t *testing.T) {
	setUpService(t)
	postId := "post1"
	expectedComments := []*model.Comment{}
	expectedLastCommentId := uint64(0)
	repository.EXPECT().GetCommentsByPostId(postId, uint64(0), 12).Return(expectedComments, uint64(0), errors.New("some error"))

	commets, lastCommentId, err := commentService.GetCommentsByPostId(postId, uint64(0), 12)

	assert.Contains(t, loggerOutput.String(), fmt.Sprintf("Error getting  %s's comments", postId))
	assert.NotNil(t, err)
	assert.ElementsMatch(t, expectedComments, commets)
	assert.Equal(t, expectedLastCommentId, lastCommentId)
}

func TestDeleteCommentWithService(t *testing.T) {
	setUpService(t)
	commentId := uint64(123456)
	repository.EXPECT().DeleteComment(commentId).Return(nil)

	commentService.DeleteComment(commentId)

	assert.Contains(t, loggerOutput.String(), "Comment with id 123456 was deleted")
}

func TestErrorOnDeleteCommentWithService(t *testing.T) {
	setUpService(t)
	commentId := uint64(123456)
	repository.EXPECT().DeleteComment(commentId).Return(errors.New("some error"))

	commentService.DeleteComment(commentId)

	assert.Contains(t, loggerOutput.String(), "Error deleting comment with id 123456")
}
