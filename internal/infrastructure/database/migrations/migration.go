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

// GetAllMigrations returns migration definitions for the PUBLIC schema only
// NOTE: Staff and Patient tables are tenant-specific and created via tenant_service.go
// when setting up a new tenant schema - they should NOT be in public schema
func GetAllMigrations() []MigrationDefinition {
	return []MigrationDefinition{
		Migration_20240101_005_CreateTenantsTable(),
		Migration_20240101_007_AddHospitalFieldsToTenants(),
	}
}

// GetAllMigrationsIncludingLegacy returns ALL migrations including legacy ones
// This is needed for rollback operations to properly find migration definitions
func GetAllMigrationsIncludingLegacy() []MigrationDefinition {
	return []MigrationDefinition{
		Migration_20240101_002_CreateStaffsTable(),
		Migration_20240101_003_CreatePatientsTable(),
		Migration_20240101_005_CreateTenantsTable(),
		Migration_20240101_007_AddHospitalFieldsToTenants(),
	}
}
