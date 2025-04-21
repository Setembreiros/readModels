package integration_test_reaction

import (
	"context"
	database "readmodels/internal/db"
	"readmodels/internal/model"
	"readmodels/internal/reaction"
	reaction_handler "readmodels/internal/reaction/handler"
	integration_test_arrange "readmodels/test/integration_test_common/arrange"
	integration_test_assert "readmodels/test/integration_test_common/assert"
	"readmodels/test/test_common"
	"testing"
)

var db *database.Database
var cache *database.Cache
var userLikedPostEventHandler *reaction_handler.UserLikedPostEventHandler

func setUp(t *testing.T) {
	// Real infrastructure and services
	ctx := context.TODO()
	db = integration_test_arrange.CreateTestDatabase(t, ctx)
	cache = integration_test_arrange.CreateTestCache(t, ctx)
	repository := reaction.NewReactionRepository(db, cache)
	service := reaction.NewReactionService(repository)
	userLikedPostEventHandler = reaction_handler.NewUserLikedPostEventHandler(service)
}

func tearDown() {
	db.Client.Clean()
	cache.Client.Clean()
}

func TestCreateLikePost_WhenDatabaseReturnsSuccess(t *testing.T) {
	setUp(t)
	defer tearDown()
	existingPost := &database.PostMetadata{
		PostId:   "post123",
		Username: "username1",
		Type:     "TEXT",
	}
	integration_test_arrange.AddPostToDatabase(t, db, existingPost)
	data := &reaction_handler.UserLikedPostEvent{
		Username: "user123",
		PostId:   existingPost.PostId,
	}
	event, _ := test_common.SerializeData(data)
	expectedLike := &model.LikePost{
		Username: data.Username,
		PostId:   data.PostId,
	}

	userLikedPostEventHandler.Handle(event)

	integration_test_assert.AssertLikePostExists(t, db, expectedLike)
	integration_test_assert.AssertPostLikesIncreased(t, db, existingPost.PostId)
}
