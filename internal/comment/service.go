package comment

import (
	"time"

	"github.com/rs/zerolog/log"
)

//go:generate mockgen -source=service.go -destination=test/mock/service.go

type Repository interface {
	AddNewComment(data *Comment) error
}

type CommentService struct {
	repository Repository
}

type Comment struct {
	CommentId string    `json:"commentId"`
	Username  string    `json:"username"`
	PostId    string    `json:"postId"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
}

func NewCommentService(repository Repository) *CommentService {
	return &CommentService{
		repository: repository,
	}
}

func (s *CommentService) CreateNewComment(data *Comment) {
	err := s.repository.AddNewComment(data)
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Error adding comment with id %s", data.CommentId)
		return
	}

	log.Info().Msgf("Comment with id %s was added", data.CommentId)
}
