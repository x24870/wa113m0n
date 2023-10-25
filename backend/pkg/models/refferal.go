package models

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type ReferralInft interface {
	GetID() uuid.UUID
	GetCode() string
	GetLimit() uint
	Create(db *gorm.DB) (ReferralInft, error)
	GetByCode(db *gorm.DB, code string) (ReferralInft, error)
}

var Referral referral

type referral struct {
	Base

	ID    uuid.UUID `gorm:"column:id;primary_key;type:uuid;default:uuid_generate_v4()"`
	Code  string    `gorm:"column:code;type:varchar(255);not null"`
	Limit uint      `gorm:"column:limit;type:integer;not null"`
}

func init() {
	registerModelForAutoMigration(&referral{})
}

func (t *referral) TableName() string {
	return "referral"
}

func (t *referral) Indexes() []CustomIndex {
	return []CustomIndex{
		{
			Name:      "created_at_idx",
			Unique:    false,
			Fields:    []string{"created_at"},
			Type:      "",
			Condition: "",
		},
		{
			Name:      "code_idx",
			Unique:    true,
			Fields:    []string{"code"},
			Type:      "",
			Condition: "",
		},
	}
}

func NewReferral(code string, limit uint) ReferralInft {
	r := referral{
		Code:  code,
		Limit: limit,
	}
	return &r
}

func (r *referral) GetID() uuid.UUID {
	return r.ID
}

func (r *referral) GetCode() string {
	return r.Code
}

func (r *referral) GetLimit() uint {
	return r.Limit
}

func (r *referral) Create(db *gorm.DB) (ReferralInft, error) {
	if err := db.Create(r).Error; err != nil {
		return nil, err
	}
	return r, nil
}

func (r *referral) GetByCode(db *gorm.DB, code string) (ReferralInft, error) {
	if err := db.Where("code = ?", code).First(r).Error; err != nil {
		return nil, err
	}
	return r, nil
}
