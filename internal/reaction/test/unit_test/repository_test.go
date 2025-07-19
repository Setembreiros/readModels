package reaction_test

import (
	database "readmodels/internal/db"
	"readmodels/internal/model"
	"readmodels/internal/reaction"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var reactionRepository *reaction.ReactionRepository

func setUpRepository(t *testing.T) {
	SetUp(t)
	reactionRepository = reaction.NewReactionRepository(database.NewDatabase(client))
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

func TestCreateReviewInRepository(t *testing.T) {
	setUpRepository(t)
	timeNow := time.Now().UTC()
	data := &model.Review{
		ReviewId:  uint64(123456),
		Username:  "user123",
		PostId:    "post123",
		Content:   "Exemplo de content",
		Rating:    3,
		CreatedAt: timeNow,
	}
	expectedPostKey := &database.PostMetadataKey{
		PostId: data.PostId,
	}
	client.EXPECT().InsertDataAndIncreaseCounter("readmodels.reviews", data, "PostMetadata", expectedPostKey, "Reviews").Return(nil)

	err := reactionRepository.CreateReview(data)

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
	client.EXPECT().GetPostLikesByIndexPostId(postId, lastUsername, limit).Return(expectedResult, expectedLastUsername, nil)

	result, lastUsername, err := reactionRepository.GetLikesMetadataByPostId(postId, lastUsername, limit)

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, result)
	assert.Equal(t, expectedLastUsername, lastUsername)
}

func TestGetPostSuperlikesMetadataInRepository_WhenDatabaseReturnsSuccess(t *testing.T) {
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
	client.EXPECT().GetPostSuperlikesByIndexPostId(postId, lastUsername, limit).Return(expectedResult, expectedLastUsername, nil)

	result, lastUsername, err := reactionRepository.GetSuperlikesMetadataByPostId(postId, lastUsername, limit)

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, result)
	assert.Equal(t, expectedLastUsername, lastUsername)
}

func TestGetReviewsByPostIdInRepository_WhenDatabaseReturnsSuccess(t *testing.T) {
	setUpRepository(t)
	postId := "post2"
	lastReviewId := uint64(7)
	limit := 3
	timeNow := time.Now().UTC()
	data := []*model.Review{
		{
			ReviewId:  uint64(5),
			Username:  "username1",
			PostId:    postId,
			Content:   "a miña review 1",
			Rating:    4,
			CreatedAt: timeNow,
		},
		{
			ReviewId:  uint64(6),
			Username:  "username2",
			PostId:    postId,
			Content:   "a miña review 2",
			Rating:    4,
			CreatedAt: timeNow,
		},
		{
			ReviewId:  uint64(7),
			Username:  "username1",
			PostId:    postId,
			Content:   "a miña review 3",
			Rating:    4,
			CreatedAt: timeNow,
		},
	}
	expectedResult := []*model.Review{
		{
			ReviewId:  uint64(5),
			Username:  "username1",
			PostId:    postId,
			Content:   "a miña review 1",
			Rating:    4,
			CreatedAt: timeNow,
		},
		{
			ReviewId:  uint64(6),
			Username:  "username2",
			PostId:    postId,
			Content:   "a miña review 2",
			Rating:    4,
			CreatedAt: timeNow,
		},
		{
			ReviewId:  uint64(7),
			Username:  "username1",
			PostId:    postId,
			Content:   "a miña review 3",
			Rating:    4,
			CreatedAt: timeNow,
		},
	}
	expectedLastReviewId := uint64(7)
	client.EXPECT().GetReviewsByIndexPostId(postId, lastReviewId, limit).Return(data, expectedLastReviewId, nil)

	result, lastReviewId, err := reactionRepository.GetReviewsByPostId(postId, lastReviewId, limit)

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, result)
	assert.Equal(t, expectedLastReviewId, lastReviewId)
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
