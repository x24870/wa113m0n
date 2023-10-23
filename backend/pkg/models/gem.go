package models

import (
	uuid "github.com/satori/go.uuid"
)

type Gem struct {
	Base

	ID     uuid.UUID `gorm:"column:id;primary_key;type:uuid;default:uuid_generate_v4()"`
	Amount uint      `gorm:"column:amount;type:integer;not null"`
}

func init() {
	registerModelForAutoMigration(&Gem{})
}

func (t *Gem) TableName() string {
	return "gem"
}

func (t *Gem) Indexes() []CustomIndex {
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
