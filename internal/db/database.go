package database

import (
	"context"
)

//go:generate mockgen -source=database.go -destination=mock/database.go

type TableAttributes struct {
	Name          string
	AttributeType string
}

type Database struct {
	Client DatabaseClient
}

type DatabaseClient interface {
	TableExists(tableName string) bool
	CreateTable(tableName string, keys *[]TableAttributes, ctx context.Context) error
	CreateIndexesOnTable(tableName, indexName string, inndexes *[]TableAttributes, ctx context.Context) error
	InsertData(tableName string, attributes any) error
	GetData(tableName string, key any, result any) error
	RemoveMultipleData(tableName string, keys []any) error
	GetPostsByIndexUser(username string, lastPostId, lastPostCreatedAt string, limit int) ([]*PostMetadata, string, string, error)
}

func NewDatabase(client DatabaseClient) *Database {
	return &Database{
		Client: client,
	}
}
