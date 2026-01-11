package migrations

import (
	"github.com/wichai2002/his_v1/internal/domain"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Migration_20240101_004_SeedDefaultData seeds the default hospital and admin user
func Migration_20240101_004_SeedDefaultData() MigrationDefinition {
	return MigrationDefinition{
		Version: "20240101_004",
		Name:    "seed_default_data",
		Up: func(db *gorm.DB) error {
			// Create default hospital
			hospital := &domain.Hospital{
				Name:            "Default Hospital",
				HospitalCode:    "HOS001",
				PhoneNumber:     "0000000000",
				Email:           "admin@hospital.com",
				Address:         "Bangkok, Thailand",
				HNRunningNumber: 0,
			}

			if err := db.Create(hospital).Error; err != nil {
				return err
			}

			// Create default admin user
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
			if err != nil {
				return err
			}

			admin := &domain.Staff{
				Username:    "admin",
				Password:    string(hashedPassword),
				StaffCode:   "ADM001",
				PhoneNumber: "0000000001",
				Email:       "admin@his.com",
				FirstName:   "System",
				LastName:    "Administrator",
				HospitalID:  hospital.ID,
				IsAdmin:     true,
			}

			return db.Create(admin).Error
		},
		Down: func(db *gorm.DB) error {
			// Remove seed data
			if err := db.Where("username = ?", "admin").Delete(&domain.Staff{}).Error; err != nil {
				return err
			}
			return db.Where("hospital_code = ?", "HOS001").Delete(&domain.Hospital{}).Error
		},
	}
}
