package comment

import (
	database "readmodels/internal/db"
)

type CommentRepository database.Database

func (r CommentRepository) AddNewComment(data *Comment) error {
	return r.Client.InsertData("readmodels.comments", data)
}
