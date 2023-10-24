package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"runtime/debug"

	"gorm.io/gorm"
)

// DB is the interface handle to a SQL database.
type DB interface {
	initialize(ctx context.Context)
	finalize()
	db() interface{}
}

// dbConfig is the config to connect to a SQL database.
type dbConfig struct {
	// The dialect of the SQL database.
	Dialect string

	// The username used to login to the database.
	Username string

	// The password used to login to the database.
	Password string

	// The address of the database service to connect to.
	Address string

	// The port of the database service to connect to.
	Port string

	// The name of the database to connect to.
	DBName string
}

// Global database interfaces.
var DBIntf DB

// Database root context.
var dbRootCtx context.Context

func init() {}

// Initialize initializes the database module and instance.
func Initialize(ctx context.Context) {
	// Save database root context.
	dbRootCtx = ctx

	// Create database according to dialect.
	Dialect := "postgres"
	switch Dialect {
	case "postgres", "cloudsqlpostgres":
		DBIntf = &postgresDB{}
	default:
		panic("invalid dialect")
	}

	// Initialize the database context.
	DBIntf.initialize(dbRootCtx)
}

// Finalize finalizes the database module and closes the database handles.
func Finalize() {
	// Make sure database instance has been initialized.
	if DBIntf == nil {
		panic("database has not been initialized")
	}

	// Finalize database instance.
	DBIntf.finalize()
}

// GetDB returns the database instance.
func GetDB() interface{} {
	return DBIntf.db()
}

// GetSQL returns the SQL database instance.
func GetSQL() *gorm.DB {
	return GetDB().(*gorm.DB)
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
