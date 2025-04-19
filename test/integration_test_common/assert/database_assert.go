package integration_test_assert

import (
	"fmt"
	"readmodels/internal/comment"
	database "readmodels/internal/db"
	"testing"

	"github.com/stretchr/testify/assert"
)

func AssertCommentExists(t *testing.T, db *database.Database, expectedCommentId uint64, expectedComment *comment.Comment) {
	commentKey := &database.CommentKey{
		CommentId: expectedCommentId,
	}
	var comment comment.Comment
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
	var comment comment.Comment
	err := db.Client.GetData("readmodels.comments", commentKey, &comment)
	assert.NotNil(t, err)
	assert.Equal(t, fmt.Sprintf("Data in table readmodels.comments not found for key %v", commentKey), err.Error())
}
