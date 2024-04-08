package controllers

import (
	"fmt"
	"net/http"
	"shines/middlewares"

	"github.com/gin-gonic/gin"
)

func RootHandler(c *gin.Context) {
	isLogged := middlewares.CheckSession(c)
	if isLogged {
		c.Redirect(
			http.StatusFound,
			"/shines/main/home-page",
		)
	} else {
		c.Redirect(
			http.StatusFound,
			"/shines/main/login-page",
		)
	}
}

func ViewHomeHandler(c *gin.Context) {
	user := middlewares.GetSession(c)
	isLogged := middlewares.CheckSession(c)
	fmt.Println(isLogged)
	if !isLogged {
		c.Redirect(
			http.StatusFound,
			"shines/main/login",
		)
		return
	}
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