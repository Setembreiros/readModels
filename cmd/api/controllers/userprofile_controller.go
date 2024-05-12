package controllers

import (
	"log"
	"net/http"
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
	controller.infoLog.Println("Handling Request")
	id := c.Param("username")
	username := string(id)

	userProfile, err := controller.service.GetUserProfile(username)
	if err != nil {
		controller.errorLog.Println(err)
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	sendOKWithResult(c, userProfile)
}
