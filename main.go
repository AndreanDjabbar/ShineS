package main

import (
	"log"
	"shines/controllers"
	"shines/models"
	"shines/routers"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("views/html/*.html")
	router.Static("/css", "./views/css")
	models.ConnectToDatabase()

	router.GET("/", controllers.RootHandler)

	shines := router.Group("shines/")
	routers.MainRouter(shines)
	
	err := router.Run("localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

}
