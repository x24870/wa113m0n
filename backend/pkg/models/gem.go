package models

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type GemInft interface {
	GetAmount() uint
	GetByTokenID(db *gorm.DB) (GemInft, error)
	GetByTokenIDAndLock(db *gorm.DB) (GemInft, error)
	CreateIfNotExists(db *gorm.DB) (GemInft, error)
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
	}
}

func NewGem(amount uint) GemInft {
	g := gem{
		Amount: amount,
	}
	return &g
}

func (t *gem) GetAmount() uint {
	return t.Amount
}

func (t *gem) GetByTokenID(db *gorm.DB) (GemInft, error) {
	if err := db.Where("token_id = ?", t.TokenID).First(t).Error; err != nil {
		return nil, err
	}
	return t, nil
}

func (t *gem) GetByTokenIDAndLock(db *gorm.DB) (GemInft, error) {
	if err := db.Set("gorm:query_option", "FOR UPDATE").
		Where("token_id = ?", t.TokenID).
		First(t).Error; err != nil {
		return nil, err
	}
	return t, nil
}

func (t *gem) CreateIfNotExists(db *gorm.DB) (GemInft, error) {
	if err := db.Where("token_id = ?", t.TokenID).FirstOrCreate(t).Error; err != nil {
		return nil, err
	}
	return t, nil
}

func (t *gem) Update(db *gorm.DB, values interface{}) error {
	if err := db.Model(t).Updates(values).Error; err != nil {
		return err
	}
	return nil
}
