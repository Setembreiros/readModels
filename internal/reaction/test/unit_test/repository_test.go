package reaction_test

import (
	database "readmodels/internal/db"
	"readmodels/internal/model"
	"readmodels/internal/reaction"
	"testing"

	"github.com/stretchr/testify/assert"
)

var reactionRepository *reaction.ReactionRepository

func setUpRepository(t *testing.T) {
	SetUp(t)
	reactionRepository = reaction.NewReactionRepository(database.NewDatabase(client), database.NewCache(cacheClient))
}

func TestCreateLikePostInRepository(t *testing.T) {
	setUpRepository(t)
	data := &model.LikePost{
		Username: "user123",
		PostId:   "post123",
	}
	expectedPostKey := &database.PostMetadataKey{
		PostId: data.PostId,
	}
	client.EXPECT().InsertDataAndIncreaseCounter("readmodels.likes", data, "PostMetadata", expectedPostKey, "Likes").Return(nil)

	err := reactionRepository.CreateLikePost(data)

	assert.Nil(t, err)
}
