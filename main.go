package main

import (
	"log"
	"shines/controllers"

	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()

	router.LoadHTMLGlob("views/html/*.html")

	router.GET("/", controllers.HomePage)

	err := router.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}

}
