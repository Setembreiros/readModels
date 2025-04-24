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
	Comments    int       `json:"comments"`
	Likes       int       `json:"likes"`
	Superlikes  int       `json:"superlikes"`
	CreatedAt   time.Time `json:"created_at"`
	LastUpdated time.Time `json:"last_updated"`
}

type CommentKey struct {
	CommentId uint64
}

type PostLikeKey struct {
	PostId   string
	Username string
}

type PostLikeMetadata struct {
	PostId   string
	Username string
	Name     string
}

type PostSuperlikeKey struct {
	PostId   string
	Username string
}

type PostSuperlikeMetadata struct {
	PostId   string
	Username string
	Name     string
}
