package model

import "time"

type Review struct {
	ReviewId  uint64    `json:"reviewId"`
	PostId    string    `json:"postId"`
	Username  string    `json:"username"`
	Content   string    `json:"content"`
	Rating    int       `json:"rating"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
