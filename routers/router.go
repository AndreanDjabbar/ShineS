package routers

import (
	"net/http"
	"shines/controllers"
	"shines/middlewares"
	"github.com/gin-gonic/gin"
)

func RootHandler(c *gin.Context) {
	c.Redirect(
		http.StatusMovedPermanently,
		"/main/home/",
	)
}

func MainRouter(c *gin.RouterGroup) {
	main := c.Group("main/", middlewares.SetSession())
	main.Use(middlewares.AuthSession())
	{
		main.GET("login/", controllers.ViewLoginHandler)
		main.POST("login/", controllers.LoginHandler)
		main.GET("register/", controllers.ViewRegisterHandler)	
		main.POST("register/", controllers.RegisterHandler)	
		main.GET("logout/", controllers.LogoutHandler)
	}

	main.GET("home/", controllers.ViewHomeHandler)
}