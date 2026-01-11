package repository

import (
	"fmt"

	"github.com/wichai2002/his_v1/internal/domain"
	"github.com/wichai2002/his_v1/internal/infrastructure/database"
	"gorm.io/gorm"
)

type patientRepository struct {
	*TenantAwareRepository
}

// NewPatientRepository creates a new patient repository
func NewPatientRepository(db *gorm.DB, dbManager *database.TenantDBManager) domain.PatientRepository {
	return &patientRepository{
		TenantAwareRepository: NewTenantAwareRepository(db, dbManager),
	}
}

// getDB returns the appropriate database based on schema
func (r *patientRepository) getDB(schemaName string) (*gorm.DB, error) {
	if schemaName == "" || schemaName == "public" {
		return r.GetDB(), nil
	}
	return r.GetTenantDB(schemaName)
}

func (r *patientRepository) GetAll(schemaName string) ([]domain.Patient, error) {
	db, err := r.getDB(schemaName)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant db: %w", err)
	}

	var patients []domain.Patient
	if err := db.Find(&patients).Error; err != nil {
		return nil, err
	}
	return patients, nil
}

func (r *patientRepository) GetByID(id uint, schemaName string) (*domain.Patient, error) {
	db, err := r.getDB(schemaName)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant db: %w", err)
	}

	var patient domain.Patient
	if err := db.First(&patient, id).Error; err != nil {
		return nil, err
	}
	return &patient, nil
}

// SearchByID searches patient by ID, national ID, or passport ID
func (r *patientRepository) SearchByID(id uint, schemaName string) (*domain.Patient, error) {
	db, err := r.getDB(schemaName)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant db: %w", err)
	}

	var patient domain.Patient
	if err := db.Where("id = ?", id).
		Or("national_id = ?", id).
		Or("passport_id = ?", id).
		First(&patient).Error; err != nil {
		return nil, err
	}
	return &patient, nil
}

// Search patient by query: first name, last name, middle name, patient HN, national ID, passport ID, phone number
func (r *patientRepository) Search(query string, schemaName string) ([]domain.Patient, error) {
	db, err := r.getDB(schemaName)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant db: %w", err)
	}

	var patients []domain.Patient

	if query == "" {
		// Return all patients if no query
		if err := db.Find(&patients).Error; err != nil {
			return nil, err
		}
		return patients, nil
	}

	searchQuery := "%" + query + "%"
	if err := db.Where("first_name_th ILIKE ?", searchQuery).
		Or("last_name_th ILIKE ?", searchQuery).
		Or("first_name_en ILIKE ?", searchQuery).
		Or("last_name_en ILIKE ?", searchQuery).
		Or("patient_hn ILIKE ?", searchQuery).
		Or("national_id ILIKE ?", searchQuery).
		Or("passport_id ILIKE ?", searchQuery).
		Or("phone_number ILIKE ?", searchQuery).
		Find(&patients).Error; err != nil {
		return nil, err
	}

	return patients, nil
}

func (r *patientRepository) Create(patient *domain.Patient, schemaName string) error {
	db, err := r.getDB(schemaName)
	if err != nil {
		return fmt.Errorf("failed to get tenant db: %w", err)
	}
	return db.Create(patient).Error
}

func (r *patientRepository) Update(patient *domain.Patient, schemaName string) error {
	db, err := r.getDB(schemaName)
	if err != nil {
		return fmt.Errorf("failed to get tenant db: %w", err)
	}
	return db.Save(patient).Error
}

func (r *patientRepository) PartialUpdate(id uint, updates map[string]interface{}, schemaName string) error {
	db, err := r.getDB(schemaName)
	if err != nil {
		return fmt.Errorf("failed to get tenant db: %w", err)
	}

	// Check if record exists and update in one operation
	result := db.Model(&domain.Patient{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}

	// Check if any rows were affected (record exists)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *patientRepository) Delete(id uint, schemaName string) error {
	db, err := r.getDB(schemaName)
	if err != nil {
		return fmt.Errorf("failed to get tenant db: %w", err)
	}

	result := db.Delete(&domain.Patient{}, id)
	if result.Error != nil {
		return result.Error
	}

	// Check if any rows were affected (record exists)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
