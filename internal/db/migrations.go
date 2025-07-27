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

	if !db.Client.TableExists("readmodels.reviews") {
		keys := []TableAttributes{
			{
				Name:          "ReviewId",
				AttributeType: "number",
			},
		}
		err := db.Client.CreateTable("readmodels.reviews", &keys, ctx)
		if err != nil {
			return err
		}

		indexes := []TableAttributes{
			{
				Name:          "PostId",
				AttributeType: "string",
			},
			{
				Name:          "ReviewId",
				AttributeType: "number",
			},
		}
		err = db.Client.CreateIndexesOnTable("readmodels.reviews", "PostIdIndex", &indexes, ctx)
		if err != nil {
			return err
		}
	}

	if !db.Client.TableExists("readmodels.postLikes") {
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
		err := db.Client.CreateTable("readmodels.postLikes", &keys, ctx)
		if err != nil {
			return err
		}
	}

	if !db.Client.TableExists("readmodels.postSuperlikes") {
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
		err := db.Client.CreateTable("readmodels.postSuperlikes", &keys, ctx)
		if err != nil {
			return err
		}
	}

	if db.Client.TableExists("readmodels.reviews") {
		// Comprobar se o índice xa existe antes de crealo
		if !db.Client.IndexExists("readmodels.reviews", "UsernamePostIndex") {
			usernamePostIndexes := []TableAttributes{
				{
					Name:          "Username",
					AttributeType: "string",
				},
				{
					Name:          "PostId",
					AttributeType: "string",
				},
			}
			err := db.Client.CreateIndexesOnTable("readmodels.reviews", "UsernamePostIndex", &usernamePostIndexes, ctx)
			if err != nil {
				log.Error().Err(err).Msg("Error creating UsernamePostIndex on readmodels.reviews")
				return err
			}
			log.Info().Msg("Created UsernamePostIndex on readmodels.reviews table")
		} else {
			log.Info().Msg("UsernamePostIndex already exists on readmodels.reviews table")
		}
	}

	return nil
}
