package controllers

import (
	"net/http"
	_"shines/middlewares"
	"github.com/gin-gonic/gin"
)

func ViewHomeHandler(c *gin.Context) {
	context := gin.H {
		"title":"Home",
	}
	c.HTML(
		http.StatusOK,
		"home.html",
		context,
	)
}