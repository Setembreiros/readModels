package database

import (
	"time"
)

type PostMetadata struct {
	PostId      string    `json:"post_id"`
	Username    string    `json:"username"`
	Type        string    `json:"type"`
	FileType    string    `json:"file_type"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	LastUpdated time.Time `json:"last_updated"`
}
