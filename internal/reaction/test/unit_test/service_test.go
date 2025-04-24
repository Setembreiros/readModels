package reaction_test

import (
	"errors"
	"readmodels/internal/model"
	"readmodels/internal/reaction"
	mock_reaction "readmodels/internal/reaction/test/mock"
	"testing"

	"github.com/stretchr/testify/assert"
)

var repositoryService *mock_reaction.MockRepository
var reactionService *reaction.ReactionService

func setUpService(t *testing.T) {
	SetUp(t)
	repositoryService = mock_reaction.NewMockRepository(ctrl)
	reactionService = reaction.NewReactionService(repositoryService)
}

func TestCreateLikePostWithService(t *testing.T) {
	setUpService(t)
	data := &model.LikePost{
		Username: "user123",
		PostId:   "post123",
	}
	repositoryService.EXPECT().CreateLikePost(data)

	reactionService.CreateLikePost(data)

	assert.Contains(t, loggerOutput.String(), "LikePost was created, username: user123 -> postId: post123")
}

func TestErrorOnCreateLikePostWithService(t *testing.T) {
	setUpService(t)
	data := &model.LikePost{
		Username: "user123",
		PostId:   "post123",
	}
	repositoryService.EXPECT().CreateLikePost(data).Return(errors.New("some error"))

	reactionService.CreateLikePost(data)

	assert.Contains(t, loggerOutput.String(), "Error creating likePost, username: user123 -> postId: post123")
}

func TestCreateSuperlikePostWithService(t *testing.T) {
	setUpService(t)
	data := &model.SuperlikePost{
		Username: "user123",
		PostId:   "post123",
	}
	repositoryService.EXPECT().CreateSuperlikePost(data)

	reactionService.CreateSuperlikePost(data)

	assert.Contains(t, loggerOutput.String(), "SuperlikePost was created, username: user123 -> postId: post123")
}

func TestErrorOnCreateSuperlikePostWithService(t *testing.T) {
	setUpService(t)
	data := &model.SuperlikePost{
		Username: "user123",
		PostId:   "post123",
	}
	repositoryService.EXPECT().CreateSuperlikePost(data).Return(errors.New("some error"))

	reactionService.CreateSuperlikePost(data)

	assert.Contains(t, loggerOutput.String(), "Error creating superlikePost, username: user123 -> postId: post123")
}

func TestDeleteLikePostWithService(t *testing.T) {
	setUpService(t)
	data := &model.LikePost{
		Username: "user123",
		PostId:   "post123",
	}
	repositoryService.EXPECT().DeleteLikePost(data)

	reactionService.DeleteLikePost(data)

	assert.Contains(t, loggerOutput.String(), "LikePost was deleted, username: user123 -> postId: post123")
}

func TestErrorOnDeleteLikePostWithService(t *testing.T) {
	setUpService(t)
	data := &model.LikePost{
		Username: "user123",
		PostId:   "post123",
	}
	repositoryService.EXPECT().DeleteLikePost(data).Return(errors.New("some error"))

	reactionService.DeleteLikePost(data)

	assert.Contains(t, loggerOutput.String(), "Error deleting likePost, username: user123 -> postId: post123")
}
