package routers

import (
	"shines/controllers"
	"shines/middlewares"
	"github.com/gin-gonic/gin"
)


func MainRouter(c *gin.RouterGroup) {
	main := c.Group("main/", middlewares.SetSession())
	{
		main.GET("", controllers.RootHandler)
		main.GET("login-page/", controllers.ViewLoginHandler)
		main.POST("login-page/", controllers.LoginHandler)
		main.GET("register-page/", controllers.ViewRegisterHandler)	
		main.POST("register-page/", controllers.RegisterHandler)	
		main.GET("logout-page/", controllers.LogoutHandler)
		main.Use(middlewares.AuthSession())
		main.GET("home-page/", controllers.ViewHomeHandler)
		main.GET("personal-information-page/", controllers.ViewPersonalHandler)
	}
}
