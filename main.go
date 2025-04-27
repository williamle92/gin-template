package main

import (
	"gin-template/database"
	"gin-template/routes"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	// Initialize the database connection.
	// Assuming database.Connect() establishes the connection and sets database.DB
	// You should add error handling here to gracefully exit if the connection fails.
	err := database.Connect() // Assuming Connect returns an error
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	router := gin.Default()
	api := router.Group("/api")
	{
		api.POST("/users", routes.CreateUser)
	}

	if err := router.Run(":8080"); err != nil {
		panic(err)
	} // Listen and serve on 0.0.0.0:8080
}
