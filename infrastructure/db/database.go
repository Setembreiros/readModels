package database

import (
	"context"
	userprofile "readmodels/internal/user_profile"
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

func (d *Database) AddNewUserProfile(data *userprofile.UserProfile) error {
	return d.Client.InsertData("UserProfile", data)
}
