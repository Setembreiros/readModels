package database

import (
	"context"

	"github.com/rs/zerolog/log"
)

func (db *Database) ApplyMigrations(ctx context.Context) error {
	log.Info().Msg("Applying migrations...")

	if !db.Client.TableExists("UserProfile") {
		keys := []TableAttributes{
			{
				Name:          "Username",
				AttributeType: "string",
			},
		}
		err := db.Client.CreateTable("UserProfile", &keys, ctx)
		if err != nil {
			return err
		}
	}

	return nil
}
