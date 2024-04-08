package controllers

import (
	"fmt"
	"net/http"
	"shines/middlewares"

	"github.com/gin-gonic/gin"
)

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