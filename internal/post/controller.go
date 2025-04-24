package post

import (
	"readmodels/internal/api"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

//go:generate mockgen -source=controller.go -destination=mock/controller.go

type Service interface {
	GetPostMetadatasByUser(username string, lastPostId, lastPostCreatedAt string, limit int) ([]*PostMetadata, string, string, error)
}

type PostController struct {
	service Service
}

type GetPostMetadatasResponse struct {
	Posts             []*PostMetadata `json:"posts"`
	Limit             int             `json:"limit"`
	LastPostId        string          `json:"lastPostId"`
	LastPostCreatedAt string          `json:"lastPostCreatedAt"`
}

func NewPostController(service Service) *PostController {
	return &PostController{
		service: service,
	}
}

func (controller *PostController) Routes(routerGroup *gin.RouterGroup) {
	routerGroup.GET("/user-posts/:username", controller.GetPostMetadatasByUser)
}

func (controller *PostController) GetPostMetadatasByUser(c *gin.Context) {
	log.Info().Msg("Handling Request GET UserProfile")
	username := c.Param("username")
	lastPostId := c.DefaultQuery("lastPostId", "")
	lastPostCreatedAt := c.DefaultQuery("lastPostCreatedAt", "")
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "6"))

	if err != nil || limit <= 0 {
		api.SendBadRequest(c, "Invalid pagination parameters, limit has to be greater than 0")
		return
	}

	if (lastPostId != "" && lastPostCreatedAt == "") || (lastPostId == "" && lastPostCreatedAt != "") {
		api.SendBadRequest(c, "Invalid pagination parameters, lastPostId and lastPostCreatedAt both have to have value or both have to be empty")
		return
	}

	postMetadatas, lastPostId, lastPostCreatedAt, err := controller.service.GetPostMetadatasByUser(username, lastPostId, lastPostCreatedAt, limit)
	if err != nil {
		api.SendInternalServerError(c, err.Error())
		return
	}

	api.SendOKWithResult(c, &GetPostMetadatasResponse{
		Posts:             postMetadatas,
		Limit:             limit,
		LastPostId:        lastPostId,
		LastPostCreatedAt: lastPostCreatedAt,
	})
}
