package userprofile

import (
	"errors"
	"readmodels/internal/api"
	database "readmodels/internal/db"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type UserProfileController struct {
	service *UserProfileService
}

func NewUserProfileController(repository Repository) *UserProfileController {
	return &UserProfileController{
		service: NewUserProfileService(repository),
	}
}

func (controller *UserProfileController) Routes(routerGroup *gin.RouterGroup) {
	routerGroup.GET("/userprofile/:username", controller.GetUserProfile)
}

func (controller *UserProfileController) GetUserProfile(c *gin.Context) {
	log.Info().Msg("Handling Request GET UserProfile")
	id := c.Param("username")
	username := string(id)

	userProfile, err := controller.service.GetUserProfile(username)
	if err != nil {
		var notFoundError *database.NotFoundError
		if errors.As(err, &notFoundError) {
			message := "User Profile not found for username " + username
			api.SendNotFound(c, message)
		} else {
			api.SendInternalServerError(c, err.Error())
		}
		return
	}

	api.SendOKWithResult(c, &userProfile)
}
