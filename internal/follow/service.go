package follow

import "github.com/rs/zerolog/log"

//go:generate mockgen -source=service.go -destination=test/mock/service.go

type Repository interface {
	GetFollowersMetadata(followerIds []string) (*[]FollowerMetadata, error)
}

type FollowService struct {
	repository Repository
}

type FollowerMetadata struct {
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
