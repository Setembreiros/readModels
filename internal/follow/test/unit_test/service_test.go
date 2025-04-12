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

func TestGetFollowerMetadatasWithService(t *testing.T) {
	setUpService(t)
	followerIds := []string{"USERA", "USERB", "USERC"}
	expectedData := []*follow.FollowerMetadata{
		{
			Username: followerIds[0],
			Fullname: "fullname1",
		},
		{
			Username: followerIds[1],
			Fullname: "fullname2",
		},
		{
			Username: followerIds[2],
			Fullname: "fullname3",
		},
	}
	repository.EXPECT().GetFollowerMetadatas(followerIds).Return(expectedData, nil)

	followService.GetFollowerMetadatas(followerIds)
}

func TestErrorOnGetFollowerMetadatasWithService(t *testing.T) {
	setUpService(t)
	followerIds := []string{"USERA", "USERB", "USERC"}
	expectedData := []*follow.FollowerMetadata{}
	repository.EXPECT().GetFollowerMetadatas(followerIds).Return(expectedData, errors.New("some error"))

	followService.GetFollowerMetadatas(followerIds)

	assert.Contains(t, loggerOutput.String(), fmt.Sprintf("Error retrieving metadata for followerIds %v", followerIds))
}
