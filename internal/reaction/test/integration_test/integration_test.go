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
var userSuperlikedPostEventHandler *reaction_handler.UserSuperlikedPostEventHandler
var userUnlikedPostEventHandler *reaction_handler.UserUnlikedPostEventHandler
var userUnsuperlikedPostEventHandler *reaction_handler.UserUnsuperlikedPostEventHandler

func setUp(t *testing.T) {
	// Real infrastructure and services
	ctx := context.TODO()
	db = integration_test_arrange.CreateTestDatabase(t, ctx)
	cache = integration_test_arrange.CreateTestCache(t, ctx)
	repository := reaction.NewReactionRepository(db, cache)
	service := reaction.NewReactionService(repository)
	userLikedPostEventHandler = reaction_handler.NewUserLikedPostEventHandler(service)
	userSuperlikedPostEventHandler = reaction_handler.NewUserSuperlikedPostEventHandler(service)
	userUnlikedPostEventHandler = reaction_handler.NewUserUnlikedPostEventHandler(service)
	userUnsuperlikedPostEventHandler = reaction_handler.NewUserUnsuperlikedPostEventHandler(service)
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
		Likes:    0,
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

func TestCreateSuperlikePost_WhenDatabaseReturnsSuccess(t *testing.T) {
	setUp(t)
	defer tearDown()
	existingPost := &database.PostMetadata{
		PostId:     "post123",
		Username:   "username1",
		Type:       "TEXT",
		Superlikes: 0,
	}
	integration_test_arrange.AddPostToDatabase(t, db, existingPost)
	data := &reaction_handler.UserSuperlikedPostEvent{
		Username: "user123",
		PostId:   existingPost.PostId,
	}
	event, _ := test_common.SerializeData(data)
	expectedSuperlike := &model.SuperlikePost{
		Username: data.Username,
		PostId:   data.PostId,
	}

	userSuperlikedPostEventHandler.Handle(event)

	integration_test_assert.AssertSuperlikePostExists(t, db, expectedSuperlike)
	integration_test_assert.AssertPostSuperlikesIncreased(t, db, existingPost.PostId)
}

func TestDeleteLikePost_WhenDatabaseReturnsSuccess(t *testing.T) {
	setUp(t)
	defer tearDown()
	existingPost := &database.PostMetadata{
		PostId:   "post123",
		Username: "username1",
		Type:     "TEXT",
		Likes:    1,
	}
	integration_test_arrange.AddPostToDatabase(t, db, existingPost)
	data := &reaction_handler.UserUnlikedPostEvent{
		Username: "user123",
		PostId:   existingPost.PostId,
	}
	event, _ := test_common.SerializeData(data)
	expectedLike := &model.LikePost{
		Username: data.Username,
		PostId:   data.PostId,
	}

	userUnlikedPostEventHandler.Handle(event)

	integration_test_assert.AssertLikePostDoesNotExists(t, db, expectedLike)
	integration_test_assert.AssertPostLikesDecreased(t, db, existingPost.PostId)
}

func TestDeleteSuperlikePost_WhenDatabaseReturnsSuccess(t *testing.T) {
	setUp(t)
	defer tearDown()
	existingPost := &database.PostMetadata{
		PostId:     "post123",
		Username:   "username1",
		Type:       "TEXT",
		Superlikes: 1,
	}
	integration_test_arrange.AddPostToDatabase(t, db, existingPost)
	data := &reaction_handler.UserUnsuperlikedPostEvent{
		Username: "user123",
		PostId:   existingPost.PostId,
	}
	event, _ := test_common.SerializeData(data)
	expectedSuperlike := &model.SuperlikePost{
		Username: data.Username,
		PostId:   data.PostId,
	}

	userUnsuperlikedPostEventHandler.Handle(event)

	integration_test_assert.AssertSuperlikePostDoesNotExists(t, db, expectedSuperlike)
	integration_test_assert.AssertPostSuperlikesDecreased(t, db, existingPost.PostId)
}
