package models

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type UserInft interface {
	GetID() uuid.UUID
	GetEmail() string
	GetAddress() string
	Create(db *gorm.DB) (UserInft, error)
}

// User is the exported static model interface.
var User user

type user struct {
	Base

	ID    uuid.UUID `gorm:"column:id;primary_key;type:uuid;default:uuid_generate_v4()"`
	Email string    `gorm:"column:email;type:varchar(256);not null"`
	// EVM address with 0x prefix
	Address string `gorm:"column:address;type:varchar(256);not null"`
}

func init() {
	registerModelForAutoMigration(&user{})
}

func (u *user) TableName() string {
	return "user"
}

func (u *user) Indexes() []CustomIndex {
	return []CustomIndex{
		{
			Name:      "created_at_idx",
			Unique:    false,
			Fields:    []string{"created_at"},
			Type:      "",
			Condition: "",
		},
		{
			Name:      "email_idx",
			Unique:    true,
			Fields:    []string{"email"},
			Type:      "",
			Condition: "",
		},
	}
}

func NewUser(email string, address string) UserInft {
	u := user{
		Email:   email,
		Address: address,
	}

	return &u
}

func (u *user) GetID() uuid.UUID {
	return u.ID
}

func (u *user) GetEmail() string {
	return u.Email
}

func (u *user) GetAddress() string {
	return u.Address
}

func (u *user) GetByEmail(db *gorm.DB, email string) (UserInft, error) {
	if err := db.Where("email = ?", email).First(u).Error; err != nil {
		return nil, err
	}
	return u, nil
}

// Create
func (u *user) Create(db *gorm.DB) (UserInft, error) {
	if err := db.Create(u).Error; err != nil {
		return nil, err
	}
	return u, nil
}

func (u *user) CreateIfNotExists(db *gorm.DB) (UserInft, error) {
	if err := db.Where("email = ?", u.Email).FirstOrCreate(u).Error; err != nil {
		return nil, err
	}
	return u, nil
}
