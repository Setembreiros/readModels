package userprofile

import database "readmodels/infrastructure/db"

type UserProfileRepository database.Database

func (r UserProfileRepository) AddNewUserProfile(data *UserProfile) error {
	return r.Client.InsertData("UserProfile", data)
}

func (r UserProfileRepository) GetUserProfile(username string) (UserProfile, error) {
	return UserProfile{}, nil
}
