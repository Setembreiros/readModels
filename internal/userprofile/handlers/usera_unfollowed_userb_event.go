package userprofile_handler

import (
	common_data "readmodels/internal/common/data"
	userprofile "readmodels/internal/userprofile"

	"github.com/rs/zerolog/log"
)

type UserAUnfollowedUserBEvent struct {
	FollowerID string `json:"followerId"`
	FolloweeID string `json:"followeeId"`
}

type UserAUnfollowedUserBEventService interface {
	DecreaseFollowers(username string)
	DecreaseFollowees(username string)
}

type UserAUnfollowedUserBEventHandler struct {
	service UserAUnfollowedUserBEventService
}

func NewUserAUnfollowedUserBEventHandler(repository userprofile.Repository) *UserAUnfollowedUserBEventHandler {
	return &UserAUnfollowedUserBEventHandler{
		service: userprofile.NewUserProfileService(repository),
	}
}

func (handler *UserAUnfollowedUserBEventHandler) Handle(event []byte) {
	var userAFollowedUserBEvent UserAUnfollowedUserBEvent
	log.Info().Msg("Handling UserAUnfollowedUserBEvent")

	err := common_data.DeserializeData(event, &userAFollowedUserBEvent)
	if err != nil {
		log.Error().Stack().Err(err).Msg("Invalid event data")
		return
	}

	handler.service.DecreaseFollowers(userAFollowedUserBEvent.FolloweeID)
	handler.service.DecreaseFollowees(userAFollowedUserBEvent.FollowerID)
}
