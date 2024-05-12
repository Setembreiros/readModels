package userprofile

import database "readmodels/infrastructure/db"

type UserProfileRepository database.Database

type UserProfileKey struct {
	Username string
}

func (r UserProfileRepository) AddNewUserProfile(data *UserProfile) error {
	return r.Client.InsertData("UserProfile", data)
}

func (r UserProfileRepository) GetUserProfile(username string) (*UserProfile, error) {
	userProfileKey := &UserProfileKey{
		Username: username,
	}
	var userProfile UserProfile
	err := r.Client.GetData("UserProfile", userProfileKey, &userProfile)
	if err != nil {
		return &userProfile, err
	}

	return &userProfile, nil
}
