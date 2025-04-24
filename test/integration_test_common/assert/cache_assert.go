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

func AssertCachedPostLikesExists(t *testing.T, db *database.Cache, postId string, lastUsername string, limit int, expectedPostLikes []*model.UserMetadata) {
	cachedPostLikes, cachedLastUsername, found := db.Client.GetPostLikes(postId, lastUsername, limit)
	assert.Equal(t, true, found)
	assert.Equal(t, expectedPostLikes, cachedPostLikes)
	assert.Equal(t, expectedPostLikes[len(expectedPostLikes)-1].Username, cachedLastUsername)
}

func AssertCachedPostSuperlikesExists(t *testing.T, db *database.Cache, postId string, lastUsername string, limit int, expectedPostSuperlikes []*model.UserMetadata) {
	cachedPostSuperlikes, cachedLastUsername, found := db.Client.GetPostSuperlikes(postId, lastUsername, limit)
	assert.Equal(t, true, found)
	assert.Equal(t, expectedPostSuperlikes, cachedPostSuperlikes)
	assert.Equal(t, expectedPostSuperlikes[len(expectedPostSuperlikes)-1].Username, cachedLastUsername)
}
