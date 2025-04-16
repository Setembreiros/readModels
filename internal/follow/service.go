package follow

import "github.com/rs/zerolog/log"

//go:generate mockgen -source=service.go -destination=test/mock/service.go

type Repository interface {
	GetFollowersMetadata(followerIds []string) (*[]FollowerMetadata, error)
	GetFolloweesMetadata(followeeIds []string) (*[]FolloweeMetadata, error)
}

type FollowService struct {
	repository Repository
}

type FollowerMetadata struct {
	Username string `json:"username"`
	Name     string `json:"fullname"`
}

type FolloweeMetadata struct {
	Username string `json:"username"`
	Name     string `json:"fullname"`
}

func NewFollowService(repository Repository) *FollowService {
	return &FollowService{
		repository: repository,
	}
}

func (s *FollowService) GetFollowersMetadata(followerIds []string) (*[]FollowerMetadata, error) {
	followersMetadata, err := s.repository.GetFollowersMetadata(followerIds)
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Error retrieving metadata for followerIds %v", followerIds)
		return &[]FollowerMetadata{}, err
	}

	return followersMetadata, nil
}

func (s *FollowService) GetFolloweesMetadata(followeeIds []string) (*[]FolloweeMetadata, error) {
	followeesMetadata, err := s.repository.GetFolloweesMetadata(followeeIds)
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Error retrieving metadata for followeeIds %v", followeeIds)
		return &[]FolloweeMetadata{}, err
	}

	return followeesMetadata, nil
}
