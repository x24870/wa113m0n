package database

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	gormpkg "wallemon/pkg/gorm"
	"wallemon/pkg/models"
)

// postgresDB is the concrete PostgresSQL handle to a SQL database.
type postgresDB struct{ *gorm.DB }

// initialize initializes the PostgreSQL database handle.
func (db *postgresDB) initialize(ctx context.Context) {
	var err error
	db.DB, err = gormpkg.NewGormPostgresConn(
		gormpkg.Config{
			// DSN:             config.GetDBArg(),
			// DSN:             "postgres://user:user@db:5432/wallemon?sslmode=disable", //TODO: use config
			DSN: "postgres://user:user@db:5432/postgres?sslmode=disable", //TODO: use config
			// DSN:             "host=db port=5432 user=user password=user dbname=postgres sslmode=disable binary_parameters=yes",
			MaxIdleConns:    2,
			MaxOpenConns:    2,
			ConnMaxLifetime: 10 * time.Minute,
			SingularTable:   true,
		},
	)
	if err != nil {
		panic(fmt.Errorf("failed to init db, err: %v", err))
	}

	// Perform database schema auto-migration.
	if err := models.AutoMigrate(db.DB); err != nil {
		panic(err)
	}

	// Load UUID extension if not loaded.
	stmt := fmt.Sprintf("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")
	if err = db.DB.Exec(stmt).Error; err != nil {
		panic(err)
	}
}

// finalize finalizes the PostgreSQL database handle.
func (db *postgresDB) finalize() {
	// Close the PostgreSQL database handle.
	sqlDB, err := db.DB.DB()
	if err != nil {
		// logging.Error(dbRootCtx, "Failed to get database handle: %v", err)
		fmt.Println("Failed to get database handle: %v", err)
	}
	if err := sqlDB.Close(); err != nil {
		// logging.Error(dbRootCtx, "Failed to close database handle: %v", err)
		fmt.Println("Failed to close database handle: %v", err)
	}
}

// db returns the PostgreSQL GORM database handle.
func (db *postgresDB) db() interface{} {
	return db.DB
}
