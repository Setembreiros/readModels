package userprofile

import database "readmodels/internal/db"

type UserProfileRepository database.Database

func (r UserProfileRepository) GetUserProfile(username string) (*UserProfile, error) {
	userProfileKey := &database.UserProfileKey{
		Username: username,
	}
	var userProfile UserProfile
	err := r.Client.GetData("UserProfile", userProfileKey, &userProfile)

	return &userProfile, err
}

func (r UserProfileRepository) AddNewUserProfile(data *UserProfile) error {
	return r.Client.InsertData("UserProfile", data)
}

func (r UserProfileRepository) UpdateUserProfile(data *UserProfile) error {
	return r.Client.InsertData("UserProfile", data)
}
