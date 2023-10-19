package api

import (
	"net/http"

	"wallemon/pkg/models" // Replace with your module's import path

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GetGemReq struct {
	Address string `json:"address"`
}

// GetGem - Handler to get the gem of a user
func GetGem(c *gin.Context) {
	var gq GetGemReq

	if err := c.ShouldBindJSON(&gq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var gem models.Gem

	db := c.MustGet("db").(*gorm.DB) // Assuming you've set up a middleware to inject the db instance into the context
	if result := db.First(&gem, gq.Address); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch gem."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"gem": gem.Amount})
}
