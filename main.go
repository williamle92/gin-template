package main

import (
	"gin-template/database"
	"gin-template/routes"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	// Initialize the database connection by calling the Connect function from the database package.
	err := database.Connect()
	if err != nil {
		// If the database connection fails, log a fatal error and exit the application.
		// `log.Fatalf` prints the error message and stops the program (os.Exit(1)).
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	router := gin.Default()

	// Group API routes under the "/api" prefix for better organization.
	// All routes defined within this group will have "/api" prepended to their path.
	api := router.Group("/api")
	{
		// Register a POST route "/api/users" that maps to the CreateUser handler function.
		// This endpoint is used for registering new users.
		api.POST("/users", routes.CreateUser)

		// Define a POST route "/api/login" that maps to the LoginUser handler function.
		// This endpoint is used for authenticating existing users.
		api.POST("/login", routes.LoginUser)
	}

	if err := router.Run(":8080"); err != nil {
		// If the server fails to start (e.g., port already in use), log the error and panic.
		// `panic` stops the normal execution flow and begins panicking, which usually prints a stack trace.
		log.Fatalf("Failed to start server: %v", err)
		panic(err)
	}
}
