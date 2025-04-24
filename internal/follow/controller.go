package follow

import (
	"readmodels/internal/api"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type FollowController struct {
	service *FollowService
}

type GetFollowersMetadataResponse struct {
	Followers *[]FollowerMetadata `json:"followers"`
}

func NewFollowController(repository Repository) *FollowController {
	return &FollowController{
		service: NewFollowService(repository),
	}
}

func (controller *FollowController) Routes(routerGroup *gin.RouterGroup) {
	routerGroup.GET("/followers", controller.GetFollowersMetadata)
}

func (controller *FollowController) GetFollowersMetadata(c *gin.Context) {
	log.Info().Msg("Handling Request GET Followers")
	followerIds := c.QueryArray("followerId")

	followersMetadata, err := controller.service.GetFollowersMetadata(followerIds)
	if err != nil {
		api.SendInternalServerError(c, err.Error())
		return
	}

	api.SendOKWithResult(c, &GetFollowersMetadataResponse{
		Followers: followersMetadata,
	})
}
