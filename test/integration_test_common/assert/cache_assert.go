package integration_test_assert

import (
	database "readmodels/internal/db"
	"readmodels/internal/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func AssertCachedPostCommentsExists(t *testing.T, db *database.Cache, postId string, lastCommentId uint64, limit int, expectedComments []*model.Comment) {
	cachedComments, cachedLastCommentId, found := db.Client.GetPostComments(postId, lastCommentId, limit)
	assert.Equal(t, true, found)
	assert.Equal(t, expectedComments, cachedComments)
	assert.Equal(t, expectedComments[len(expectedComments)-1].CommentId, cachedLastCommentId)
}
