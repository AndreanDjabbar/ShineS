package controllers

import (
	"net/http"
	"shines/middlewares"
	_ "shines/middlewares"

	"github.com/gin-gonic/gin"
)

func ViewHomeHandler(c *gin.Context) {
	user := middlewares.GetSession(c)
	context := gin.H {
		"title":"Home",
		"user":user,
	}
	c.HTML(
		http.StatusOK,
		"home.html",
		context,
	)
}