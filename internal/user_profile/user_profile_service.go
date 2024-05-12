package userprofile

import "log"

type UserProfileRepository interface {
	AddNewUserProfile(data *UserProfile) error
}

type UserProfileService struct {
	infoLog    *log.Logger
	errorLog   *log.Logger
	repository UserProfileRepository
}

type UserProfile struct {
	UserId   string `json:"user_id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Bio      string `json:"bio"`
	Link     string `json:"link"`
}

func NewUserProfileService(infoLog, errorLog *log.Logger, repository UserProfileRepository) *UserProfileService {
	return &UserProfileService{
		infoLog:    infoLog,
		errorLog:   errorLog,
		repository: repository,
	}
}

func (s *UserProfileService) CreateNewUserProfile(data *UserProfile) {
	err := s.repository.AddNewUserProfile(data)
	if err != nil {
		s.errorLog.Printf("Error adding user, err: %s\n", err)
		return
	}

	s.infoLog.Printf("User Profile for user %s was added", data.Username)
}
