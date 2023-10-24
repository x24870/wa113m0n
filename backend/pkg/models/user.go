package models

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type UserInft interface {
	GetEmail()
	GetAddress()
}

type User struct {
	Base

	ID    uuid.UUID `gorm:"column:id;primary_key;type:uuid;default:uuid_generate_v4()"`
	Email string    `gorm:"column:email;type:varchar(256);not null"`
	// EVM address with 0x prefix
	Address string `gorm:"column:address;type:varchar(256);not null"`
}

func init() {
	registerModelForAutoMigration(&User{})
}

func (u *User) TableName() string {
	return "user"
}

func (u *User) Indexes() []CustomIndex {
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

func NewUser(email string, address string) *User {
	return &User{
		Email:   email,
		Address: address,
	}
}

func (u *User) GetEmail() string {
	return u.Email
}

func (u *User) GetAddress() string {
	return u.Address
}

func (u *User) CreateIfNotExists(db *gorm.DB) error {
	return db.Where("email = ?", u.Email).FirstOrCreate(u).Error
}
