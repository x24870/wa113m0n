package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"wallemon/pkg/api"
	"wallemon/pkg/database"
	utils "wallemon/pkg/utils"
)

// Define a type Greeting to manage the JSON response
type Greeting struct {
	Message string `json:"message"`
}

func init() {
	err := utils.LoadSecrets("config/.secrets")
	if err != nil {
		panic(fmt.Errorf("failed to load secrets: %v", err))
	}
}

// Main function
func main() {
	// Create root context.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	database.Initialize(ctx)
	defer database.Finalize()

	r := gin.Default()

	// Define a route for the greeting
	r.GET("/greet/:name", func(c *gin.Context) {
		name := c.Param("name")
		greetingMessage := "Hello, " + name + "!"

		// Respond with a JSON
		c.JSON(http.StatusOK, Greeting{Message: greetingMessage})
	})

	// Setup the routes
	api.SetupRoutes(r)

	// Start the server on port 8080
	r.Run(":8080")
	// r.Run(":80")
}
