package reaction

import (
	database "readmodels/internal/db"
	"readmodels/internal/model"
)

type ReactionRepository struct {
	database *database.Database
}

func NewReactionRepository(database *database.Database) *ReactionRepository {
	return &ReactionRepository{
		database: database,
	}
}

func (r *ReactionRepository) CreatePostLike(postLike *model.PostLike) error {
	userFullname, err := r.getUserFullname(postLike.User.Username)
	if err != nil {
		return err
	}

	postLikeMetadata := &database.PostLikeMetadata{
		PostId:   postLike.PostId,
		Username: postLike.User.Username,
		Name:     userFullname,
	}

	postKey := &database.PostMetadataKey{
		PostId: postLikeMetadata.PostId,
	}
	return r.database.Client.InsertDataAndIncreaseCounter("readmodels.postLikes", postLikeMetadata, "PostMetadata", postKey, "Likes")
}

func (r *ReactionRepository) CreatePostSuperlike(postSuperlike *model.PostSuperlike) error {
	userFullname, err := r.getUserFullname(postSuperlike.User.Username)
	if err != nil {
		return err
	}

	postSuperlikeMetadata := &database.PostSuperlikeMetadata{
		PostId:   postSuperlike.PostId,
		Username: postSuperlike.User.Username,
		Name:     userFullname,
	}

	postKey := &database.PostMetadataKey{
		PostId: postSuperlike.PostId,
	}
	return r.database.Client.InsertDataAndIncreaseCounter("readmodels.postSuperlikes", postSuperlikeMetadata, "PostMetadata", postKey, "Superlikes")
}

func (r ReactionRepository) CreateReview(data *model.Review) error {
	postKey := &database.PostMetadataKey{
		PostId: data.PostId,
	}
	return r.database.Client.InsertDataAndIncreaseCounter("readmodels.reviews", data, "PostMetadata", postKey, "Reviews")
}

func (r *ReactionRepository) GetPostLikesMetadata(postId string, lastUsername string, limit int) ([]*model.UserMetadata, string, error) {
	likes, newLastUsername, err := r.database.Client.GetPostLikesByIndexPostId(postId, lastUsername, limit)
	if err != nil {
		return []*model.UserMetadata{}, "", err
	}

	return likes, newLastUsername, nil
}

func (r *ReactionRepository) GetPostSuperlikesMetadata(postId string, lastUsername string, limit int) ([]*model.UserMetadata, string, error) {
	superlikes, newLastUsername, err := r.database.Client.GetPostSuperlikesByIndexPostId(postId, lastUsername, limit)
	if err != nil {
		return []*model.UserMetadata{}, "", err
	}

	return superlikes, newLastUsername, nil
}

func (r *ReactionRepository) DeletePostLike(postLike *model.PostLike) error {
	postKey := &database.PostMetadataKey{
		PostId: postLike.PostId,
	}

	postLikeMetadataKey := &database.PostLikeKey{
		PostId:   postLike.PostId,
		Username: postLike.User.Username,
	}
	return r.database.Client.RemoveDataAndDecreaseCounter("readmodels.postLikes", postLikeMetadataKey, "PostMetadata", postKey, "Likes")
}

func (r *ReactionRepository) DeletePostSuperlike(postSuperLike *model.PostSuperlike) error {
	postSuperLikeMetadataKey := &database.PostSuperlikeKey{
		PostId:   postSuperLike.PostId,
		Username: postSuperLike.User.Username,
	}
	postKey := &database.PostMetadataKey{
		PostId: postSuperLike.PostId,
	}
	return r.database.Client.RemoveDataAndDecreaseCounter("readmodels.postSuperlikes", postSuperLikeMetadataKey, "PostMetadata", postKey, "Superlikes")
}

func (r *ReactionRepository) getUserFullname(username string) (string, error) {
	userKey := &database.UserProfileKey{
		Username: username,
	}

	userFullname := &struct {
		Name string `json:"name"`
	}{}

	err := r.database.Client.GetData("UserProfile", userKey, userFullname)
	if err != nil {
		return "", err
	}

	return userFullname.Name, nil
}
