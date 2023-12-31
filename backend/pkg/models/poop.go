package models

import (
	"errors"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type PoopInft interface {
	GetAmount() uint
	GetTokenID() uint
	GetByTokenID(db *gorm.DB, tokenID uint) (PoopInft, error)
	GetByTokenIDAndLock(db *gorm.DB, tokenID uint) (PoopInft, error)
	CreateIfNotExists(db *gorm.DB, tokenID uint) (PoopInft, error)
	Update(db *gorm.DB, values interface{}) error
	ListShouldSick(db *gorm.DB) ([]uint, error)
	List(db *gorm.DB) ([]PoopInft, error)
}

const PoopDuration = 300 // seconds
const PoopMaxAmount = 6

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
		{
			Name:      "token_id_idx",
			Unique:    false,
			Fields:    []string{"token_id"},
			Type:      "",
			Condition: "",
		},
	}
}

func NewPoop(tokenID uint) PoopInft {
	p := poop{
		TokenID: tokenID,
		Amount:  0,
	}
	return &p
}

func (t *poop) GetAmount() uint {
	return t.Amount
}

func (t *poop) GetTokenID() uint {
	return t.TokenID
}

func (t *poop) GetByTokenID(db *gorm.DB, tokenID uint) (PoopInft, error) {
	if err := db.Where("token_id = ?", t.TokenID).First(t).Error; err != nil {
		return nil, err
	}
	return t, nil
}

func (t *poop) GetByTokenIDAndLock(db *gorm.DB, tokenID uint) (PoopInft, error) {
	if err := db.Set("gorm:query_option", "FOR UPDATE").
		Where("token_id = ?", tokenID).
		First(t).Error; err != nil {
		return nil, err
	}
	return t, nil
}

func (t *poop) CreateIfNotExists(db *gorm.DB, tokenID uint) (PoopInft, error) {
	if err := db.Where("token_id = ?", tokenID).FirstOrCreate(t).Error; err != nil {
		return nil, err
	}
	return t, nil
}

func (t *poop) Update(db *gorm.DB, values interface{}) error {
	if err := db.Model(t).Where("token_id = ?", t.TokenID).Updates(values).Error; err != nil {
		return err
	}
	return nil
}

// ListShouldSick returns a list of tokenIDs that should be sick.
// If the amount of poop is greater than 6, it returns the tokenID.
func (t *poop) ListShouldSick(db *gorm.DB) ([]uint, error) {
	var poops []poop
	if err := db.Where("amount >= 6").Find(&poops).Error; err != nil {
		return nil, err
	}

	var ret []uint
	for _, p := range poops {
		t := NewToken(p.TokenID)
		t, err := t.GetByTokenID(db, p.TokenID)
		if err != nil {
			continue // TODO
		}
		if t.GetState() != 0 {
			continue
		}

		ret = append(ret, p.TokenID)
	}

	return ret, nil
}

func (t *poop) List(db *gorm.DB) ([]PoopInft, error) {
	poops := []*poop{}
	err := db.Find(&poops).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	poopInfts := []PoopInft{}
	for _, p := range poops {
		poopInfts = append(poopInfts, p)
	}

	return poopInfts, nil
}
