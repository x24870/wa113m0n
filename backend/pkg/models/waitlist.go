package models

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type WaitlistInft interface {
	Create(db *gorm.DB) (WaitlistInft, error)
}

// Block is the exported static model interface.
var Waitlist waitlist

type waitlist struct {
	Base

	ID         uuid.UUID `gorm:"column:id;primary_key;type:uuid;default:uuid_generate_v4()"`
	UserID     uuid.UUID `gorm:"column:user_id;type:uuid;not null"`
	ReferralID uuid.UUID `gorm:"column:referral_id;type:uuid"`
}

func init() {
	registerModelForAutoMigration(&waitlist{})
}

func (t *waitlist) TableName() string {
	return "waitlist"
}

func (t *waitlist) Indexes() []CustomIndex {
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

func NewWaitlist(userID uuid.UUID, referralID uuid.UUID) WaitlistInft {
	w := waitlist{
		UserID:     userID,
		ReferralID: referralID,
	}
	return &w
}

func (t *waitlist) GetUserID() uuid.UUID {
	return t.UserID
}

func (t *waitlist) GetReferralID() uuid.UUID {
	return t.ReferralID
}

func (t *waitlist) Create(db *gorm.DB) (WaitlistInft, error) {
	if err := db.Create(t).Error; err != nil {
		return nil, err
	}
	return t, nil
}
