package post

import (
	"time"

	"github.com/rs/zerolog/log"
)

//go:generate mockgen -source=service.go -destination=mock/service.go

type Repository interface {
	AddNewPostMetadata(data *PostMetadata) error
	GetPostMetadatasByUser(username string) ([]*PostMetadata, error)
	RemovePostMetadata(postIds []string) error
}

type PostService struct {
	repository Repository
}

type PostMetadata struct {
	PostId      string    `json:"post_id"`
	Username    string    `json:"username"`
	Type        string    `json:"type"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	LastUpdated time.Time `json:"last_updated"`
}

func NewPostService(repository Repository) *PostService {
	return &PostService{
		repository: repository,
	}
}

func (s *PostService) CreateNewPostMetadata(data *PostMetadata) {
	err := s.repository.AddNewPostMetadata(data)
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Error adding post metadata for id %s", data.PostId)
		return
	}

	log.Info().Msgf("Post metadata for id %s was added", data.PostId)
}

func (s *PostService) GetPostMetadatasByUser(username string) ([]*PostMetadata, error) {
	postMetadatas, err := s.repository.GetPostMetadatasByUser(username)
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Error getting post metadatas for username %s", username)
		return postMetadatas, err
	}

	return postMetadatas, nil
}

func (s *PostService) RemovePostMetadata(postIds []string) {
	err := s.repository.RemovePostMetadata(postIds)
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Error removing post metadatas for id %v", postIds)
		return
	}

	log.Info().Msgf("Post metadatas for ids %v were removed", postIds)
}
