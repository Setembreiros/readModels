package database

import "context"

func (db Database) ApplyMigrations(ctx context.Context) error {
	db.infoLog.Println("Applying migrations...")

	if !db.Client.TableExists("UserProfile") {
		attributes := []TableAttributes{
			{
				Name:          "UserId",
				AttributeType: "string",
			},
			{
				Name:          "Username",
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
