package models

import (
	"gorm.io/gorm"
)

type TokenInft interface {
	GetID() uint
	Create(db *gorm.DB) (TokenInft, error)
	CreateIfNotExists(db *gorm.DB) (TokenInft, error)
	GetByTokenID(db *gorm.DB) (TokenInft, error)
	GetByTokenIDAndLock(db *gorm.DB) (TokenInft, error)
}

// Token is the exported static model interface.
var Token token

type token struct {
	Base

	// TokenID is ERC721 token ID.
	// Which is a unique uint value.
	TokenID uint `gorm:"column:token_id;type:integer;not null"`
}

func init() {
	registerModelForAutoMigration(&token{})
}

func (t *token) TableName() string {
	return "token"
}

func (t *token) Indexes() []CustomIndex {
	return []CustomIndex{
		{
			Name:      "created_at_idx",
			Unique:    false,
			Fields:    []string{"created_at"},
			Type:      "",
			Condition: "",
		},
	}
}

func NewToken(tokenID uint) TokenInft {
	t := token{
		TokenID: tokenID,
	}
	return &t
}

func (t *token) GetID() uint {
	return t.TokenID
}

func (t *token) Create(db *gorm.DB) (TokenInft, error) {
	if err := db.Create(t).Error; err != nil {
		return nil, err
	}
	return t, nil
}

func (t *token) CreateIfNotExists(db *gorm.DB) (TokenInft, error) {
	if err := db.Where("token_id = ?", t.TokenID).FirstOrCreate(t).Error; err != nil {
		return nil, err
	}
	return t, nil
}

func (t *token) GetByTokenID(db *gorm.DB) (TokenInft, error) {
	if err := db.Where("token_id = ?", t.TokenID).
		First(t).Error; err != nil {
		return nil, err
	}
	return t, nil
}

func (t *token) GetByTokenIDAndLock(db *gorm.DB) (TokenInft, error) {
	if err := db.Set("gorm:query_option", "FOR UPDATE").
		Where("token_id = ?", t.TokenID).
		First(t).Error; err != nil {
		return nil, err
	}
	return t, nil
}
