package userprofile

import (
	database "readmodels/internal/db"
	"readmodels/internal/model"
)

type UserProfileRepository database.Database

func (r UserProfileRepository) GetUserProfile(username string) (*model.UserProfile, error) {
	userProfileKey := &database.UserProfileKey{
		Username: username,
	}
	var userProfile model.UserProfile
	err := r.Client.GetData("UserProfile", userProfileKey, &userProfile)

	return &userProfile, err
}

func (r UserProfileRepository) AddNewUserProfile(data *model.UserProfile) error {
	return r.Client.InsertData("UserProfile", data)
}

func (r UserProfileRepository) UpdateUserProfile(data *model.UserProfile) error {
	userProfileKey := &database.UserProfileKey{
		Username: data.Username,
	}

	updateAttributes := map[string]interface{}{
		"Name": data.Name,
		"Bio":  data.Bio,
		"Link": data.Link,
	}

	return r.Client.UpdateData("UserProfile", userProfileKey, updateAttributes)
}

func (r UserProfileRepository) IncreaseFollowers(username string) error {
	userProfileKey := &database.UserProfileKey{
		Username: username,
	}

	return r.Client.IncrementCounter("UserProfile", userProfileKey, "Followers", 1)
}

func (r UserProfileRepository) IncreaseFollowees(username string) error {
	userProfileKey := &database.UserProfileKey{
		Username: username,
	}

	return r.Client.IncrementCounter("UserProfile", userProfileKey, "Followees", 1)
}

func (r UserProfileRepository) DecreaseFollowers(username string) error {
	userProfileKey := &database.UserProfileKey{
		Username: username,
	}

	return r.Client.IncrementCounter("UserProfile", userProfileKey, "Followers", -1)
}

func (r UserProfileRepository) DecreaseFollowees(username string) error {
	userProfileKey := &database.UserProfileKey{
		Username: username,
	}

	return r.Client.IncrementCounter("UserProfile", userProfileKey, "Followees", -1)
}
