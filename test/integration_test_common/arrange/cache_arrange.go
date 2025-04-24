package integration_test_arrange

import (
	"context"
	"readmodels/cmd/provider"
	database "readmodels/internal/db"
	"readmodels/internal/model"
	"testing"
)

func CreateTestCache(t *testing.T, ctx context.Context) *database.Cache {
	provider := provider.NewProvider("test")
	return provider.ProvideCache(ctx)
}

func AddCachedCommentsToCache(t *testing.T, cache *database.Cache, postId string, lastCommentId uint64, limit int, comments []*model.Comment) {
	cache.Client.SetPostComments(postId, lastCommentId, limit, comments)
}

func AddCachedPostLikesToCache(t *testing.T, cache *database.Cache, postId string, lastUsername string, limit int, postLikes []*model.UserMetadata) {
	cache.Client.SetPostLikes(postId, lastUsername, limit, postLikes)
}

func AddCachedPostSuperlikesToCache(t *testing.T, cache *database.Cache, postId string, lastUsername string, limit int, postSuperlikes []*model.UserMetadata) {
	cache.Client.SetPostSuperlikes(postId, lastUsername, limit, postSuperlikes)
}
