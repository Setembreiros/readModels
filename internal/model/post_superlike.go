package model

type PostSuperlike struct {
	PostId string        `json:"postId"`
	User   *UserMetadata `json:"user"`
}
