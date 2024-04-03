package main

import (
	"log"
	"shines/routers"
)

func main() {

	router := routers.AddRouter()

	err := router.Run("localhost:8080")
	if err != nil {
		log.Fatal(err)
	}

}
