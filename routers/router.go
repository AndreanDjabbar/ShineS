package routers

import (
	"net/http"
	"shines/controllers"
	"shines/middlewares"

	"github.com/gin-gonic/gin"
)

func RootHandler(c *gin.Context) {
	isLogged := middlewares.CheckSession(c)
	if !isLogged {
		c.Redirect(
			http.StatusMovedPermanently,
			"/authentication/register",
		)
		return
	} else {
		c.Redirect(
			http.StatusMovedPermanently,
			"/main/home/",
		)
		return
	}
}

func AddRouter() *gin.Engine {
	router := gin.Default()
	router.Use(middlewares.SetSession())
	router.LoadHTMLGlob("views/html/*.html")
	router.Static("/css", "./views/css")

	router.GET("", RootHandler)

	authRouter := router.Group("authentication/")
	{
		authRouter.GET("login/", controllers.ViewLoginHandler)
		authRouter.GET("register/", controllers.ViewRegisterHandler)	
		authRouter.Use(middlewares.AuthSession())
	}

	// mainRouter := router.Group("main/")
	// {
	
	// }

	return router
}