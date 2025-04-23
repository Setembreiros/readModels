package model

type PostLike struct {
	PostId string        `json:"postId"`
	User   *UserMetadata `json:"userMetadata"`
}
