package integration_test_arrange

import (
	"context"
	"readmodels/cmd/provider"
	database "readmodels/internal/db"
	"readmodels/internal/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func CreateTestDatabase(t *testing.T, ctx context.Context) *database.Database {
	provider := provider.NewProvider("test")
	db, _ := provider.ProvideDb(ctx)
	db.ApplyMigrations(ctx)
	return db
}

func AddUserProfileToDatabase(t *testing.T, db *database.Database, data *model.UserProfile) {
	err := db.Client.InsertData("UserProfile", data)
	assert.Nil(t, err)
	userProfileKey := &database.UserProfileKey{
		Username: data.Username,
	}
	var userProfile model.UserProfile
	err = db.Client.GetData("UserProfile", userProfileKey, &userProfile)
	assert.Nil(t, err)
	assert.Equal(t, userProfile.Username, data.Username)
	assert.Equal(t, userProfile.Name, data.Name)
	assert.Equal(t, userProfile.Bio, data.Bio)
	assert.Equal(t, userProfile.Link, data.Link)
}

func AddCommentToDatabase(t *testing.T, db *database.Database, data *model.Comment) {
	err := db.Client.InsertData("readmodels.comments", data)
	assert.Nil(t, err)
	commentKey := &database.CommentKey{
		CommentId: data.CommentId,
	}
	var comment model.Comment
	err = db.Client.GetData("readmodels.comments", commentKey, &comment)
	assert.Nil(t, err)
	assert.Equal(t, comment.CommentId, data.CommentId)
	assert.Equal(t, comment.Username, data.Username)
	assert.Equal(t, comment.PostId, data.PostId)
	assert.Equal(t, comment.Content, data.Content)

	post := &database.PostMetadata{
		PostId:   data.PostId,
		Username: data.Username,
		Type:     "TEXT",
		Comments: 1,
	}
	err = db.Client.InsertData("PostMetadata", post)
	assert.Nil(t, err)
	var existingPost database.PostMetadata
	postKey := &database.PostMetadataKey{
		PostId: post.PostId,
	}
	err = db.Client.GetData("PostMetadata", postKey, &existingPost)
	assert.Nil(t, err)
	assert.Equal(t, existingPost.PostId, post.PostId)
	assert.Equal(t, existingPost.Comments, 1)
}

func AddPostToDatabase(t *testing.T, db *database.Database, data *database.PostMetadata) {
	err := db.Client.InsertData("PostMetadata", data)
	assert.Nil(t, err)
	postKey := &database.PostMetadataKey{
		PostId: data.PostId,
	}
	var post database.PostMetadata
	err = db.Client.GetData("PostMetadata", postKey, &post)
	assert.Nil(t, err)
	assert.Equal(t, data.Title, post.Title)
	assert.Equal(t, data.Username, post.Username)
	assert.Equal(t, data.PostId, post.PostId)
	assert.Equal(t, data.Description, post.Description)
	assert.Equal(t, data.Comments, post.Comments)
	assert.Equal(t, data.Likes, post.Likes)
}

func AddPostLikeToDatabase(t *testing.T, db *database.Database, data *database.PostLikeMetadata) {
	err := db.Client.InsertData("readmodels.postLikes", data)
	assert.Nil(t, err)
	likeKey := &database.PostLikeKey{
		PostId:   data.PostId,
		Username: data.Username,
	}
	var postLike database.PostLikeMetadata
	err = db.Client.GetData("readmodels.postLikes", likeKey, &postLike)
	assert.Nil(t, err)
	assert.Equal(t, postLike.PostId, data.PostId)
	assert.Equal(t, postLike.Username, data.Username)
	assert.Equal(t, postLike.Name, data.Name)
}

func AddPostSuperlikeToDatabase(t *testing.T, db *database.Database, data *database.PostSuperlikeMetadata) {
	err := db.Client.InsertData("readmodels.postSuperlikes", data)
	assert.Nil(t, err)
	likeKey := &database.PostSuperlikeKey{
		PostId:   data.PostId,
		Username: data.Username,
	}
	var postSuperlike database.PostSuperlikeMetadata
	err = db.Client.GetData("readmodels.postSuperlikes", likeKey, &postSuperlike)
	assert.Nil(t, err)
	assert.Equal(t, postSuperlike.PostId, data.PostId)
	assert.Equal(t, postSuperlike.Username, data.Username)
	assert.Equal(t, postSuperlike.Name, data.Name)
}
