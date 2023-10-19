package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ClaimReq struct {
	Address string `json:"address"` // EVM address with 0x prefix
	RefCode string `json:"refCode"` // Referral code
}

// Claim - Handler to claim a wallemon
func Claim(c *gin.Context) {
	var gq ClaimReq

	if err := c.ShouldBindJSON(&gq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// check 0x prefix, if not present, add it
	if gq.Address[:2] != "0x" {
		gq.Address = "0x" + gq.Address
	}

	// check if address is valid EVM address
	if len(gq.Address) < 42 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid address."})
		return
	}

	// check if refCode is valid, for now, just check if it's 'wallemon'
	if gq.RefCode != "wallemon" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid referral code."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"signature": "0x1234567890abcdef"})
}
