package migrations

import (
	"gorm.io/gorm"
)

// MigrationFunc is the function signature for migration up/down functions
type MigrationFunc func(db *gorm.DB) error

// MigrationDefinition defines a single migration
type MigrationDefinition struct {
	Version string
	Name    string
	Up      MigrationFunc
	Down    MigrationFunc
}

// GetAllMigrations returns all migration definitions
// Add new migrations to this slice when creating new migration files
func GetAllMigrations() []MigrationDefinition {
	return []MigrationDefinition{
		Migration_20240101_001_CreateHospitalsTable(),
		Migration_20240101_002_CreateStaffsTable(),
		Migration_20240101_003_CreatePatientsTable(),
		Migration_20240101_004_SeedDefaultData(),
	}
}
