package reaction

import (
	"readmodels/internal/model"

	"github.com/rs/zerolog/log"
)

//go:generate mockgen -source=service.go -destination=test/mock/service.go

type Repository interface {
	CreateLikePost(data *model.LikePost) error
	CreateSuperlikePost(data *model.SuperlikePost) error
	DeleteLikePost(data *model.LikePost) error
	DeleteSuperlikePost(data *model.SuperlikePost) error
}

type ReactionService struct {
	repository Repository
}

func NewReactionService(repository Repository) *ReactionService {
	return &ReactionService{
		repository: repository,
	}
}

func (s *ReactionService) CreateLikePost(data *model.LikePost) {
	err := s.repository.CreateLikePost(data)
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Error creating likePost, username: %s -> postId: %s", data.Username, data.PostId)
		return
	}

	log.Info().Msgf("LikePost was created, username: %s -> postId: %s", data.Username, data.PostId)
}

func (s *ReactionService) CreateSuperlikePost(data *model.SuperlikePost) {
	err := s.repository.CreateSuperlikePost(data)
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Error creating superlikePost, username: %s -> postId: %s", data.Username, data.PostId)
		return
	}

	log.Info().Msgf("SuperlikePost was created, username: %s -> postId: %s", data.Username, data.PostId)
}

func (s *ReactionService) DeleteLikePost(data *model.LikePost) {
	err := s.repository.DeleteLikePost(data)
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Error deleting likePost, username: %s -> postId: %s", data.Username, data.PostId)
		return
	}

	log.Info().Msgf("LikePost was deleted, username: %s -> postId: %s", data.Username, data.PostId)
}

func (s *ReactionService) DeleteSuperlikePost(data *model.SuperlikePost) {
	err := s.repository.DeleteSuperlikePost(data)
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Error deleting superlikePost, username: %s -> postId: %s", data.Username, data.PostId)
		return
	}

	log.Info().Msgf("SuperlikePost was deleted, username: %s -> postId: %s", data.Username, data.PostId)
}
