package database

import (
	"context"
	"readmodels/internal/model"
)

//go:generate mockgen -source=database.go -destination=test/mock/database.go

type TableAttributes struct {
	Name          string
	AttributeType string
}

type Database struct {
	Client DatabaseClient
}

type DatabaseClient interface {
	Clean()
	Truncate()
	TableExists(tableName string) bool
	CreateTable(tableName string, keys *[]TableAttributes, ctx context.Context) error
	CreateIndexesOnTable(tableName, indexName string, inndexes *[]TableAttributes, ctx context.Context) error
	InsertData(tableName string, attributes any) error
	InsertDataAndIncreaseCounter(tableName string, attributes any, counterTableName string, counterKey any, counterFieldName string) error
	GetData(tableName string, key any, result any) error
	GetMultipleData(tableName string, keys []any, results any) error
	GetPostsByIndexUser(username string, currentUsername string, lastPostId, lastPostCreatedAt string, limit int) ([]*PostMetadata, string, string, error)
	GetCommentsByIndexPostId(postID string, lastCommentId uint64, limit int) ([]*model.Comment, uint64, error)
	GetPostLikesByIndexPostId(postID string, lastUsername string, limit int) ([]*model.UserMetadata, string, error)
	GetPostSuperlikesByIndexPostId(postID string, lastUsername string, limit int) ([]*model.UserMetadata, string, error)
	UpdateData(tableName string, key any, updateAttributes map[string]any) error
	IncrementCounter(tableName string, key any, counterFieldName string, incrementValue int) error
	RemoveDataAndDecreaseCounter(tableName string, key any, counterTableName string, counterKey any, counterFieldName string) error
	RemoveMultipleDataAndDecreaseCounter(tableName string, keys []any, counterTableName string, counterKey any, counterFieldName string) error
	RemoveMultipleData(tableName string, keys []any) error
}

func NewDatabase(client DatabaseClient) *Database {
	return &Database{
		Client: client,
	}
}
