package migrations

import (
	"github.com/wichai2002/his_v1/internal/domain"

	"gorm.io/gorm"
)

// Migration_20240101_002_CreateStaffsTable creates the staffs table
func Migration_20240101_002_CreateStaffsTable() MigrationDefinition {
	return MigrationDefinition{
		Version: "20240101_002",
		Name:    "create_staffs_table",
		Up: func(db *gorm.DB) error {
			return db.AutoMigrate(&domain.Staff{})
		},
		Down: func(db *gorm.DB) error {
			return db.Migrator().DropTable(&domain.Staff{})
		},
	}
}
