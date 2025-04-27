package main

import (
	"gin-template/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	api := router.Group("/api")
	{
		api.POST("/users", routes.CreateUser)
	}

	if err := router.Run(":8080"); err != nil {
		panic(err)
	} // Listen and serve on 0.0.0.0:8080
}
