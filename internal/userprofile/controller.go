package userprofile

import (
	"errors"
	"log"
	"readmodels/internal/api"
	database "readmodels/internal/db"

	"github.com/gin-gonic/gin"
)

type UserProfileController struct {
	infoLog  *log.Logger
	errorLog *log.Logger
	service  *UserProfileService
}

func NewUserProfileController(infoLog, errorLog *log.Logger, repository Repository) *UserProfileController {
	return &UserProfileController{
		service:  NewUserProfileService(infoLog, errorLog, repository),
		infoLog:  infoLog,
		errorLog: errorLog,
	}
}

func (controller *UserProfileController) Routes(router *gin.Engine) {
	router.GET("/userprofile/:username", controller.GetUserProfile)
}

func (controller *UserProfileController) GetUserProfile(c *gin.Context) {
	controller.infoLog.Println("Handling Request GET UserProfile")
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
