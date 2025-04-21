package reaction

import (
	"readmodels/internal/model"

	"github.com/rs/zerolog/log"
)

//go:generate mockgen -source=service.go -destination=test/mock/service.go

type Repository interface {
	CreateLikePost(data *model.LikePost) error
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
		log.Error().Stack().Err(err).Msgf("Error creating like, username: %s -> postId: %s", data.Username, data.PostId)
		return
	}

	log.Info().Msgf("Like was created, username: %s -> postId: %s", data.Username, data.PostId)
}
