package api

import (
	"net/http"

	"wallemon/pkg/database"
	"wallemon/pkg/models"
	"wallemon/pkg/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const playMessage = "Let's play!"

type GetGemReq struct {
	Address string `json:"address"`
	TokenID uint   `json:"token_id"`
}

// GetGem - Handler to get the gem of a user
func GetGem(c *gin.Context) {
	var req GetGemReq

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: maybe check if this address owns this token

	db := database.GetSQL()
	t := models.NewToken(req.TokenID)
	t, err := t.CreateIfNotExists(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	g := models.NewGem(t.GetID())
	g, err = g.CreateIfNotExists(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{"amount": g.GetAmount()})
}

type GetPlayReq struct {
	Address string `json:"address"`
}

// GetPlay - Handler to generate play message for a user
func GetPlay(c *gin.Context) {
	var req GetPlayReq
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": playMessage})
}

type PlayReq struct {
	TokenID   uint   `json:"token_id"`
	Address   string `json:"address"`
	Signature string `json:"signature"`
}

func Play(c *gin.Context) {
	var req PlayReq
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sig, err := utils.SignatureToBytes(req.Signature)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// address string to common.Address
	addr := common.HexToAddress(req.Address)

	// verify signature
	valid, err := utils.VerifySignature(addr, playMessage, sig)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Signature verification failed."})
		return
	}

	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid signature."})
		return
	}

	// increase gem
	db := database.GetSQL()
	err = database.Transaction(db, func(tx *gorm.DB) error {
		g := models.NewGem(req.TokenID)
		g, err := g.GetByTokenIDAndLock(db)
		if err != nil {
			return err
		}
		if err := g.Update(db, map[string]interface{}{
			"amount": g.GetAmount() + 1,
		}); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}
