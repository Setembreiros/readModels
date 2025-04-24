package follow

import (
	database "readmodels/internal/db"
)

type FollowRepository database.Database

func (r FollowRepository) GetFollowersMetadata(followerIds []string) (*[]FollowerMetadata, error) {
	followerKeys := make([]any, len(followerIds))
	for i, v := range followerIds {
		followerKeys[i] = database.UserProfileKey{
			Username: v,
		}
	}

	followersMetadata := &[]FollowerMetadata{} // mandatory inizialiting like this otherwise it will failed
	err := r.Client.GetMultipleData("UserProfile", followerKeys, followersMetadata)
	if err != nil {
		return followersMetadata, err
	}

	return followersMetadata, nil
}
