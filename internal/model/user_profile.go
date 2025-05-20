package model

type UserProfile struct {
	Username        string `json:"username"`
	Name            string `json:"name"`
	Bio             string `json:"bio"`
	Link            string `json:"link"`
	FollowersAmount int    `json:"followersAmount"`
	FolloweesAmount int    `json:"followeesAmount"`
	PostsAmount     int    `json:"postsAmount"`
}
