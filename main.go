package main

import (
	"log"
	"shines/models"
	"shines/routers"
)

func main() {

	router := routers.AddRouter()
	models.ConnectToDatabase()
	
	err := router.Run("localhost:8080")
	if err != nil {
		log.Fatal(err)
	}

}
