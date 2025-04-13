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

func (r FollowRepository) GetFolloweesMetadata(followeeIds []string) (*[]FolloweeMetadata, error) {
	followeeKeys := make([]any, len(followeeIds))
	for i, v := range followeeIds {
		followeeKeys[i] = database.UserProfileKey{
			Username: v,
		}
	}

	followeesMetadata := &[]FolloweeMetadata{} // mandatory inizialiting like this otherwise it will failed
	err := r.Client.GetMultipleData("UserProfile", followeeKeys, followeesMetadata)
	if err != nil {
		return followeesMetadata, err
	}

	return followeesMetadata, nil
}
