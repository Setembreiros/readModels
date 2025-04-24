package model

type UserProfile struct {
	Username  string `json:"username"`
	Name      string `json:"name"`
	Bio       string `json:"bio"`
	Link      string `json:"link"`
	Followers int    `json:"followers"`
	Followees int    `json:"followees"`
}
