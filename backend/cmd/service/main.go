package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"wallemon/pkg/api"
	gormpkg "wallemon/pkg/gorm"
	"wallemon/pkg/models"
	utils "wallemon/pkg/utils"
)

// Define a type Greeting to manage the JSON response
type Greeting struct {
	Message string `json:"message"`
}

var (
	gormdb *gorm.DB
)

func init() {
	var err error
	gormdb, err = gormpkg.NewGormPostgresConn(
		gormpkg.Config{
			// DSN:             config.GetDBArg(),
			// DSN:             "postgres://user:user@db:5432/wallemon?sslmode=disable", //TODO: use config
			DSN: "postgres://user:user@db:5432/postgres?sslmode=disable", //TODO: use config
			// DSN:             "host=db port=5432 user=user password=user dbname=postgres sslmode=disable binary_parameters=yes",
			MaxIdleConns:    2,
			MaxOpenConns:    2,
			ConnMaxLifetime: 10 * time.Minute,
			SingularTable:   true,
		},
	)
	if err != nil {
		panic(fmt.Errorf("failed to init db, err: %v", err))
	}

	// Perform database schema auto-migration.
	if err := models.AutoMigrate(gormdb); err != nil {
		panic(err)
	}

	err = utils.LoadSecrets("config/.secrets")
	if err != nil {
		panic(fmt.Errorf("failed to load secrets: %v", err))
	}
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

	// Setup the routes
	api.SetupRoutes(r)

	// Start the server on port 8080
	r.Run(":8080")
	// r.Run(":80")
}
