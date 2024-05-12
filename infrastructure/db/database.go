package database

import (
	"context"
)

type TableAttributes struct {
	Name          string
	AttributeType string
}

type Database struct {
	Client DatabaseClient
}

type DatabaseClient interface {
	TableExists(tableName string) bool
	CreateTable(tableName string, attributes []TableAttributes, ctx context.Context) error
	InsertData(tableName string, attributes any) error
}
