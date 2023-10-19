package api

import (
	hdlr "wallemon/pkg/api/handlers"

	middleware "wallemon/pkg/api/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// Register the claim handler
	r.POST("/claim", middleware.MaxBodySize(1024), hdlr.Claim)

	// Register the gem handlers
	gemGroup := r.Group("/gem")
	{
		gemGroup.GET("/gem", hdlr.GetGem)
	}
}
