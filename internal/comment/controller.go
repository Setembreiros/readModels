package comment

import (
	"readmodels/internal/api"
	"readmodels/internal/model"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type CommentController struct {
	service *CommentService
}

type GetCommentsResponse struct {
	Comments      []*model.Comment `json:"comments"`
	LastCommentId uint64           `json:"lastCommentId"`
}

func NewCommentController(repository Repository) *CommentController {
	return &CommentController{
		service: NewCommentService(repository),
	}
}

func (controller *CommentController) Routes(routerGroup *gin.RouterGroup) {
	routerGroup.GET("/comments/:postId", controller.GetCommentsByPostId)
}

func (controller *CommentController) GetCommentsByPostId(c *gin.Context) {
	log.Info().Msg("Handling Request GET Comments")
	postId := c.Param("postId")
	if postId == "" {
		api.SendBadRequest(c, "Missing parameter postId")
		return
	}

	lastCommentId, limit, err := getQueryParameters(c)
	if err != nil || limit <= 0 {
		return
	}

	comments, lastCommentId, err := controller.service.GetCommentsByPostId(postId, lastCommentId, limit)
	if err != nil {
		api.SendInternalServerError(c, err.Error())
		return
	}

	api.SendOKWithResult(c, &GetCommentsResponse{
		Comments:      comments,
		LastCommentId: lastCommentId,
	})
}

func getQueryParameters(c *gin.Context) (uint64, int, error) {
	lastCommentId, err := strconv.ParseUint(c.DefaultQuery("lastCommentId", "0"), 10, 64)
	if err != nil {
		api.SendBadRequest(c, "Invalid pagination parameters, lastCommentId must be a positive number")
		return 0, 0, err
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "12"))
	if err != nil || limit <= 0 {
		api.SendBadRequest(c, "Invalid pagination parameters, limit must be greater than 0")
		return 0, 0, err
	}

	return lastCommentId, limit, nil
}
