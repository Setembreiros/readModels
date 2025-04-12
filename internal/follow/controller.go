package follow

import (
	"readmodels/internal/api"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type FollowController struct {
	service *FollowService
}

type GetFollowerMetadatasResponse struct {
	Followers []*FollowerMetadata `json:"followers"`
}

func NewFollowController(repository Repository) *FollowController {
	return &FollowController{
		service: NewFollowService(repository),
	}
}

func (controller *FollowController) Routes(routerGroup *gin.RouterGroup) {
	routerGroup.GET("/followers", controller.GetFollowerMetadatas)
}

func (controller *FollowController) GetFollowerMetadatas(c *gin.Context) {
	log.Info().Msg("Handling Request GET Followers")
	followerIds := c.QueryArray("followerId")

	followerMetadatas, err := controller.service.GetFollowerMetadatas(followerIds)
	if err != nil {
		api.SendInternalServerError(c, err.Error())
		return
	}

	api.SendOKWithResult(c, &GetFollowerMetadatasResponse{
		Followers: followerMetadatas,
	})
}
