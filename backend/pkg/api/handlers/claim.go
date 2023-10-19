package api

import (
	"encoding/hex"
	"net/http"
	"os"

	utils "wallemon/pkg/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
)

type ClaimReq struct {
	Address string `json:"address"`  // EVM address with 0x prefix
	RefCode string `json:"ref_code"` // Referral code
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
	to := common.HexToAddress(gq.Address)
	msg := utils.EncodePacked(to.Bytes(), []byte(gq.RefCode))
	msgHash := crypto.Keccak256(msg)
	signedMsg := utils.EncodePacked([]byte("\x19Ethereum Signed Message:\n32"), msgHash)
	signedMsgHash := crypto.Keccak256(signedMsg)
	signature, err := crypto.Sign(signedMsgHash, pk)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Sign error."})
		return
	}

	// return hex string signature
	c.JSON(http.StatusOK, gin.H{"signature": hex.EncodeToString(signature)})
}
