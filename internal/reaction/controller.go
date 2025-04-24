package reaction

import (
	"readmodels/internal/api"
	"readmodels/internal/model"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

//go:generate mockgen -source=controller.go -destination=test/mock/controller.go

type ReactionController struct {
	service ControllerService
}

type ControllerService interface {
	GetPostLikesMetadata(postId, lastUsername string, limit int) ([]*model.UserMetadata, string, error)
	GetPostSuperlikesMetadata(postId, lastUsername string, limit int) ([]*model.UserMetadata, string, error)
}

type GetPostLikesMetadataResponse struct {
	Users        []*model.UserMetadata `json:"postLikes"`
	LastUsername string                `json:"lastUsername"`
}

type GetPostSuperlikesMetadataResponse struct {
	Users        []*model.UserMetadata `json:"postSuperlikes"`
	LastUsername string                `json:"lastUsername"`
}

func NewReactionController(service ControllerService) *ReactionController {
	return &ReactionController{
		service: service,
	}
}

func (controller *ReactionController) Routes(routerGroup *gin.RouterGroup) {
	routerGroup.GET("/postLikes/:postId", controller.GetPostLikesMetadata)
	routerGroup.GET("/postSuperlikes/:postId", controller.GetPostSuperlikesMetadata)
}

func (controller *ReactionController) GetPostLikesMetadata(c *gin.Context) {
	log.Info().Msg("Handling Request GET PostLikes")
	postId := c.Param("postId")
	if postId == "" {
		api.SendBadRequest(c, "Missing parameter postId")
		return
	}

	lastUsername, limit, err := getQueryParameters(c)
	if err != nil || limit <= 0 {
		return
	}

	users, lastUsername, err := controller.service.GetPostLikesMetadata(postId, lastUsername, limit)
	if err != nil {
		api.SendInternalServerError(c, err.Error())
		return
	}

	api.SendOKWithResult(c, &GetPostLikesMetadataResponse{
		Users:        users,
		LastUsername: lastUsername,
	})
}

func (controller *ReactionController) GetPostSuperlikesMetadata(c *gin.Context) {
	log.Info().Msg("Handling Request GET PostSuperlikes")
	postId := c.Param("postId")
	if postId == "" {
		api.SendBadRequest(c, "Missing parameter postId")
		return
	}

	lastUsername, limit, err := getQueryParameters(c)
	if err != nil || limit <= 0 {
		return
	}

	users, lastUsername, err := controller.service.GetPostSuperlikesMetadata(postId, lastUsername, limit)
	if err != nil {
		api.SendInternalServerError(c, err.Error())
		return
	}

	api.SendOKWithResult(c, &GetPostSuperlikesMetadataResponse{
		Users:        users,
		LastUsername: lastUsername,
	})
}

func getQueryParameters(c *gin.Context) (string, int, error) {
	lastUsername := c.DefaultQuery("lastUsername", "")

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "12"))
	if err != nil || limit <= 0 {
		api.SendBadRequest(c, "Invalid pagination parameters, limit must be greater than 0")
		return "", 0, err
	}

	return lastUsername, limit, nil
}
