package follow

import (
	database "readmodels/internal/db"
)

type FollowRepository database.Database

func (r FollowRepository) GetFollowerMetadatas(followerIds []string) (*[]FollowerMetadata, error) {
	followerKeys := make([]any, len(followerIds))
	for i, v := range followerIds {
		followerKeys[i] = database.UserProfileKey{
			Username: v,
		}
	}

	followerMetadatas := &[]FollowerMetadata{} // mandatory inizialiting like this otherwise it will failed
	err := r.Client.GetMultipleData("UserProfile", followerKeys, followerMetadatas)
	if err != nil {
		return followerMetadatas, err
	}

	return followerMetadatas, nil
}
