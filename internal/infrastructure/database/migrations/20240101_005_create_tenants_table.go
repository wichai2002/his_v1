package migrations

import (
	"github.com/wichai2002/his_v1/internal/domain"
	"gorm.io/gorm"
)

// Migration_20240101_005_CreateTenantsTable creates the tenants table in public schema
func Migration_20240101_005_CreateTenantsTable() MigrationDefinition {
	return MigrationDefinition{
		Version: "20240101_005",
		Name:    "create_tenants_table",
		Up: func(db *gorm.DB) error {
			return db.AutoMigrate(&domain.Tenant{})
		},
		Down: func(db *gorm.DB) error {
			return db.Migrator().DropTable(&domain.Tenant{})
		},
	}
}
