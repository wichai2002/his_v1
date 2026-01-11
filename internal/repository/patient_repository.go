package repository

import (
	"github.com/wichai2002/his_v1/internal/domain"

	"gorm.io/gorm"
)

type patientRepository struct {
	db *gorm.DB
}

func NewPatientRepository(db *gorm.DB) domain.PatientRepository {
	return &patientRepository{db: db}
}

func (r *patientRepository) GetAll() ([]domain.Patient, error) {
	var patients []domain.Patient
	if err := r.db.Preload("Hospital").Find(&patients).Error; err != nil {
		return nil, err
	}
	return patients, nil
}

func (r *patientRepository) GetByID(id uint) (*domain.Patient, error) {
	var patient domain.Patient
	if err := r.db.Preload("Hospital").First(&patient, id).Error; err != nil {
		return nil, err
	}
	return &patient, nil
}

// Search patient by ID, national ID, passport ID
func (r *patientRepository) SearchByID(id uint, hospitalID uint) (*domain.Patient, error) {
	var patient domain.Patient
	db := r.db.Preload("Hospital").Where("hospital_id = ?", hospitalID)

	if err := db.Where("id = ?", id).
		Or("national_id = ?", id).
		Or("passport_id = ?", id).
		First(&patient).Error; err != nil {
		return nil, err
	}
	return &patient, nil
}

// Search patient by query, first name, last name, middle name, patient HN, national ID, passport ID, phone number
func (r *patientRepository) Search(query string, hospitalID uint) ([]domain.Patient, error) {
	var patients []domain.Patient
	searchQuery := "%" + query + "%"

	db := r.db.Preload("Hospital").Where("hospital_id = ?", hospitalID)

	if query != "" {
		if err := db.Where(
			r.db.Where("first_name_th ILIKE ?", searchQuery).
				Or("last_name_th ILIKE ?", searchQuery).
				Or("first_name_en ILIKE ?", searchQuery).
				Or("last_name_en ILIKE ?", searchQuery).
				Or("patient_hn ILIKE ?", searchQuery).
				Or("national_id ILIKE ?", searchQuery).
				Or("passport_id ILIKE ?", searchQuery).
				Or("phone_number ILIKE ?", searchQuery),
		).Find(&patients).Error; err != nil {
			return nil, err
		}
	}

	return patients, nil
}

func (r *patientRepository) Create(patient *domain.Patient) error {
	return r.db.Create(patient).Error
}

func (r *patientRepository) Update(patient *domain.Patient) error {
	return r.db.Save(patient).Error
}

func (r *patientRepository) PartialUpdate(id uint, updates map[string]interface{}) error {
	return r.db.Model(&domain.Patient{}).Where("id = ?", id).Updates(updates).Error
}

func (r *patientRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Patient{}, id).Error
}
