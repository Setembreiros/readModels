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

	if !db.Client.TableExists("PostMetadata") {
		keys := []TableAttributes{
			{
				Name:          "PostId",
				AttributeType: "string",
			},
		}
		err := db.Client.CreateTable("PostMetadata", &keys, ctx)
		if err != nil {
			return err
		}
		indexes := []TableAttributes{
			{
				Name:          "Username",
				AttributeType: "string",
			},
			{
				Name:          "CreatedAt",
				AttributeType: "string",
			},
		}
		err = db.Client.CreateIndexesOnTable("PostMetadata", "UserIndex", &indexes, ctx)
		if err != nil {
			return err
		}
		indexes = []TableAttributes{
			{
				Name:          "Type",
				AttributeType: "string",
			},
			{
				Name:          "CreatedAt",
				AttributeType: "string",
			},
		}
		err = db.Client.CreateIndexesOnTable("PostMetadata", "TypeIndex", &indexes, ctx)
		if err != nil {
			return err
		}
	}

	if !db.Client.TableExists("readmodels.comments") {
		keys := []TableAttributes{
			{
				Name:          "CommentId",
				AttributeType: "number",
			},
		}
		err := db.Client.CreateTable("readmodels.comments", &keys, ctx)
		if err != nil {
			return err
		}

		indexes := []TableAttributes{
			{
				Name:          "PostId",
				AttributeType: "string",
			},
			{
				Name:          "CommentId",
				AttributeType: "number",
			},
		}
		err = db.Client.CreateIndexesOnTable("readmodels.comments", "PostIdIndex", &indexes, ctx)
		if err != nil {
			return err
		}
	}

	if !db.Client.TableExists("readmodels.likePosts") {
		keys := []TableAttributes{
			{
				Name:          "PostId",
				AttributeType: "string",
			},
			{
				Name:          "Username",
				AttributeType: "string",
			},
		}
		err := db.Client.CreateTable("readmodels.likePosts", &keys, ctx)
		if err != nil {
			return err
		}
	}

	if !db.Client.TableExists("readmodels.superlikePosts") {
		keys := []TableAttributes{
			{
				Name:          "PostId",
				AttributeType: "string",
			},
			{
				Name:          "Username",
				AttributeType: "string",
			},
		}
		err := db.Client.CreateTable("readmodels.superlikePosts", &keys, ctx)
		if err != nil {
			return err
		}
	}

	return nil
}
