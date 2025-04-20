package model

import "time"

type Comment struct {
	CommentId uint64    `json:"commentId"`
	PostId    string    `json:"postId"`
	Username  string    `json:"username"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
