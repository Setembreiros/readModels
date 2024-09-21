package post

import (
	"readmodels/internal/api"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type PostController struct {
	service *PostService
}

type GetPostMetadatasResponse struct {
	Posts []*PostMetadata `json:"posts"`
}

func NewPostController(repository Repository) *PostController {
	return &PostController{
		service: NewPostService(repository),
	}
}

func (controller *PostController) Routes(routerGroup *gin.RouterGroup) {
	routerGroup.GET("/user-posts/:username", controller.GetPostMetadatasByUser)
}

func (controller *PostController) GetPostMetadatasByUser(c *gin.Context) {
	log.Info().Msg("Handling Request GET UserProfile")
	id := c.Param("username")
	username := string(id)

	postMetadatas, err := controller.service.GetPostMetadatasByUser(username)
	if err != nil {
		api.SendInternalServerError(c, err.Error())
		return
	}

	api.SendOKWithResult(c, &GetPostMetadatasResponse{
		Posts: postMetadatas,
	})
}
