package comment

import (
	database "readmodels/internal/db"
)

type CommentRepository database.Database

func (r CommentRepository) AddNewComment(data *Comment) error {
	return r.Client.InsertData("readmodels.comments", data)
}

func (r CommentRepository) DeleteComment(commentId uint64) error {
	key := &database.CommentKey{
		CommentId: commentId,
	}
	return r.Client.RemoveData("readmodels.comments", key)
}
