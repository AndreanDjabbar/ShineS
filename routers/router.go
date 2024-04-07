package routers

import (
	"net/http"
	"shines/controllers"
	"github.com/gin-gonic/gin"
)

func RootHandler(c *gin.Context) {
	c.Redirect(
		http.StatusMovedPermanently,
		"/main/home/",
	)
}

func AddRouter() *gin.Engine {
	router := gin.Default()
	router.LoadHTMLGlob("views/html/*.html")
	router.Static("/css", "./views/css")

	router.GET("", RootHandler)

	authRouter := router.Group("/authentication/")
	{
		authRouter.GET("login/", controllers.ViewLoginHandler)
		authRouter.POST("login/", controllers.LoginHandler)
		authRouter.GET("register/", controllers.ViewRegisterHandler)	
		authRouter.POST("register/", controllers.RegisterHandler)	
		authRouter.GET("logoutss", controllers.LogoutHandler)
	}

	mainRouter := router.Group("/main/")
	{
		mainRouter.GET("home/", controllers.ViewHomeHandler)
	}

	return router
}