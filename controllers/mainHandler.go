package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func ViewHomeHandler(c *gin.Context) {
	context := gin.H {
		"title":"Home",
		// "user":user,
	}
	c.HTML(
		http.StatusOK,
		"home.html",
		context,
	)
}