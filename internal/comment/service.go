package comment

import (
	"time"

	"github.com/rs/zerolog/log"
)

//go:generate mockgen -source=service.go -destination=test/mock/service.go

type Repository interface {
	AddNewComment(data *Comment) error
	DeleteComment(commentId uint64) error
}

type CommentService struct {
	repository Repository
}

type Comment struct {
	CommentId uint64    `json:"commentId"`
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
		log.Error().Stack().Err(err).Msgf("Error adding comment with id %d", data.CommentId)
		return
	}

	log.Info().Msgf("Comment with id %d was added", data.CommentId)
}

func (s *CommentService) DeleteComment(commentId uint64) {
	err := s.repository.DeleteComment(commentId)
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Error deleting comment with id %d", commentId)
		return
	}

	log.Info().Msgf("Comment with id %d was deleted", commentId)
}
