package controllers

import (
	"errors"
	"log"
	database "readmodels/infrastructure/db"
	"readmodels/internal/userprofile"

	"github.com/gin-gonic/gin"
)

type UserProfileController struct {
	infoLog  *log.Logger
	errorLog *log.Logger
	service  *userprofile.UserProfileService
}

func NewUserProfileController(infoLog, errorLog *log.Logger, repository userprofile.Repository) *UserProfileController {
	return &UserProfileController{
		service:  userprofile.NewUserProfileService(infoLog, errorLog, repository),
		infoLog:  infoLog,
		errorLog: errorLog,
	}
}

func (controller *UserProfileController) Routes(router *gin.Engine) {
	router.GET("/userprofile/:username", controller.getUserProfile)
}

func (controller *UserProfileController) getUserProfile(c *gin.Context) {
	controller.infoLog.Println("Handling Request GET UserProfile")
	id := c.Param("username")
	username := string(id)

	userProfile, err := controller.service.GetUserProfile(username)
	if err != nil {
		var notFoundError database.NotFoundError
		if errors.Is(err, notFoundError) {
			message := "User Profile not found for username " + username
			sendNotFound(c, message)
		} else {
			sendInternalServerError(c, err.Error())
		}
		return
	}

	sendOKWithResult(c, userProfile)
}
