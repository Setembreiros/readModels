package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type response struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Content any    `json:"content"`
}

func sendOKWithResult(c *gin.Context, result any) {
	var payload response

	payload.Error = false
	payload.Message = "200 OK"
	payload.Content = result

	c.IndentedJSON(http.StatusOK, payload)
}

func sendFailure(c *gin.Context, httpStatus int, errorMessage string) {
	var payload response

	payload.Error = true
	payload.Message = errorMessage

	c.IndentedJSON(httpStatus, payload)
}

func sendNotFound(c *gin.Context, errorMessage string) {
	sendFailure(c, http.StatusNotFound, errorMessage)
}

func sendInternalServerError(c *gin.Context, errorMessage string) {
	sendFailure(c, http.StatusInternalServerError, errorMessage)
}
