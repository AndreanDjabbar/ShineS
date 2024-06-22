package main

import (
	"log"
	"shines/controllers"
	"shines/models"
	"shines/routers"
	"text/template"
	repositories "shines/repositories"
	"github.com/gin-gonic/gin"
)

func init() {
	models.ConnectToDatabase()
}

func main() {
	router := gin.Default()
	router.SetFuncMap(template.FuncMap{
		"add1": repositories.Add1,
	})
	router.LoadHTMLGlob("views/html/*.html")
	router.Static("/images", "./views/images")
	router.Static("/css", "./views/css")
	router.MaxMultipartMemory = 8 << 20
	shines := router.Group("shines/")
	routers.MainRouter(shines)
	router.GET("/", controllers.RootHandler)

	err := router.Run("localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

}
