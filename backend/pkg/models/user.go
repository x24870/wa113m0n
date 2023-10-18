package models

import (
	uuid "github.com/satori/go.uuid"
)

type User struct {
	Base

	ID    uuid.UUID `gorm:"column:id;primary_key;type:uuid;default:uuid_generate_v4()"`
	Email string    `gorm:"column:email;type:varchar(256);not null"`
	// EVM address with 0x prefix
	Address string `gorm:"column:address;type:varchar(256);not null"`
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
			Condition: " WHERE disabled = false",
		},
	}
}
