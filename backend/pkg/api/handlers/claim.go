package api

import (
	"encoding/hex"
	"net/http"
	"os"

	"wallemon/pkg/database"
	"wallemon/pkg/models"
	utils "wallemon/pkg/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type ClaimReq struct {
	Address string `json:"address"`  // EVM address with 0x prefix
	RefCode string `json:"ref_code"` // Referral code
}

// Claim - Handler to claim a wallemon
func Claim(c *gin.Context) {
	var req ClaimReq

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// check 0x prefix, if not present, add it
	if req.Address[:2] != "0x" {
		req.Address = "0x" + req.Address
	}

	// check if address is valid EVM address
	if len(req.Address) != 42 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid address."})
		return
	}

	// check if refCode is valid, for now, just check if it's 'wallemon'
	if req.RefCode != "wallemon" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid referral code."})
		return
	}

	// get signer private key string from environment variable
	k := os.Getenv("SIGNER_KEY")
	if k[:2] == "0x" {
		k = k[2:]
	}
	if len(k) != 64 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error. signer key error"})
		return
	}

	// key string to private key then get signer
	pk, err := utils.KeyStringToPrivateKey(k)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error. pk error"})
		return
	}

	// sign
	to := common.HexToAddress(req.Address)
	msg := utils.EncodePacked(to.Bytes(), []byte(req.RefCode))
	msgHash := crypto.Keccak256(msg)
	signedMsg := utils.EncodePacked([]byte("\x19Ethereum Signed Message:\n32"), msgHash)
	signedMsgHash := crypto.Keccak256(signedMsg)
	signature, err := crypto.Sign(signedMsgHash, pk)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Sign error."})
		return
	}

	// modify last byte of signature to make it compatible with EVM
	signature[64] += 27

	// return hex string signature
	c.JSON(http.StatusOK, gin.H{"signature": hex.EncodeToString(signature)})
}

type JoinWaitlistReq struct {
	Email string `json:"email"`
}

// JoinWaitlist - Handler to join waitlist
func JoinWaitlist(c *gin.Context) {
	var req JoinWaitlistReq

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// check if email is valid
	if !utils.IsValidEmail(req.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email."})
		return
	}

	db := database.GetSQL()
	// check if email already exists
	_, err := models.User.GetByEmail(db, req.Email)
	if err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error."})
		return
	}

	err = database.Transaction(db, func(tx *gorm.DB) error {
		// create user
		user := models.NewUser(req.Email, "")
		user, err := user.Create(db)
		if err != nil {
			return err
		}
		// create waitlist
		wl := models.NewWaitlist(user.GetID(), uuid.Nil)
		if _, err := wl.Create(db); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error."})
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}
