package models

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type PoopInft interface {
	GetAmount() uint
	GetByTokenID(db *gorm.DB) (PoopInft, error)
	GetByTokenIDAndLock(db *gorm.DB) (PoopInft, error)
	CreateIfNotExists(db *gorm.DB) (PoopInft, error)
	Update(db *gorm.DB, values interface{}) error
}

// Poop is the exported static model interface.
var Poop poop

type poop struct {
	Base

	ID      uuid.UUID `gorm:"column:id;primary_key;type:uuid;default:uuid_generate_v4()"`
	TokenID uint      `gorm:"column:token_id;type:integer;not null"`
	Amount  uint      `gorm:"column:amount;type:integer;not null"`
}

func init() {
	registerModelForAutoMigration(&poop{})
}

func (t *poop) TableName() string {
	return "poop"
}

func (t *poop) Indexes() []CustomIndex {
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

func NewPoop(amount uint) PoopInft {
	p := poop{
		Amount: amount,
	}
	return &p
}

func (t *poop) GetAmount() uint {
	return t.Amount
}

func (t *poop) GetByTokenID(db *gorm.DB) (PoopInft, error) {
	if err := db.Where("token_id = ?", t.TokenID).First(t).Error; err != nil {
		return nil, err
	}
	return t, nil
}

func (t *poop) GetByTokenIDAndLock(db *gorm.DB) (PoopInft, error) {
	if err := db.Set("gorm:query_option", "FOR UPDATE").
		Where("token_id = ?", t.TokenID).
		First(t).Error; err != nil {
		return nil, err
	}
	return t, nil
}

func (t *poop) CreateIfNotExists(db *gorm.DB) (PoopInft, error) {
	if err := db.Where("token_id = ?", t.TokenID).FirstOrCreate(t).Error; err != nil {
		return nil, err
	}
	return t, nil
}

func (t *poop) Update(db *gorm.DB, values interface{}) error {
	if err := db.Model(t).Updates(values).Error; err != nil {
		return err
	}
	return nil
}
