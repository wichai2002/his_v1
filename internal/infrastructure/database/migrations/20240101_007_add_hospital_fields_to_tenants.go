package migrations

import (
	"gorm.io/gorm"
)

// Migration_20240101_007_AddHospitalFieldsToTenants adds hospital fields to tenants table
func Migration_20240101_007_AddHospitalFieldsToTenants() MigrationDefinition {
	return MigrationDefinition{
		Version: "20240101_007",
		Name:    "add_hospital_fields_to_tenants",
		Up: func(db *gorm.DB) error {
			// Add hospital_name column
			if !db.Migrator().HasColumn(&tenantTable{}, "hospital_name") {
				if err := db.Exec("ALTER TABLE tenants ADD COLUMN hospital_name VARCHAR(150) NOT NULL DEFAULT ''").Error; err != nil {
					return err
				}
			}

			// Add hospital_code column with unique constraint
			if !db.Migrator().HasColumn(&tenantTable{}, "hospital_code") {
				if err := db.Exec("ALTER TABLE tenants ADD COLUMN hospital_code VARCHAR(8) NOT NULL DEFAULT ''").Error; err != nil {
					return err
				}
			}

			// Add address column (nullable)
			if !db.Migrator().HasColumn(&tenantTable{}, "address") {
				if err := db.Exec("ALTER TABLE tenants ADD COLUMN address TEXT").Error; err != nil {
					return err
				}
			}

			// Add hn_running column
			if !db.Migrator().HasColumn(&tenantTable{}, "hn_running") {
				if err := db.Exec("ALTER TABLE tenants ADD COLUMN hn_running BIGINT NOT NULL DEFAULT 0").Error; err != nil {
					return err
				}
			}

			return nil
		},
		Down: func(db *gorm.DB) error {
			// Drop columns in reverse order
			if db.Migrator().HasColumn(&tenantTable{}, "hn_running") {
				if err := db.Exec("ALTER TABLE tenants DROP COLUMN hn_running").Error; err != nil {
					return err
				}
			}
			if db.Migrator().HasColumn(&tenantTable{}, "address") {
				if err := db.Exec("ALTER TABLE tenants DROP COLUMN address").Error; err != nil {
					return err
				}
			}
			if db.Migrator().HasColumn(&tenantTable{}, "hospital_code") {
				if err := db.Exec("DROP INDEX IF EXISTS idx_tenants_hospital_code").Error; err != nil {
					return err
				}
				if err := db.Exec("ALTER TABLE tenants DROP COLUMN hospital_code").Error; err != nil {
					return err
				}
			}
			if db.Migrator().HasColumn(&tenantTable{}, "hospital_name") {
				if err := db.Exec("ALTER TABLE tenants DROP COLUMN hospital_name").Error; err != nil {
					return err
				}
			}
			return nil
		},
	}
}

// tenantTable is used for migration column checking
type tenantTable struct{}

func (tenantTable) TableName() string {
	return "tenants"
}
