package reaction_test

import (
	"errors"
	"fmt"
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

func TestCreatePostLikeWithService(t *testing.T) {
	setUpService(t)
	data := &model.PostLike{
		User: &model.UserMetadata{
			Username: "user123",
		},
		PostId: "post123",
	}
	repositoryService.EXPECT().CreatePostLike(data)

	reactionService.CreatePostLike(data)

	assert.Contains(t, loggerOutput.String(), "PostLike was created, username: user123 -> postId: post123")
}

func TestErrorOnCreatePostLikeWithService(t *testing.T) {
	setUpService(t)
	data := &model.PostLike{
		User: &model.UserMetadata{
			Username: "user123",
		},
		PostId: "post123",
	}
	repositoryService.EXPECT().CreatePostLike(data).Return(errors.New("some error"))

	reactionService.CreatePostLike(data)

	assert.Contains(t, loggerOutput.String(), "Error creating postLike, username: user123 -> postId: post123")
}

func TestCreatePostSuperlikeWithService(t *testing.T) {
	setUpService(t)
	data := &model.PostSuperlike{
		User: &model.UserMetadata{
			Username: "user123",
		},
		PostId: "post123",
	}
	repositoryService.EXPECT().CreatePostSuperlike(data)

	reactionService.CreatePostSuperlike(data)

	assert.Contains(t, loggerOutput.String(), "PostSuperlike was created, username: user123 -> postId: post123")
}

func TestErrorOnCreatePostSuperlikeWithService(t *testing.T) {
	setUpService(t)
	data := &model.PostSuperlike{
		User: &model.UserMetadata{
			Username: "user123",
		},
		PostId: "post123",
	}
	repositoryService.EXPECT().CreatePostSuperlike(data).Return(errors.New("some error"))

	reactionService.CreatePostSuperlike(data)

	assert.Contains(t, loggerOutput.String(), "Error creating postSuperlike, username: user123 -> postId: post123")
}

func TestGetPostLikesMetadataWithService(t *testing.T) {
	setUpService(t)
	postId := "post1"
	expectedPostLikes := []*model.UserMetadata{
		{
			Username: "username1",
			Name:     "fullname1",
		},
		{
			Username: "username2",
			Name:     "fullname2",
		},
		{
			Username: "username3",
			Name:     "fullname3",
		},
	}
	expectedLastUsername := "username3"
	repositoryService.EXPECT().GetPostLikesMetadata(postId, "", 12).Return(expectedPostLikes, expectedLastUsername, nil)

	postLikes, lastUsername, err := reactionService.GetPostLikesMetadata(postId, "", 12)
	assert.Nil(t, err)
	assert.ElementsMatch(t, expectedPostLikes, postLikes)
	assert.Equal(t, expectedLastUsername, lastUsername)
}

func TestErrorOnGetPostLikesMetadataWithService(t *testing.T) {
	setUpService(t)
	postId := "post1"
	expectedPostLikes := []*model.UserMetadata{}
	expectedLastUsername := ""
	repositoryService.EXPECT().GetPostLikesMetadata(postId, "", 12).Return(expectedPostLikes, expectedLastUsername, errors.New("some error"))

	postLikes, lastUsername, err := reactionService.GetPostLikesMetadata(postId, "", 12)

	assert.Contains(t, loggerOutput.String(), fmt.Sprintf(`Error getting post %s's likes`, postId))
	assert.NotNil(t, err)
	assert.ElementsMatch(t, expectedPostLikes, postLikes)
	assert.Equal(t, expectedLastUsername, lastUsername)
}

func TestDeletePostLikeWithService(t *testing.T) {
	setUpService(t)
	data := &model.PostLike{
		User: &model.UserMetadata{
			Username: "user123",
		},
		PostId: "post123",
	}
	repositoryService.EXPECT().DeletePostLike(data)

	reactionService.DeletePostLike(data)

	assert.Contains(t, loggerOutput.String(), "PostLike was deleted, username: user123 -> postId: post123")
}

func TestErrorOnDeletePostLikeWithService(t *testing.T) {
	setUpService(t)
	data := &model.PostLike{
		User: &model.UserMetadata{
			Username: "user123",
		},
		PostId: "post123",
	}
	repositoryService.EXPECT().DeletePostLike(data).Return(errors.New("some error"))

	reactionService.DeletePostLike(data)

	assert.Contains(t, loggerOutput.String(), "Error deleting postLike, username: user123 -> postId: post123")
}
