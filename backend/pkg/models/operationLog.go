package models

import (
	"database/sql/driver"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// operation log table is for logging the operations of the game.
// like play, clean etc.

type OperationType string

const (
	OperationTypePlay  OperationType = "play"
	OperationTypeClean OperationType = "clean"
)

var operationType = map[OperationType]struct{}{
	OperationTypePlay:  {},
	OperationTypeClean: {},
}

// IsValid return true if totpauthAction is valid.
func (t *OperationType) IsValid() bool {
	if _, ok := operationType[*t]; ok {
		return true
	}
	return false
}

// Scan implements database/sql.Scanner.
func (x *OperationType) Scan(value interface{}) error {
	strValue, ok := value.(string)
	if !ok {
		return fmt.Errorf("OperationType: value.(string) failed. value(%+v)", value)
	}

	newValue := OperationType(strValue)
	if newValue.IsValid() {
		*x = newValue
		return nil
	}

	return fmt.Errorf("OperationType: value is invalid. strValue(%s)", strValue)
}

// Value implements database/sql/driver.Valuer.
func (x OperationType) Value() (driver.Value, error) {
	if len(x) != 0 {
		return string(x), nil
	}
	err := fmt.Errorf("OperationType: value is empty string")
	return nil, err
}

type OpLogInft interface {
	Create(db *gorm.DB) (OpLogInft, error)
	PlayLimitReached(db *gorm.DB, tokenID uint) (bool, error)
}

// OpLog is the exported static model interface.
var OpLog opLog

type opLog struct {
	Base

	ID      uuid.UUID `gorm:"column:id;primary_key;type:uuid;default:uuid_generate_v4()"`
	TokenID uint      `gorm:"column:token_id;type:integer;not null"`
	OpType  string    `gorm:"column:op_type;type:varchar(255);not null"`
}

func init() {
	registerModelForAutoMigration(&opLog{})
}

func (t *opLog) TableName() string {
	return "operation_log"
}

func (t *opLog) Indexes() []CustomIndex {
	return []CustomIndex{
		{
			Name:      "created_at_idx",
			Unique:    false,
			Fields:    []string{"created_at"},
			Type:      "",
			Condition: "",
		},
		{
			Name:      "token_id_op_type_created_at_idx",
			Unique:    false,
			Fields:    []string{"token_id", "op_type", "created_at"},
			Type:      "",
			Condition: "",
		},
	}
}

func NewOpLog(tokenID uint, opType string) OpLogInft {
	t := opLog{
		TokenID: tokenID,
		OpType:  opType,
	}
	return &t
}

func (t *opLog) Create(db *gorm.DB) (OpLogInft, error) {
	if err := db.Create(t).Error; err != nil {
		return nil, err
	}
	return t, nil
}

// PlayLimitReached returns true if the play limit is reached for the token.
// The play limit is 3 times per day.
func (t *opLog) PlayLimitReached(db *gorm.DB, tokenID uint) (bool, error) {
	var count int64
	if err := db.Model(&opLog{}).
		Where("token_id = ?", tokenID).
		Where("op_type = ?", OperationTypePlay).
		Where("created_at > now() - interval '1 day'").
		Count(&count).Error; err != nil {
		return false, err
	}
	if count >= 3 {
		return true, nil
	}
	return false, nil
}
