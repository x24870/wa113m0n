package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"wallemon/pkg/api"
	"wallemon/pkg/database"
	utils "wallemon/pkg/utils"
)

// Define a type Greeting to manage the JSON response
type Greeting struct {
	Message string `json:"message"`
}

var (
	env string
)

func init() {
	err := utils.LoadEnvConfig("config/.env")
	if err != nil {
		panic(fmt.Errorf("failed to load config: %v", err))
	}

	env = os.Getenv("ENV")
	if env != "local" {
		env = "cloud"
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
	if env == "local" {
		r.Run(":8080")
	} else {
		r.Run(":80")
	}
}
