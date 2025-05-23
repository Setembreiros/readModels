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
	repository.EXPECT().CreateComment(data)

	commentService.CreateComment(data)

	assert.Contains(t, loggerOutput.String(), "Comment with id 123456 in post post123 was created")
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
	repository.EXPECT().CreateComment(data).Return(errors.New("some error"))

	commentService.CreateComment(data)

	assert.Contains(t, loggerOutput.String(), "Error creating comment with id 123456")
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

func TestUpdateCommentWithService(t *testing.T) {
	setUpService(t)
	timeNow := time.Now().UTC()
	data := &model.Comment{
		CommentId: uint64(123456),
		Content:   "Exemplo de content",
		UpdatedAt: timeNow,
	}
	repository.EXPECT().UpdateComment(data)

	commentService.UpdateComment(data)

	assert.Contains(t, loggerOutput.String(), "Comment with id 123456 was updated")
}

func TestErrorOnUpdateNewCommentWithService(t *testing.T) {
	setUpService(t)
	timeNow := time.Now().UTC()
	data := &model.Comment{
		CommentId: uint64(123456),
		Content:   "Exemplo de content",
		UpdatedAt: timeNow,
	}
	repository.EXPECT().UpdateComment(data).Return(errors.New("some error"))

	commentService.UpdateComment(data)

	assert.Contains(t, loggerOutput.String(), "Error updating comment with id 123456")
}

func TestDeleteCommentWithService(t *testing.T) {
	setUpService(t)
	commentId := uint64(123456)
	postId := "post1"
	repository.EXPECT().DeleteComment(postId, commentId).Return(nil)

	commentService.DeleteComment(postId, commentId)

	assert.Contains(t, loggerOutput.String(), "Comment with id 123456 in post post1 was deleted")
}

func TestErrorOnDeleteCommentWithService(t *testing.T) {
	setUpService(t)
	commentId := uint64(123456)
	posId := "post1"
	repository.EXPECT().DeleteComment(posId, commentId).Return(errors.New("some error"))

	commentService.DeleteComment(posId, commentId)

	assert.Contains(t, loggerOutput.String(), "Error deleting comment with id 123456")
}
