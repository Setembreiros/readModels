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
