package database

import (
	"context"
	"fmt"
	"os"
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
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))

	db.DB, err = gormpkg.NewGormPostgresConn(
		gormpkg.Config{
			DSN:             dsn,
			MaxIdleConns:    2,
			MaxOpenConns:    2,
			ConnMaxLifetime: 10 * time.Minute,
			SingularTable:   true,
		},
	)
	if err != nil {
		panic(fmt.Errorf("failed to init db, dsn: \n%v \nerr: %v", dsn, err))
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
