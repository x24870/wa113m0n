package api

import (
	hdlr "wallemon/pkg/api/handlers"

	middleware "wallemon/pkg/api/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// Configuring and using CORS middleware
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true // adjust this to your needs
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization", "address", "token_id"}
	r.Use(cors.New(config))

	// health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// Register the claim handler
	r.POST("/claim", middleware.MaxBodySize(1024), hdlr.Claim)
	r.POST("/joinWaitlist", middleware.MaxBodySize(1024), hdlr.JoinWaitlist)

	// Register the play handlers
	r.GET("/gem", hdlr.GetGem)
	r.GET("/play", hdlr.GetPlay)
	r.POST("/play", middleware.MaxBodySize(1024), hdlr.Play)
	r.GET("/poop", hdlr.GetPoop)
	r.GET("/clean", hdlr.GetClean)
	r.POST("/clean", middleware.MaxBodySize(1024), hdlr.Clean)
}
