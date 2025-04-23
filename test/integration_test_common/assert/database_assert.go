package integration_test_assert

import (
	"fmt"
	database "readmodels/internal/db"
	"readmodels/internal/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func AssertCommentExists(t *testing.T, db *database.Database, expectedCommentId uint64, expectedComment *model.Comment) {
	commentKey := &database.CommentKey{
		CommentId: expectedCommentId,
	}
	var comment model.Comment
	err := db.Client.GetData("readmodels.comments", commentKey, &comment)
	assert.Nil(t, err)
	assert.Equal(t, expectedCommentId, comment.CommentId)
	assert.Equal(t, expectedComment.PostId, comment.PostId)
	assert.Equal(t, expectedComment.Username, comment.Username)
	assert.Equal(t, expectedComment.Content, comment.Content)
}

func AssertCommentDoesNotExist(t *testing.T, db *database.Database, expectedCommentId uint64) {
	commentKey := &database.CommentKey{
		CommentId: expectedCommentId,
	}
	var comment model.Comment
	err := db.Client.GetData("readmodels.comments", commentKey, &comment)
	assert.NotNil(t, err)
	assert.Equal(t, fmt.Sprintf("Data in table readmodels.comments not found for key %v", commentKey), err.Error())
}

func AssertPostCommentsIncreased(t *testing.T, db *database.Database, postId string) {
	postKey := &database.PostMetadataKey{
		PostId: postId,
	}
	var post database.PostMetadata
	err := db.Client.GetData("PostMetadata", postKey, &post)
	assert.Nil(t, err)
	assert.Equal(t, 1, post.Comments)
}

func AssertPostCommentsDecreased(t *testing.T, db *database.Database, postId string) {
	postKey := &database.PostMetadataKey{
		PostId: postId,
	}
	var post database.PostMetadata
	err := db.Client.GetData("PostMetadata", postKey, &post)
	assert.Nil(t, err)
	assert.Equal(t, 0, post.Comments)
}

func AssertPostLikeExists(t *testing.T, db *database.Database, expectedPostLike *model.PostLike) {
	postLikeKey := &database.PostLikeKey{
		PostId:   expectedPostLike.PostId,
		Username: expectedPostLike.User.Username,
	}
	var postLike database.PostLikeMetadata
	err := db.Client.GetData("readmodels.postLikes", postLikeKey, &postLike)
	assert.Nil(t, err)
	assert.Equal(t, expectedPostLike.PostId, postLike.PostId)
	assert.Equal(t, expectedPostLike.User.Username, postLike.Username)
	assert.Equal(t, expectedPostLike.User.Name, postLike.Name)
}

func AssertPostLikeDoesNotExists(t *testing.T, db *database.Database, expectedPostLike *model.PostLike) {
	postLikeKey := &database.PostLikeKey{
		PostId:   expectedPostLike.PostId,
		Username: expectedPostLike.User.Username,
	}
	var postLike database.PostLikeMetadata
	err := db.Client.GetData("readmodels.postLikes", postLikeKey, &postLike)
	assert.NotNil(t, err)
	assert.Equal(t, fmt.Sprintf("Data in table readmodels.postLikes not found for key %v", postLikeKey), err.Error())
}

func AssertPostLikesIncreased(t *testing.T, db *database.Database, postId string) {
	postKey := &database.PostMetadataKey{
		PostId: postId,
	}
	var post database.PostMetadata
	err := db.Client.GetData("PostMetadata", postKey, &post)
	assert.Nil(t, err)
	assert.Equal(t, 1, post.Likes)
}

func AssertPostLikesDecreased(t *testing.T, db *database.Database, postId string) {
	postKey := &database.PostMetadataKey{
		PostId: postId,
	}
	var post database.PostMetadata
	err := db.Client.GetData("PostMetadata", postKey, &post)
	assert.Nil(t, err)
	assert.Equal(t, 0, post.Likes)
}

func AssertPostSuperlikeExists(t *testing.T, db *database.Database, expectedPostSuperlike *model.PostSuperlike) {
	postSuperlikeKey := &database.PostSuperlikeKey{
		PostId:   expectedPostSuperlike.PostId,
		Username: expectedPostSuperlike.User.Username,
	}
	var postSuperlike database.PostSuperlikeMetadata
	err := db.Client.GetData("readmodels.postSuperlikes", postSuperlikeKey, &postSuperlike)
	assert.Nil(t, err)
	assert.Equal(t, expectedPostSuperlike.PostId, postSuperlike.PostId)
	assert.Equal(t, expectedPostSuperlike.User.Username, postSuperlike.Username)
	assert.Equal(t, expectedPostSuperlike.User.Name, postSuperlike.Name)
}

func AssertPostSuperlikeDoesNotExists(t *testing.T, db *database.Database, expectedPostSuperlike *model.PostSuperlike) {
	postLikeKey := &database.PostSuperlikeKey{
		PostId:   expectedPostSuperlike.PostId,
		Username: expectedPostSuperlike.User.Username,
	}
	var postSuperlike database.PostSuperlikeMetadata
	err := db.Client.GetData("readmodels.postSuperlikes", postLikeKey, &postSuperlike)
	assert.NotNil(t, err)
	assert.Equal(t, fmt.Sprintf("Data in table readmodels.postSuperlikes not found for key %v", postLikeKey), err.Error())
}

func AssertPostSuperlikesIncreased(t *testing.T, db *database.Database, postId string) {
	postKey := &database.PostMetadataKey{
		PostId: postId,
	}
	var post database.PostMetadata
	err := db.Client.GetData("PostMetadata", postKey, &post)
	assert.Nil(t, err)
	assert.Equal(t, 1, post.Superlikes)
}

func AssertPostSuperlikesDecreased(t *testing.T, db *database.Database, postId string) {
	postKey := &database.PostMetadataKey{
		PostId: postId,
	}
	var post database.PostMetadata
	err := db.Client.GetData("PostMetadata", postKey, &post)
	assert.Nil(t, err)
	assert.Equal(t, 0, post.Superlikes)
}
