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

func (r CommentRepository) AddNewComment(data *model.Comment) error {
	return r.database.Client.InsertData("readmodels.comments", data)
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

func (r CommentRepository) DeleteComment(commentId uint64) error {
	key := &database.CommentKey{
		CommentId: commentId,
	}
	return r.database.Client.RemoveData("readmodels.comments", key)
}
