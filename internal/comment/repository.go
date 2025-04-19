package comment

import (
	database "readmodels/internal/db"
)

type CommentRepository database.Database

func (r CommentRepository) AddNewComment(data *Comment) error {
	return r.Client.InsertData("readmodels.comments", data)
}

func (r CommentRepository) GetCommentsByPostId(postId string, lastCommentId uint64, limit int) ([]*Comment, uint64, error) {
	data, lastCommentId, err := r.Client.GetCommentsByIndexPostId(postId, lastCommentId, limit)
	if err != nil {
		return []*Comment{}, uint64(0), err
	}

	var comments []*Comment
	for _, comment := range data {
		comments = append(comments, mapToDomain(comment))
	}

	return comments, lastCommentId, nil
}

func (r CommentRepository) DeleteComment(commentId uint64) error {
	key := &database.CommentKey{
		CommentId: commentId,
	}
	return r.Client.RemoveData("readmodels.comments", key)
}

func mapToDomain(data *database.Comment) *Comment {
	return &Comment{
		CommentId: data.CommentId,
		PostId:    data.PostId,
		Username:  data.Username,
		Content:   data.Content,
		CreatedAt: data.CreatedAt,
	}
}
