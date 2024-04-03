package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ViewLoginHandler(c *gin.Context) {
	context := gin.H {
		"title":"Login",
	}
	c.HTML(
		http.StatusOK,
		"login.html",
		context,
	)
}

func ViewRegisterHandler(c *gin.Context) {
	context := gin.H {
		"title":"Sign Up",
	}
	c.HTML(
		http.StatusOK,
		"register.html",
		context,
	)
}
