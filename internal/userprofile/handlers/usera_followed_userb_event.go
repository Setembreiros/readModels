package userprofile_handler

import (
	userprofile "readmodels/internal/userprofile"

	"github.com/rs/zerolog/log"
)

type UserAFollowedUserBEvent struct {
	FollowerID string `json:"followerId"`
	FolloweeID string `json:"followeeId"`
}

type UserAFollowedUserBEventService interface {
	IncreaseFollowers(username string)
	IncreaseFollowees(username string)
}

type UserAFollowedUserBEventHandler struct {
	service UserAFollowedUserBEventService
}

func NewUserAFollowedUserBEventHandler(repository userprofile.Repository) *UserAFollowedUserBEventHandler {
	return &UserAFollowedUserBEventHandler{
		service: userprofile.NewUserProfileService(repository),
	}
}

func (handler *UserAFollowedUserBEventHandler) Handle(event []byte) {
	var userAFollowedUserBEvent UserAFollowedUserBEvent
	log.Info().Msg("Handling UserAFollowedUserBEvent")

	err := Decode(event, &userAFollowedUserBEvent)
	if err != nil {
		log.Error().Stack().Err(err).Msg("Invalid event data")
		return
	}

	handler.service.IncreaseFollowers(userAFollowedUserBEvent.FolloweeID)
	handler.service.IncreaseFollowees(userAFollowedUserBEvent.FollowerID)
}
