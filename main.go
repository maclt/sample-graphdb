package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"maclt/graphdb/neo4j/database"
	"maclt/graphdb/neo4j/user"
)

func main() {
	c := context.Background()

	// Initialize database driver
	dbDriver := database.Connect()
	defer dbDriver.Close(c)

	userService := &user.UserService{DatabaseDriver: dbDriver}

	// Create a Gin router
	router := gin.Default()

	// Define a GET route
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to Gin Web Framework!",
		})
	})

	// Define a POST route
	router.POST("/users", userService.RegisterUser)

	// Define a GET route
	router.GET("/users/:name", userService.GetUser)

	// Define a DELETE route
	router.DELETE("/users:name", userService.DeleteUser)

	// Define a POST route
	router.POST("/marriage", userService.MarryUser)

	// Start the server
	router.Run(":8080") // Default runs on localhost:8080
}
