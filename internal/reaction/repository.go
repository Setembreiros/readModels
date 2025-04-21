package reaction

import (
	database "readmodels/internal/db"
	"readmodels/internal/model"
)

type ReactionRepository struct {
	cache    *database.Cache
	database *database.Database
}

func NewReactionRepository(database *database.Database, cache *database.Cache) *ReactionRepository {
	return &ReactionRepository{
		cache:    cache,
		database: database,
	}
}

func (r ReactionRepository) CreateLikePost(data *model.LikePost) error {
	postKey := &database.PostMetadataKey{
		PostId: data.PostId,
	}
	return r.database.Client.InsertDataAndIncreaseCounter("readmodels.likePosts", data, "PostMetadata", postKey, "Likes")
}
