package database

import "context"

func (db Database) ApplyMigrations(ctx context.Context) error {
	if !db.Client.TableExists("UserProfile") {
		attributes := []TableAttributes{
			{
				Name:          "userId",
				AttributeType: "string",
			},
			{
				Name:          "username",
				AttributeType: "string",
			},
		}
		err := db.Client.CreateTable("UserProfile", attributes, ctx)
		if err != nil {
			return err
		}
	}

	return nil
}
