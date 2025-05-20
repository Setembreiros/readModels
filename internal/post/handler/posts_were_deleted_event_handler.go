package post_handler

import (
	common_data "readmodels/internal/common/data"
	"readmodels/internal/post"

	"github.com/rs/zerolog/log"
)

type PostsWereDeletedEvent struct {
	Username string   `json:"username"`
	PostIds  []string `json:"postIds"`
}

type PostsWereDeletedEventService interface {
	CreateNewPostMetadata(data *post.PostMetadata)
	RemovePostMetadata(username string, postIds []string)
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

	err := common_data.DeserializeData(event, &postsWereDeletedEvent)
	if err != nil {
		log.Error().Stack().Err(err).Msg("Invalid event data")
		return
	}

	handler.service.RemovePostMetadata(postsWereDeletedEvent.Username, postsWereDeletedEvent.PostIds)
}
