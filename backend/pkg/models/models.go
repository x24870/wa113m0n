package models

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// Model ...
type model interface {
	Indexes() []CustomIndex
	TableName() string
}

// models ...
var models = []model{}

// registerModelForMigration...
func registerModelForAutoMigration(model model) {
	models = append(models, model)
}

// AutoMigrate ...
func AutoMigrate(db *gorm.DB) error {
	// Turn on logging for migration.
	db = db.Debug()

	// Begin a new transaction.
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		return tx.Error
	}

	// Perform migration on models.
	for _, model := range models {
		err := tx.AutoMigrate(model)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// Create indexes, and foreign keys for each model.
	for _, m := range models {
		for _, idx := range m.Indexes() {
			unique := ""
			extension := ""
			if idx.Unique {
				unique = "UNIQUE"
			}
			if len(idx.Type) != 0 {
				extension = "USING " + idx.Type
			}
			columns := strings.Join(idx.Fields, ",")
			idxStat := fmt.Sprintf(
				`CREATE %s INDEX IF NOT EXISTS %s_%s ON "%s" %s(%s) %s`,
				unique, m.TableName(), idx.Name, m.TableName(), extension, columns, idx.Condition)
			err := tx.Exec(idxStat).Error
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	// Commit the transaction.
	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
