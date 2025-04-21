package comment

import (
	database "readmodels/internal/db"
	"readmodels/internal/model"
)

type CommentRepository struct {
	cache    *database.Cache
	database *database.Database
}

func NewCommentRepository(database *database.Database, cache *database.Cache) *CommentRepository {
	return &CommentRepository{
		cache:    cache,
		database: database,
	}
}

func (r CommentRepository) CreateComment(data *model.Comment) error {
	postKey := &database.PostMetadataKey{
		PostId: data.PostId,
	}
	return r.database.Client.InsertDataAndIncreaseCounter("readmodels.comments", data, "PostMetadata", postKey, "Comments")
}

func (r CommentRepository) GetCommentsByPostId(postId string, lastCommentId uint64, limit int) ([]*model.Comment, uint64, error) {
	comments, newLastCommentId, found := r.cache.Client.GetPostComments(postId, lastCommentId, limit)
	if found {
		return comments, newLastCommentId, nil
	}

	comments, newLastCommentId, err := r.database.Client.GetCommentsByIndexPostId(postId, lastCommentId, limit)
	if err != nil {
		return []*model.Comment{}, uint64(0), err
	}

	r.cache.Client.SetPostComments(postId, lastCommentId, limit, comments)

	return comments, newLastCommentId, nil
}

type postIdFromCommentData struct {
	PostId string `json:"postId"`
}

func (r CommentRepository) GetPostIdFromComment(commentId uint64) (string, error) {
	commentKey := &database.CommentKey{
		CommentId: commentId,
	}

	data := &postIdFromCommentData{}
	err := r.database.Client.GetData("readmodels.comments", commentKey, data)
	if err != nil {
		return "", err
	}

	return data.PostId, nil
}

func (r CommentRepository) UpdateComment(data *model.Comment) error {
	commentKey := &database.CommentKey{
		CommentId: data.CommentId,
	}

	updateAttributes := map[string]interface{}{
		"Content":   data.Content,
		"UpdatedAt": data.UpdatedAt,
	}

	return r.database.Client.UpdateData("readmodels.comments", commentKey, updateAttributes)
}

func (r CommentRepository) DeleteComment(postId string, commentId uint64) error {
	commentKey := &database.CommentKey{
		CommentId: commentId,
	}

	postKey := &database.PostMetadataKey{
		PostId: postId,
	}

	return r.database.Client.RemoveDataAndDecreaseCounter("readmodels.comments", commentKey, "PostMetadata", postKey, "Comments")
}
