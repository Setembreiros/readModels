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
	GetLikesMetadataByPostId(postId, lastUsername string, limit int) ([]*model.UserMetadata, string, error)
	GetSuperlikesMetadataByPostId(postId, lastUsername string, limit int) ([]*model.UserMetadata, string, error)
	GetReviewsByPostId(postId string, lastReviewId uint64, limit int) ([]*model.Review, uint64, error)
}

type GetPostLikesMetadataResponse struct {
	Users        []*model.UserMetadata `json:"postLikes"`
	LastUsername string                `json:"lastUsername"`
}

type GetPostSuperlikesMetadataResponse struct {
	Users        []*model.UserMetadata `json:"postSuperlikes"`
	LastUsername string                `json:"lastUsername"`
}

type GetReviewsResponse struct {
	Reviews      []*model.Review `json:"reviews"`
	LastReviewId uint64          `json:"lastReviewId"`
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

	users, lastUsername, err := controller.service.GetLikesMetadataByPostId(postId, lastUsername, limit)
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

	users, lastUsername, err := controller.service.GetSuperlikesMetadataByPostId(postId, lastUsername, limit)
	if err != nil {
		api.SendInternalServerError(c, err.Error())
		return
	}

	api.SendOKWithResult(c, &GetPostSuperlikesMetadataResponse{
		Users:        users,
		LastUsername: lastUsername,
	})
}

func (controller *ReactionController) GetReviewsByPostId(c *gin.Context) {
	log.Info().Msg("Handling Request GET Reviews")
	postId := c.Param("postId")
	if postId == "" {
		api.SendBadRequest(c, "Missing parameter postId")
		return
	}

	lastReviewId, err := strconv.ParseUint(c.DefaultQuery("lastReviewId", "0"), 10, 64)
	if err != nil {
		api.SendBadRequest(c, "Invalid pagination parameters, lastReviewId must be a positive number")
		return
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "12"))
	if err != nil || limit <= 0 {
		api.SendBadRequest(c, "Invalid pagination parameters, limit must be greater than 0")
		return
	}

	reviews, lastReviewId, err := controller.service.GetReviewsByPostId(postId, lastReviewId, limit)
	if err != nil {
		api.SendInternalServerError(c, err.Error())
		return
	}

	api.SendOKWithResult(c, &GetReviewsResponse{
		Reviews:      reviews,
		LastReviewId: lastReviewId,
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
