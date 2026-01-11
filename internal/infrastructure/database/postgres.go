package database

import (
	"fmt"
	"log"

	"github.com/wichai2002/his_v1/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewPostgresDB(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port, cfg.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Connected to PostgreSQL database successfully")
	return db, nil
}

// RunMigrations runs all pending migrations using the migration system
func RunMigrations(db *gorm.DB) error {
	migrator := NewMigrator(db)
	return migrator.MigrateUp()
}

// RollbackMigration rolls back the last migration
func RollbackMigration(db *gorm.DB) error {
	migrator := NewMigrator(db)
	return migrator.MigrateDown()
}

// RollbackAllMigrations rolls back all migrations
func RollbackAllMigrations(db *gorm.DB) error {
	migrator := NewMigrator(db)
	return migrator.MigrateDownAll()
}

// MigrationStatus prints the current migration status
func MigrationStatus(db *gorm.DB) {
	migrator := NewMigrator(db)
	migrator.Status()
}
