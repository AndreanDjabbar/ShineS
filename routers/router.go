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

func MainRouter(c *gin.RouterGroup) {
	main := c.Group("main/", middlewares.SetSession())
	main.Use(middlewares.AuthSession())

	{
		main.GET("register/", controllers.ViewRegisterHandler)
		main.POST("register/", controllers.RegisterHandler)
		main.GET("login/", controllers.ViewLoginHandler)
		main.POST("login/", controllers.LoginHandler)
	}
	main.GET("home/", controllers.ViewHomeHandler)
}