package models

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type Waitlist struct {
	Base

	ID         uuid.UUID `gorm:"column:id;primary_key;type:uuid;default:uuid_generate_v4()"`
	UserID     uuid.UUID `gorm:"column:user_id;type:uuid;not null"`
	ReferralID uuid.UUID `gorm:"column:referral_id;type:uuid"`
}

func init() {
	registerModelForAutoMigration(&Waitlist{})
}

func (t *Waitlist) TableName() string {
	return "waitlist"
}

func (t *Waitlist) Indexes() []CustomIndex {
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

func NewWaitlist(userID uuid.UUID, referralID uuid.UUID) *Waitlist {
	return &Waitlist{
		UserID:     userID,
		ReferralID: referralID,
	}
}

func (t *Waitlist) GetUserID() uuid.UUID {
	return t.UserID
}

func (t *Waitlist) GetReferralID() uuid.UUID {
	return t.ReferralID
}

func (t *Waitlist) CreateIfNotExists(db *gorm.DB) error {
	return db.FirstOrCreate(t, Waitlist{UserID: t.UserID}).Error
}
