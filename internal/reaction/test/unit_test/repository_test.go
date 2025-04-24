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

func TestGetPostLikesMetadataInRepository_WhenDatabaseReturnsSuccess(t *testing.T) {
	setUpRepository(t)
	postId := "post2"
	lastUsername := "username0"
	limit := 3
	expectedResult := []*model.UserMetadata{
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
	cacheClient.EXPECT().GetPostLikes(postId, lastUsername, limit).Return([]*model.UserMetadata{}, "", false)
	client.EXPECT().GetPostLikesByIndexPostId(postId, lastUsername, limit).Return(expectedResult, expectedLastUsername, nil)
	cacheClient.EXPECT().SetPostLikes(postId, lastUsername, limit, expectedResult)

	result, lastUsername, err := reactionRepository.GetPostLikesMetadata(postId, lastUsername, limit)

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, result)
	assert.Equal(t, expectedLastUsername, lastUsername)
}

func TestGetPostLikesMetadataInRepository_WhenCacheeturnsSuccess(t *testing.T) {
	setUpRepository(t)
	postId := "post2"
	lastUsername := "username0"
	limit := 3
	expectedResult := []*model.UserMetadata{
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
	cacheClient.EXPECT().GetPostLikes(postId, lastUsername, limit).Return(expectedResult, expectedLastUsername, true)

	result, lastUsername, err := reactionRepository.GetPostLikesMetadata(postId, lastUsername, limit)

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, result)
	assert.Equal(t, expectedLastUsername, lastUsername)
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
