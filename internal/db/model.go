package database

import (
	"time"
)

type UserProfileKey struct {
	Username string
}

type PostMetadataKey struct {
	PostId string
}

type PostMetadata struct {
	PostId      string    `json:"post_id"`
	Username    string    `json:"username"`
	Type        string    `json:"type"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	LastUpdated time.Time `json:"last_updated"`
}

type CommentKey struct {
	CommentId uint64
}

type Comment struct {
	CommentId uint64    `json:"commentId"`
	PostId    string    `json:"postId"`
	Username  string    `json:"username"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
}
