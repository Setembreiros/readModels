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

func AssertLikePostExists(t *testing.T, db *database.Database, expectedLikePost *model.LikePost) {
	likePostKey := &database.LikePostKey{
		PostId:   expectedLikePost.PostId,
		Username: expectedLikePost.Username,
	}
	var likePost model.LikePost
	err := db.Client.GetData("readmodels.likePosts", likePostKey, &likePost)
	assert.Nil(t, err)
	assert.Equal(t, expectedLikePost.PostId, likePost.PostId)
	assert.Equal(t, expectedLikePost.Username, likePost.Username)
}

func AssertLikePostDoesNotExists(t *testing.T, db *database.Database, expectedLikePost *model.LikePost) {
	likePostKey := &database.LikePostKey{
		PostId:   expectedLikePost.PostId,
		Username: expectedLikePost.Username,
	}
	var likePost model.LikePost
	err := db.Client.GetData("readmodels.likePosts", likePostKey, &likePost)
	assert.NotNil(t, err)
	assert.Equal(t, fmt.Sprintf("Data in table readmodels.likePosts not found for key %v", likePostKey), err.Error())
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

func AssertSuperlikePostExists(t *testing.T, db *database.Database, expectedSuperlikePost *model.SuperlikePost) {
	superlikePostKey := &database.SuperlikePostKey{
		PostId:   expectedSuperlikePost.PostId,
		Username: expectedSuperlikePost.Username,
	}
	var superlikePost model.SuperlikePost
	err := db.Client.GetData("readmodels.superlikePosts", superlikePostKey, &superlikePost)
	assert.Nil(t, err)
	assert.Equal(t, expectedSuperlikePost.PostId, superlikePost.PostId)
	assert.Equal(t, expectedSuperlikePost.Username, superlikePost.Username)
}

func AssertSuperlikePostDoesNotExists(t *testing.T, db *database.Database, expectedSuperlikePost *model.SuperlikePost) {
	likePostKey := &database.SuperlikePostKey{
		PostId:   expectedSuperlikePost.PostId,
		Username: expectedSuperlikePost.Username,
	}
	var superlikePost model.SuperlikePost
	err := db.Client.GetData("readmodels.superlikePosts", likePostKey, &superlikePost)
	assert.NotNil(t, err)
	assert.Equal(t, fmt.Sprintf("Data in table readmodels.superlikePosts not found for key %v", likePostKey), err.Error())
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
