package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type response struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Content interface{} `json:"content"`
}

func sendOKWithResult(c *gin.Context, result interface{}) {
	var payload response

	payload.Error = false
	payload.Message = "200 OK"
	payload.Content = result

	c.IndentedJSON(http.StatusOK, payload)
}
