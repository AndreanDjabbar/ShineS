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
	router.Static("/images", "./views/images")
	router.Static("/css", "./views/css")
	router.MaxMultipartMemory = 8 << 20
	models.ConnectToDatabase()

	shines := router.Group("shines/")
	routers.MainRouter(shines)
	router.GET("/", controllers.RootHandler)

	err := router.Run("localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

}
