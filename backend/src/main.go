package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Define a type Greeting to manage the JSON response
type Greeting struct {
	Message string `json:"message"`
}

// Main function
func main() {
	r := gin.Default()

	// Define a route for the greeting
	r.GET("/greet/:name", func(c *gin.Context) {
		name := c.Param("name")
		greetingMessage := "Hello, " + name + "!"

		// Respond with a JSON
		c.JSON(http.StatusOK, Greeting{Message: greetingMessage})
	})

	// Start the server on port 8080
	r.Run(":8080")
}
