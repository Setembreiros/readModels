package database

import "context"

func (db Database) ApplyMigrations(ctx context.Context) error {
	db.infoLog.Println("Applying migrations...")

	if !db.Client.TableExists("UserProfile") {
		keys := []TableAttributes{
			{
				Name:          "Username",
				AttributeType: "string",
			},
		}
		err := db.Client.CreateTable("UserProfile", keys, ctx)
		if err != nil {
			return err
		}
	}

	return nil
}
