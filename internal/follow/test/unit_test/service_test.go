package unit_test_follow

import (
	"errors"
	"fmt"
	"readmodels/internal/follow"
	"testing"

	"github.com/stretchr/testify/assert"
)

var followService *follow.FollowService

func setUpService(t *testing.T) {
	setUp(t)
	followService = follow.NewFollowService(repository)
}

func TestGetFollowersMetadataWithService(t *testing.T) {
	setUpService(t)
	followerIds := []string{"USERA", "USERB", "USERC"}
	expectedData := &[]follow.FollowerMetadata{
		{
			Username: followerIds[0],
			Name:     "fullname1",
		},
		{
			Username: followerIds[1],
			Name:     "fullname2",
		},
		{
			Username: followerIds[2],
			Name:     "fullname3",
		},
	}
	repository.EXPECT().GetFollowersMetadata(followerIds).Return(expectedData, nil)

	followService.GetFollowersMetadata(followerIds)
}

func TestErrorOnGetFollowersMetadataWithService(t *testing.T) {
	setUpService(t)
	followerIds := []string{"USERA", "USERB", "USERC"}
	expectedData := &[]follow.FollowerMetadata{}
	repository.EXPECT().GetFollowersMetadata(followerIds).Return(expectedData, errors.New("some error"))

	followService.GetFollowersMetadata(followerIds)

	assert.Contains(t, loggerOutput.String(), fmt.Sprintf("Error retrieving metadata for followerIds %v", followerIds))
}
