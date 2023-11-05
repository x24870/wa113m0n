package models

import (
	"fmt"

	"gorm.io/gorm"
)

type TokenInft interface {
	GetID() uint
	Create(db *gorm.DB) (TokenInft, error)
	CreateIfNotExists(db *gorm.DB) (TokenInft, error)
	GetByTokenID(db *gorm.DB) (TokenInft, error)
	GetByTokenIDAndLock(db *gorm.DB) (TokenInft, error)
	GetState() uint
	Update(db *gorm.DB, values interface{}) error
}

// ErrInvalidTokenID is returned when token ID is invalid.
var ErrInvalidTokenID = fmt.Errorf("invalid token ID")

// Token is the exported static model interface.
var Token token

type token struct {
	Base

	// TokenID is ERC721 token ID.
	// Which is a unique uint value.
	TokenID uint `gorm:"column:token_id;type:integer;not null"`
	// State is the state of the token.
	// 0: healthy
	// 1: sick
	// 2: dead
	State uint `gorm:"column:state;type:integer;not null;default:0"`
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
	if t.TokenID >= 1000 {
		return nil, ErrInvalidTokenID
	}
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

func (t *token) GetState() uint {
	return t.State
}

func (t *token) Update(db *gorm.DB, values interface{}) error {
	if err := db.Model(t).Where("token_id = ?", t.TokenID).Updates(values).Error; err != nil {
		return err
	}
	return nil
}
