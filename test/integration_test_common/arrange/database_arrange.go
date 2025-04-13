package integration_test_arrange

import (
	"context"
	"readmodels/cmd/provider"
	database "readmodels/internal/db"
	"readmodels/internal/userprofile"
	"testing"

	"github.com/stretchr/testify/assert"
)

func CreateTestDatabase(t *testing.T, ctx context.Context) *database.Database {
	provider := provider.NewProvider("test")
	db, _ := provider.ProvideDb(ctx)
	db.ApplyMigrations(ctx)
	return db
}

func AddUserProfileToDatabase(t *testing.T, db *database.Database, data *userprofile.UserProfile) {
	err := db.Client.InsertData("UserProfile", data)
	assert.Nil(t, err)
	userProfileKey := &database.UserProfileKey{
		Username: data.Username,
	}
	var userProfile userprofile.UserProfile
	err = db.Client.GetData("UserProfile", userProfileKey, &userProfile)
	assert.Nil(t, err)
	assert.Equal(t, userProfile.Username, data.Username)
	assert.Equal(t, userProfile.Name, data.Name)
	assert.Equal(t, userProfile.Bio, data.Bio)
	assert.Equal(t, userProfile.Link, data.Link)
}
