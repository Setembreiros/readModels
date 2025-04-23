package reaction_test

import (
	database "readmodels/internal/db"
	"readmodels/internal/model"
	"readmodels/internal/reaction"
	"testing"

	"github.com/stretchr/testify/assert"
)

var reactionRepository *reaction.ReactionRepository

func setUpRepository(t *testing.T) {
	SetUp(t)
	reactionRepository = reaction.NewReactionRepository(database.NewDatabase(client), database.NewCache(cacheClient))
}

func TestCreatePostLikeInRepository(t *testing.T) {
	setUpRepository(t)
	data := &model.PostLike{
		User: &model.UserMetadata{
			Username: "user123",
		},
		PostId: "post123",
	}
	expectedData := &database.PostLikeMetadata{
		PostId:   data.PostId,
		Username: data.User.Username,
	}
	expectedUserKey := &database.UserProfileKey{
		Username: data.User.Username,
	}
	expectedPostKey := &database.PostMetadataKey{
		PostId: data.PostId,
	}
	expectedUserFullname := &struct {
		Name string `json:"name"`
	}{}
	client.EXPECT().GetData("UserProfile", expectedUserKey, expectedUserFullname).Return(nil)
	client.EXPECT().InsertDataAndIncreaseCounter("readmodels.postLikes", expectedData, "PostMetadata", expectedPostKey, "Likes").Return(nil)

	err := reactionRepository.CreatePostLike(data)

	assert.Nil(t, err)
}

func TestCreatePostSuperlikeInRepository(t *testing.T) {
	setUpRepository(t)
	data := &model.PostSuperlike{
		User: &model.UserMetadata{
			Username: "user123",
		},
		PostId: "post123",
	}
	expectedData := &database.PostSuperlikeMetadata{
		PostId:   data.PostId,
		Username: data.User.Username,
	}
	expectedUserKey := &database.UserProfileKey{
		Username: data.User.Username,
	}
	expectedPostKey := &database.PostMetadataKey{
		PostId: data.PostId,
	}
	expectedUserFullname := &struct {
		Name string `json:"name"`
	}{}
	client.EXPECT().GetData("UserProfile", expectedUserKey, expectedUserFullname).Return(nil)
	client.EXPECT().InsertDataAndIncreaseCounter("readmodels.postSuperlikes", expectedData, "PostMetadata", expectedPostKey, "Superlikes").Return(nil)

	err := reactionRepository.CreatePostSuperlike(data)

	assert.Nil(t, err)
}

func TestDeletePostLikeInRepository(t *testing.T) {
	setUpRepository(t)
	data := &model.PostLike{
		User: &model.UserMetadata{
			Username: "user123",
		},
		PostId: "post123",
	}
	expectedData := &database.PostLikeKey{
		PostId:   data.PostId,
		Username: data.User.Username,
	}
	expectedPostKey := &database.PostMetadataKey{
		PostId: data.PostId,
	}
	client.EXPECT().RemoveDataAndDecreaseCounter("readmodels.postLikes", expectedData, "PostMetadata", expectedPostKey, "Likes").Return(nil)

	err := reactionRepository.DeletePostLike(data)

	assert.Nil(t, err)
}

func TestDeletePostSUperlikeInRepository(t *testing.T) {
	setUpRepository(t)
	data := &model.PostSuperlike{
		User: &model.UserMetadata{
			Username: "user123",
		},
		PostId: "post123",
	}
	expectedData := &database.PostSuperlikeKey{
		PostId:   data.PostId,
		Username: data.User.Username,
	}
	expectedPostKey := &database.PostMetadataKey{
		PostId: data.PostId,
	}
	client.EXPECT().RemoveDataAndDecreaseCounter("readmodels.postSuperlikes", expectedData, "PostMetadata", expectedPostKey, "Superlikes").Return(nil)

	err := reactionRepository.DeletePostSuperlike(data)

	assert.Nil(t, err)
}
