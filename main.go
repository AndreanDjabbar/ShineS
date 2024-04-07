package main

import (
	"log"
	"shines/routers"

	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()
	router.LoadHTMLGlob("views/html/*.html")
	router.Static("/css", "./views/css")

	shinesRouter := router.Group("shines/")
	routers.MainRouter(shinesRouter)

	err := router.Run("localhost:8080")
	if err != nil {
		log.Fatal(err)
	}

}
