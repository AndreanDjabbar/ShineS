package main

import (
	"log"
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
	
	err := router.Run("localhost:8080")
	if err != nil {
		log.Fatal(err)
	}

}
