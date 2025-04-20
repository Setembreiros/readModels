package comment

import (
	"readmodels/internal/model"

	"github.com/rs/zerolog/log"
)

//go:generate mockgen -source=service.go -destination=test/mock/service.go

type Repository interface {
	AddNewComment(data *model.Comment) error
	GetCommentsByPostId(postId string, lastCommentId uint64, limit int) ([]*model.Comment, uint64, error)
	DeleteComment(commentId uint64) error
}

type CommentService struct {
	repository Repository
}

func NewCommentService(repository Repository) *CommentService {
	return &CommentService{
		repository: repository,
	}
}

func (s *CommentService) CreateNewComment(data *model.Comment) {
	err := s.repository.AddNewComment(data)
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Error adding comment with id %d", data.CommentId)
		return
	}

	log.Info().Msgf("Comment with id %d was added", data.CommentId)
}

func (s *CommentService) GetCommentsByPostId(postId string, lastCommentId uint64, limit int) ([]*model.Comment, uint64, error) {
	comments, lastCommentId, err := s.repository.GetCommentsByPostId(postId, lastCommentId, limit)
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Error getting  %s's comments", postId)
		return comments, lastCommentId, err
	}

	return comments, lastCommentId, nil
}

func (s *CommentService) DeleteComment(commentId uint64) {
	err := s.repository.DeleteComment(commentId)
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Error deleting comment with id %d", commentId)
		return
	}

	log.Info().Msgf("Comment with id %d was deleted", commentId)
}
