package migrations

import (
	"github.com/wichai2002/his_v1/internal/domain"

	"gorm.io/gorm"
)

// Migration_20240101_001_CreateHospitalsTable creates the hospitals table
func Migration_20240101_001_CreateHospitalsTable() MigrationDefinition {
	return MigrationDefinition{
		Version: "20240101_001",
		Name:    "create_hospitals_table",
		Up: func(db *gorm.DB) error {
			return db.AutoMigrate(&domain.Hospital{})
		},
		Down: func(db *gorm.DB) error {
			return db.Migrator().DropTable(&domain.Hospital{})
		},
	}
}
