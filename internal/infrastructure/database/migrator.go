package database

import (
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/wichai2002/his_v1/internal/infrastructure/database/migrations"

	"gorm.io/gorm"
)

// MigrationRecord represents a database migration record
type MigrationRecord struct {
	ID        uint      `gorm:"primaryKey"`
	Version   string    `gorm:"uniqueIndex;not null"`
	Name      string    `gorm:"not null"`
	AppliedAt time.Time `gorm:"autoCreateTime"`
}

// TableName specifies the table name for migration records
func (MigrationRecord) TableName() string {
	return "schema_migrations"
}

// Migrator handles database migrations
type Migrator struct {
	db         *gorm.DB
	migrations []migrations.MigrationDefinition
}

// NewMigrator creates a new migrator instance
func NewMigrator(db *gorm.DB) *Migrator {
	return &Migrator{
		db:         db,
		migrations: migrations.GetAllMigrations(),
	}
}

// Initialize creates the migrations table if it doesn't exist
func (m *Migrator) Initialize() error {
	return m.db.AutoMigrate(&MigrationRecord{})
}

// GetAppliedMigrations returns all applied migrations
func (m *Migrator) GetAppliedMigrations() ([]MigrationRecord, error) {
	var records []MigrationRecord
	if err := m.db.Order("version ASC").Find(&records).Error; err != nil {
		return nil, err
	}
	return records, nil
}

// IsMigrationApplied checks if a migration version has been applied
func (m *Migrator) IsMigrationApplied(version string) bool {
	var count int64
	m.db.Model(&MigrationRecord{}).Where("version = ?", version).Count(&count)
	return count > 0
}

// MigrateUp runs all pending migrations
func (m *Migrator) MigrateUp() error {
	if err := m.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize migrations table: %w", err)
	}

	// Sort migrations by version
	sort.Slice(m.migrations, func(i, j int) bool {
		return m.migrations[i].Version < m.migrations[j].Version
	})

	for _, migration := range m.migrations {
		if m.IsMigrationApplied(migration.Version) {
			log.Printf("Migration %s (%s) already applied, skipping...", migration.Version, migration.Name)
			continue
		}

		log.Printf("Applying migration %s: %s", migration.Version, migration.Name)

		// Run migration in transaction
		err := m.db.Transaction(func(tx *gorm.DB) error {
			if err := migration.Up(tx); err != nil {
				return err
			}

			// Record migration
			record := MigrationRecord{
				Version: migration.Version,
				Name:    migration.Name,
			}
			return tx.Create(&record).Error
		})

		if err != nil {
			return fmt.Errorf("failed to apply migration %s: %w", migration.Version, err)
		}

		log.Printf("Migration %s applied successfully", migration.Version)
	}

	return nil
}

// MigrateDown rolls back the last migration
func (m *Migrator) MigrateDown() error {
	if err := m.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize migrations table: %w", err)
	}

	// Get the last applied migration
	var lastMigration MigrationRecord
	if err := m.db.Order("version DESC").First(&lastMigration).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Println("No migrations to rollback")
			return nil
		}
		return err
	}

	// Find the migration definition
	var migrationDef *migrations.MigrationDefinition
	for _, mig := range m.migrations {
		if mig.Version == lastMigration.Version {
			migrationDef = &mig
			break
		}
	}

	if migrationDef == nil {
		return fmt.Errorf("migration definition not found for version %s", lastMigration.Version)
	}

	log.Printf("Rolling back migration %s: %s", migrationDef.Version, migrationDef.Name)

	// Run rollback in transaction
	err := m.db.Transaction(func(tx *gorm.DB) error {
		if err := migrationDef.Down(tx); err != nil {
			return err
		}

		// Remove migration record
		return tx.Where("version = ?", migrationDef.Version).Delete(&MigrationRecord{}).Error
	})

	if err != nil {
		return fmt.Errorf("failed to rollback migration %s: %w", migrationDef.Version, err)
	}

	log.Printf("Migration %s rolled back successfully", migrationDef.Version)
	return nil
}

// MigrateDownAll rolls back all migrations
func (m *Migrator) MigrateDownAll() error {
	applied, err := m.GetAppliedMigrations()
	if err != nil {
		return err
	}

	for i := len(applied) - 1; i >= 0; i-- {
		if err := m.MigrateDown(); err != nil {
			return err
		}
	}

	return nil
}

// Status prints the migration status
func (m *Migrator) Status() {
	applied, _ := m.GetAppliedMigrations()
	appliedMap := make(map[string]MigrationRecord)
	for _, mig := range applied {
		appliedMap[mig.Version] = mig
	}

	fmt.Println("\n=== Migration Status ===")
	fmt.Printf("%-20s %-40s %-10s %-25s\n", "VERSION", "NAME", "STATUS", "APPLIED AT")
	fmt.Println("--------------------------------------------------------------------------------")

	for _, mig := range m.migrations {
		status := "Pending"
		appliedAt := ""
		if record, ok := appliedMap[mig.Version]; ok {
			status = "Applied"
			appliedAt = record.AppliedAt.Format("2006-01-02 15:04:05")
		}
		fmt.Printf("%-20s %-40s %-10s %-25s\n", mig.Version, mig.Name, status, appliedAt)
	}
	fmt.Println()
}
