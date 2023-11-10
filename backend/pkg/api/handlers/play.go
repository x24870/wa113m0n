package api

import (
	"fmt"
	"net/http"
	"strconv"

	"wallemon/pkg/database"
	"wallemon/pkg/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const playMessage = "Let's play!"
const cleanMessage = "smells good!"

type GetGemReq struct {
	TokenID uint `json:"token_id"`
}

type GetGemResp struct {
	Amount uint `json:"amount"`
}

// GetGem - Handler to get the gem of a user
func GetGem(c *gin.Context) {
	id := c.Request.Header.Get("token_id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "token_id is required"})
		return
	}

	val, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "token_id is invalid"})
		return
	}

	req := GetGemReq{
		TokenID: uint(val),
	}

	// TODO: maybe check if this address owns this token

	db := database.GetSQL()
	t := models.NewToken(req.TokenID)
	t, err = t.CreateIfNotExists(db, req.TokenID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	g := models.NewGem(t.GetID())
	g, err = g.CreateIfNotExists(db, t.GetID())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, GetGemResp{
		Amount: g.GetAmount(),
	})
}

type GetPlayReq struct {
	Address string `json:"address"`
}

type GetPlayResp struct {
	Message string `json:"message"`
}

// GetPlay - Handler to generate play message for a user
func GetPlay(c *gin.Context) {
	addr := c.Request.Header.Get("address")
	if addr == "" {
		// address is unused currently
		c.JSON(http.StatusBadRequest, gin.H{"error": "address is required"})
		return
	}

	c.JSON(http.StatusOK, GetPlayResp{
		Message: playMessage,
	})
}

type PlayReq struct {
	TokenID   uint   `json:"token_id"`
	Address   string `json:"address"`
	Signature string `json:"signature"`
}

type PlayResp struct {
	Message string `json:"message"`
}

func Play(c *gin.Context) {
	var req PlayReq
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// sig, err := utils.SignatureToBytes(req.Signature)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }

	// // address string to common.Address
	// addr := common.HexToAddress(req.Address)

	// // verify signature
	// valid, err := utils.VerifySignature(addr, playMessage, sig)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Signature verification failed."})
	// 	return
	// }

	// if !valid {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "not signed by the given address"})
	// 	return
	// }

	db := database.GetSQL()
	// check if token is healthy
	t := models.NewToken(req.TokenID)
	// t, err := t.GetByTokenID(db)
	t, err := t.CreateIfNotExists(db, req.TokenID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error."})
		return
	}
	if t.GetState() != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Your partner is not healthy."})
		return
	}

	// check if play limit reached
	reached, err := models.OpLog.PlayLimitReached(db, req.TokenID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error."})
		return
	}
	if reached {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Play limit reached."})
		return
	}

	// increase gem
	err = database.Transaction(db, func(tx *gorm.DB) error {
		g := models.NewGem(req.TokenID)
		g, err := g.CreateIfNotExists(db, req.TokenID)
		if err != nil {
			return err
		}

		g, err = g.GetByTokenIDAndLock(db, req.TokenID)
		if err != nil {
			return err
		}
		if err := g.Update(db, map[string]interface{}{
			"amount": g.GetAmount() + 1,
		}); err != nil {
			return err
		}

		l := models.NewOpLog(req.TokenID, string(models.OperationTypePlay))
		if _, err := l.Create(db); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error."})
		return
	}

	c.JSON(http.StatusOK, PlayResp{
		Message: "｡◕‿◕｡",
	})
}

type GetPoopReq struct {
	TokenID uint `json:"token_id"`
}

type GetPoopResp struct {
	Amount uint `json:"amount"`
}

func GetPoop(c *gin.Context) {
	id := c.Request.Header.Get("token_id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "token_id is required"})
		return
	}
	fmt.Println("!!!GetPoop id: ", id)

	val, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "token_id is invalid"})
		return
	}

	req := GetPoopReq{
		TokenID: uint(val),
	}
	fmt.Println("!!!req.token_id: ", req.TokenID)

	// TODO: maybe check if this address owns this token

	// create poop if not exists
	db := database.GetSQL()
	p := models.NewPoop(req.TokenID)
	p, err = p.CreateIfNotExists(db, req.TokenID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, GetGemResp{
		Amount: p.GetAmount(),
	})
}

type GetCleanReq struct {
	Address string `json:"address"`
}

type GetCleanResp struct {
	Message string `json:"message"`
}

func GetClean(c *gin.Context) {
	addr := c.Request.Header.Get("address")
	if addr == "" {
		// address is unused currently
		c.JSON(http.StatusBadRequest, gin.H{"error": "address is required"})
		return
	}

	c.JSON(http.StatusOK, GetCleanResp{
		Message: cleanMessage,
	})
}

type CleanReq struct {
	TokenID   uint   `json:"token_id"`
	Address   string `json:"address"`
	Signature string `json:"signature"`
}

type CleanResp struct {
	Message string `json:"message"`
}

func Clean(c *gin.Context) {
	var req CleanReq
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("!!!Clean id: ", req.TokenID)

	// sig, err := utils.SignatureToBytes(req.Signature)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }

	// // address string to common.Address
	// addr := common.HexToAddress(req.Address)

	// // verify signature
	// valid, err := utils.VerifySignature(addr, cleanMessage, sig)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Signature verification failed."})
	// 	return
	// }

	// if !valid {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "not signed by the given address"})
	// 	return
	// }

	// clean poops
	db := database.GetSQL()
	p := models.NewPoop(req.TokenID)
	// p, err := p.GetByTokenID(db)
	p, err := p.CreateIfNotExists(db, req.TokenID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = database.Transaction(db, func(tx *gorm.DB) error {
		p, err = p.GetByTokenIDAndLock(db, req.TokenID)
		if err != nil {
			return err
		}
		if err := p.Update(db, map[string]interface{}{
			"amount": 0,
		}); err != nil {
			return err
		}

		l := models.NewOpLog(req.TokenID, string(models.OperationTypeClean))
		if _, err := l.Create(db); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error."})
		return
	}

	c.JSON(http.StatusOK, CleanResp{
		Message: "(〃'▽'〃)",
	})
}
