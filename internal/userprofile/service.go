package userprofile

import "github.com/rs/zerolog/log"

//go:generate mockgen -source=service.go -destination=test/mock/service.go

type Repository interface {
	AddNewUserProfile(data *UserProfile) error
	UpdateUserProfile(data *UserProfile) error
	GetUserProfile(username string) (*UserProfile, error)
	IncreaseFollowers(username string) error
	IncreaseFollowees(username string) error
	DecreaseFollowers(username string) error
	DecreaseFollowees(username string) error
}

type UserProfileService struct {
	repository Repository
}

type UserProfile struct {
	Username  string `json:"username"`
	Name      string `json:"name"`
	Bio       string `json:"bio"`
	Link      string `json:"link"`
	Followers int    `json:"followers"`
	Followees int    `json:"followees"`
}

func NewUserProfileService(repository Repository) *UserProfileService {
	return &UserProfileService{
		repository: repository,
	}
}

func (s *UserProfileService) CreateNewUserProfile(data *UserProfile) {
	err := s.repository.AddNewUserProfile(data)
	if err != nil {
		log.Error().Stack().Err(err).Msg("Error adding user")
		return
	}

	log.Info().Msgf("User Profile for user %s was added", data.Username)
}

func (s *UserProfileService) UpdateUserProfile(data *UserProfile) {
	err := s.repository.UpdateUserProfile(data)
	if err != nil {
		log.Error().Stack().Err(err).Msg("Error updating user")
		return
	}

	log.Info().Msgf("User Profile for user %s was updated", data.Username)
}

func (s *UserProfileService) GetUserProfile(username string) (*UserProfile, error) {
	userprofile, err := s.repository.GetUserProfile(username)
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Error getting userprofile for username %s", username)
		return userprofile, err
	}

	return userprofile, nil
}

func (s *UserProfileService) IncreaseFollowers(username string) {
	err := s.repository.IncreaseFollowers(username)
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Error increasing %s's followers", username)
	}
}

func (s *UserProfileService) IncreaseFollowees(username string) {
	err := s.repository.IncreaseFollowees(username)
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Error increasing %s's followees", username)
	}
}

func (s *UserProfileService) DecreaseFollowers(username string) {
	err := s.repository.DecreaseFollowers(username)
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Error decreasing %s's followers", username)
	}
}

func (s *UserProfileService) DecreaseFollowees(username string) {
	err := s.repository.DecreaseFollowees(username)
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Error decreasing %s's followees", username)
	}
}
