package models

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type GemInft interface {
	GetAmount() uint
	GetByTokenID(db *gorm.DB, tokenID uint) (GemInft, error)
	GetByTokenIDAndLock(db *gorm.DB, tokenID uint) (GemInft, error)
	CreateIfNotExists(db *gorm.DB, tokenID uint) (GemInft, error)
	Update(db *gorm.DB, values interface{}) error
}

// Gem is the exported static model interface.
var Gem gem

type gem struct {
	Base

	ID      uuid.UUID `gorm:"column:id;primary_key;type:uuid;default:uuid_generate_v4()"`
	TokenID uint      `gorm:"column:token_id;type:integer;not null"`
	Amount  uint      `gorm:"column:amount;type:integer;not null"`
}

func init() {
	registerModelForAutoMigration(&gem{})
}

func (t *gem) TableName() string {
	return "gem"
}

func (t *gem) Indexes() []CustomIndex {
	return []CustomIndex{
		{
			Name:      "created_at_idx",
			Unique:    false,
			Fields:    []string{"created_at"},
			Type:      "",
			Condition: "",
		},
		{
			Name:      "token_id_idx",
			Unique:    false,
			Fields:    []string{"token_id"},
			Type:      "",
			Condition: "",
		},
	}
}

func NewGem(tokenID uint) GemInft {
	g := gem{
		TokenID: tokenID,
		Amount:  0,
	}
	return &g
}

func (t *gem) GetAmount() uint {
	return t.Amount
}

func (t *gem) GetByTokenID(db *gorm.DB, tokenID uint) (GemInft, error) {
	if err := db.Where("token_id = ?", tokenID).First(t).Error; err != nil {
		return nil, err
	}
	return t, nil
}

func (t *gem) GetByTokenIDAndLock(db *gorm.DB, tokenID uint) (GemInft, error) {
	if err := db.Set("gorm:query_option", "FOR UPDATE").
		Where("token_id = ?", tokenID).
		First(t).Error; err != nil {
		return nil, err
	}
	return t, nil
}

func (t *gem) CreateIfNotExists(db *gorm.DB, tokenID uint) (GemInft, error) {
	if err := db.Where("token_id = ?", tokenID).FirstOrCreate(t).Error; err != nil {
		return nil, err
	}
	return t, nil
}

func (t *gem) Update(db *gorm.DB, values interface{}) error {
	if err := db.Model(t).Where("token_id = ?", t.TokenID).Updates(values).Error; err != nil {
		return err
	}
	return nil
}
