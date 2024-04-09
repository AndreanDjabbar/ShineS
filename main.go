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

	shines := router.Group("shines/")
	routers.MainRouter(shines)
	router.GET("/", controllers.RootHandler)

	err := router.Run("localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

}
