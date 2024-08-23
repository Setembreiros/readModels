package post_handler

import (
	"readmodels/internal/post"

	"github.com/rs/zerolog/log"
)

type PostsWereDeletedEvent struct {
	PostIds []string `json:"post_ids"`
}

type PostsWereDeletedEventService interface {
	CreateNewPostMetadata(data *post.PostMetadata)
	RemovePostMetadata(postIds []string)
}

type PostsWereDeletedEventHandler struct {
	service PostsWereDeletedEventService
}

func NewPostsWereDeletedEventHandler(repository post.Repository) *PostsWereDeletedEventHandler {
	return &PostsWereDeletedEventHandler{
		service: post.NewPostService(repository),
	}
}

func (handler *PostsWereDeletedEventHandler) Handle(event []byte) {
	var postsWereDeletedEvent PostsWereDeletedEvent
	log.Info().Msg("Handling PostWasCreatedEvent")

	err := Decode(event, &postsWereDeletedEvent)
	if err != nil {
		log.Error().Stack().Err(err).Msg("Invalid event data")
		return
	}

	handler.service.RemovePostMetadata(postsWereDeletedEvent.PostIds)
}
