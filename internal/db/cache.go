package database

import "readmodels/internal/model"

//go:generate mockgen -source=cache.go -destination=test/mock/cache.go

type Cache struct {
	Client CacheClient
}

type CacheClient interface {
	Clean()
	SetPostComments(postId string, lastCommentId uint64, limit int, comments []*model.Comment)
	GetPostComments(postId string, lastCommentId uint64, limit int) ([]*model.Comment, uint64, bool)
}

func NewCache(client CacheClient) *Cache {
	return &Cache{
		Client: client,
	}
}
