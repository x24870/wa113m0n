package gorm

import (
	"database/sql"
	"errors"
	"fmt"
	"runtime/debug"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func NewGormPostgresConn(config Config) (*gorm.DB, error) {
	config = padDefault(config)
	db, err := gorm.Open(
		postgres.New(
			postgres.Config{
				// We use pgBouncer with transaction pooling mode,
				// which not supports extended query protocol.
				// ref: https://www.pgbouncer.org/faq.html
				// ref: https://gorm.io/docs/connecting_to_the_database.html
				PreferSimpleProtocol: true,
				DSN:                  config.DSN,
			}),
		&gorm.Config{
			Logger: logger.Default.LogMode(logger.Warn),
			NamingStrategy: schema.NamingStrategy{
				SingularTable: config.SingularTable,
			},
			// In order to use the dialect translated errors(like ErrDuplicatedKey),
			// enable the TranslateError flag when opening a db connection.
			TranslateError: true,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("gorm.Open error: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("db.DB error: %w", err)
	}

	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(config.ConnMaxLifetime)

	return db, nil
}

// Transaction wraps the database transaction and to proper error handling.
func Transaction(db *gorm.DB, body func(*gorm.DB) error) (err error) {
	tx := db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("Transaction: Cannot open transaction %v", tx.Error)
	}

	// Handle runtime.Goexit. err won't be set when Goexit is called in body.
	errDefault := errors.New("init")
	err = errDefault

	// Error checking and panic safenet.
	defer func() {
		if recovered := recover(); recovered != nil {
			err = fmt.Errorf("Transaction: rollback due to panic: %v, %s",
				recovered, string(debug.Stack()))
		}

		if err != nil {
			rollbackErr := tx.Rollback().Error
			if rollbackErr == nil || errors.Is(err, sql.ErrTxDone) {
				return
			}
			err = fmt.Errorf("Transaction: rollback due to error: %v, %w", err, rollbackErr)
			return
		}
	}()

	// Execute main body.
	if err = body(tx); err != nil {
		return err
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return nil
}
