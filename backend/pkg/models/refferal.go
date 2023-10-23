package models

import uuid "github.com/satori/go.uuid"

type Referral struct {
	Base

	ID     uuid.UUID `gorm:"column:id;primary_key;type:uuid;default:uuid_generate_v4()"`
	UserID uuid.UUID `gorm:"column:user_id;type:uuid;not null"`
	Code   string    `gorm:"column:code;type:varchar(255);not null"`
	Limit  uint      `gorm:"column:limit;type:integer;not null"`
	Count  uint      `gorm:"column:count;type:integer;not null"`
}

func init() {
	registerModelForAutoMigration(&Referral{})
}

func (t *Referral) TableName() string {
	return "referral"
}

func (t *Referral) Indexes() []CustomIndex {
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

type ReferralHistory struct {
	Base

	ID         uuid.UUID `gorm:"column:id;primary_key;type:uuid;default:uuid_generate_v4()"`
	ReferralID uuid.UUID `gorm:"column:referral_id;type:uuid;not null"`
	UserID     uuid.UUID `gorm:"column:user_id;type:uuid;not null"`
}
