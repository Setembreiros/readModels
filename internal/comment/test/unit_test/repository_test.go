package comment_test

import (
	"readmodels/internal/comment"
	database "readmodels/internal/db"
	"readmodels/internal/model"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var commentRepository *comment.CommentRepository

func setUpRepository(t *testing.T) {
	SetUp(t)
	commentRepository = comment.NewCommentRepository(database.NewDatabase(client), database.NewCache(cacheClient))
}

func TestCreateCommentInRepository(t *testing.T) {
	setUpRepository(t)
	timeNow := time.Now().UTC()
	data := &model.Comment{
		CommentId: uint64(123456),
		Username:  "user123",
		PostId:    "post123",
		Content:   "Exemplo de content",
		CreatedAt: timeNow,
	}
	client.EXPECT().InsertData("readmodels.comments", data).Return(nil)

	err := commentRepository.AddNewComment(data)

	assert.Nil(t, err)
}

func TestGetCommentsByPostIdInRepository_WhenDatabaseReturnsSuccess(t *testing.T) {
	setUpRepository(t)
	postId := "post2"
	lastCommentId := uint64(7)
	limit := 3
	timeNow := time.Now().UTC()
	data := []*model.Comment{
		{
			CommentId: uint64(5),
			Username:  "username1",
			PostId:    postId,
			Content:   "o meu comentario 1",
			CreatedAt: timeNow,
		},
		{
			CommentId: uint64(6),
			Username:  "username2",
			PostId:    postId,
			Content:   "o meu comentario 2",
			CreatedAt: timeNow,
		},
		{
			CommentId: uint64(7),
			Username:  "username1",
			PostId:    postId,
			Content:   "o meu comentario 3",
			CreatedAt: timeNow,
		},
	}
	expectedResult := []*model.Comment{
		{
			CommentId: uint64(5),
			Username:  "username1",
			PostId:    postId,
			Content:   "o meu comentario 1",
			CreatedAt: timeNow,
		},
		{
			CommentId: uint64(6),
			Username:  "username2",
			PostId:    postId,
			Content:   "o meu comentario 2",
			CreatedAt: timeNow,
		},
		{
			CommentId: uint64(7),
			Username:  "username1",
			PostId:    postId,
			Content:   "o meu comentario 3",
			CreatedAt: timeNow,
		},
	}
	expectedLastCommentId := uint64(7)
	cacheClient.EXPECT().GetPostComments(postId, lastCommentId, limit).Return([]*model.Comment{}, uint64(0), false)
	client.EXPECT().GetCommentsByIndexPostId(postId, lastCommentId, limit).Return(data, expectedLastCommentId, nil)
	cacheClient.EXPECT().SetPostComments(postId, lastCommentId, limit, data)

	result, lastCommentId, err := commentRepository.GetCommentsByPostId(postId, lastCommentId, limit)

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, result)
	assert.Equal(t, expectedLastCommentId, lastCommentId)
}

func TestGetCommentsByPostIdInRepository_WhenCacheeturnsSuccess(t *testing.T) {
	setUpRepository(t)
	postId := "post2"
	lastCommentId := uint64(7)
	limit := 3
	timeNow := time.Now().UTC()
	data := []*model.Comment{
		{
			CommentId: uint64(5),
			Username:  "username1",
			PostId:    postId,
			Content:   "o meu comentario 1",
			CreatedAt: timeNow,
		},
		{
			CommentId: uint64(6),
			Username:  "username2",
			PostId:    postId,
			Content:   "o meu comentario 2",
			CreatedAt: timeNow,
		},
		{
			CommentId: uint64(7),
			Username:  "username1",
			PostId:    postId,
			Content:   "o meu comentario 3",
			CreatedAt: timeNow,
		},
	}
	expectedResult := []*model.Comment{
		{
			CommentId: uint64(5),
			Username:  "username1",
			PostId:    postId,
			Content:   "o meu comentario 1",
			CreatedAt: timeNow,
		},
		{
			CommentId: uint64(6),
			Username:  "username2",
			PostId:    postId,
			Content:   "o meu comentario 2",
			CreatedAt: timeNow,
		},
		{
			CommentId: uint64(7),
			Username:  "username1",
			PostId:    postId,
			Content:   "o meu comentario 3",
			CreatedAt: timeNow,
		},
	}
	expectedLastCommentId := uint64(7)
	cacheClient.EXPECT().GetPostComments(postId, lastCommentId, limit).Return(data, expectedLastCommentId, true)

	result, lastCommentId, err := commentRepository.GetCommentsByPostId(postId, lastCommentId, limit)

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, result)
	assert.Equal(t, expectedLastCommentId, lastCommentId)
}

func TestDeleteCommentInRepository(t *testing.T) {
	setUpRepository(t)
	commentId := uint64(7)
	expectedKey := &database.CommentKey{
		CommentId: commentId,
	}
	client.EXPECT().RemoveData("readmodels.comments", expectedKey)

	err := commentRepository.DeleteComment(commentId)

	assert.Nil(t, err)
}
