package database

func (db Database) ApplyMigrations() error {
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
		err := db.Client.CreateTable("UserProfile", attributes)
		if err != nil {
			return err
		}
	}

	return nil
}
