package database

import (
	"context"
	"log"
)

type TableAttributes struct {
	Name          string
	AttributeType string
}

type Database struct {
	infoLog *log.Logger
	Client  DatabaseClient
}

type DatabaseClient interface {
	TableExists(tableName string) bool
	CreateTable(tableName string, attributes []TableAttributes, ctx context.Context) error
	InsertData(tableName string, attributes any) error
}

func NewDatabase(client DatabaseClient, infoLog *log.Logger) *Database {
	return &Database{
		Client:  client,
		infoLog: infoLog,
	}
}
