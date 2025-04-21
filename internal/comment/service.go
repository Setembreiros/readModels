package comment

import (
	"readmodels/internal/model"

	"github.com/rs/zerolog/log"
)

//go:generate mockgen -source=service.go -destination=test/mock/service.go

type Repository interface {
	CreateComment(data *model.Comment) error
	GetCommentsByPostId(postId string, lastCommentId uint64, limit int) ([]*model.Comment, uint64, error)
	UpdateComment(data *model.Comment) error
	DeleteComment(postId string, commentId uint64) error
}

type CommentService struct {
	repository Repository
}

func NewCommentService(repository Repository) *CommentService {
	return &CommentService{
		repository: repository,
	}
}

func (s *CommentService) CreateComment(data *model.Comment) {
	err := s.repository.CreateComment(data)
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Error creating comment with id %d in post %s", data.CommentId, data.PostId)
		return
	}

	log.Info().Msgf("Comment with id %d in post %s was created", data.CommentId, data.PostId)
}

func (s *CommentService) GetCommentsByPostId(postId string, lastCommentId uint64, limit int) ([]*model.Comment, uint64, error) {
	comments, lastCommentId, err := s.repository.GetCommentsByPostId(postId, lastCommentId, limit)
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Error getting  %s's comments", postId)
		return comments, lastCommentId, err
	}

	return comments, lastCommentId, nil
}

func (s *CommentService) UpdateComment(data *model.Comment) {
	err := s.repository.UpdateComment(data)
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Error updating comment with id %d", data.CommentId)
		return
	}

	log.Info().Msgf("Comment with id %d was updated", data.CommentId)
}

func (s *CommentService) DeleteComment(postId string, commentId uint64) {
	err := s.repository.DeleteComment(postId, commentId)
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Error deleting comment with id %d in post %s", commentId, postId)
		return
	}

	log.Info().Msgf("Comment with id %d in post %s was deleted", commentId, postId)
}
