package database

import (
	"context"
	"log"
)

//go:generate mockgen -source=database.go -destination=mock/database.go

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
	CreateTable(tableName string, keys *[]TableAttributes, ctx context.Context) error
	InsertData(tableName string, attributes any) error
	GetData(tableName string, key any, result any) error
}

func NewDatabase(client DatabaseClient, infoLog *log.Logger) *Database {
	return &Database{
		Client:  client,
		infoLog: infoLog,
	}
}
