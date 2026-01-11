package migrations

import (
	"github.com/wichai2002/his_v1/internal/domain"

	"gorm.io/gorm"
)

// Migration_20240101_003_CreatePatientsTable creates the patients table
func Migration_20240101_003_CreatePatientsTable() MigrationDefinition {
	return MigrationDefinition{
		Version: "20240101_003",
		Name:    "create_patients_table",
		Up: func(db *gorm.DB) error {
			return db.AutoMigrate(&domain.Patient{})
		},
		Down: func(db *gorm.DB) error {
			return db.Migrator().DropTable(&domain.Patient{})
		},
	}
}
