package repository

import (
	"github.com/wichai2002/his_v1/internal/domain"

	"gorm.io/gorm"
)

type hospitalRepository struct {
	db *gorm.DB
}

func NewHospitalRepository(db *gorm.DB) domain.HospitalRepository {
	return &hospitalRepository{db: db}
}

func (r *hospitalRepository) GetAllPublic() ([]domain.HospitalPublicResponse, error) {
	var hospitals []domain.HospitalPublicResponse
	if err := r.db.Model(&domain.Hospital{}).Select("id", "name").Find(&hospitals).Error; err != nil {
		return nil, err
	}
	return hospitals, nil
}

func (r *hospitalRepository) GetAll() ([]domain.Hospital, error) {
	var hospitals []domain.Hospital
	if err := r.db.Find(&hospitals).Error; err != nil {
		return nil, err
	}
	return hospitals, nil
}

func (r *hospitalRepository) GetByID(id uint) (*domain.Hospital, error) {
	var hospital domain.Hospital
	if err := r.db.First(&hospital, id).Error; err != nil {
		return nil, err
	}
	return &hospital, nil
}

func (r *hospitalRepository) Create(hospital *domain.Hospital) error {
	return r.db.Create(hospital).Error
}

func (r *hospitalRepository) Update(hospital *domain.Hospital) error {
	return r.db.Save(hospital).Error
}

func (r *hospitalRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Hospital{}, id).Error
}

func (r *hospitalRepository) IncrementHNRunningNumber(id uint) (int, error) {
	var hospital domain.Hospital
	if err := r.db.First(&hospital, id).Error; err != nil {
		return 0, err
	}

	hospital.HNRunningNumber++
	if err := r.db.Save(&hospital).Error; err != nil {
		return 0, err
	}

	return hospital.HNRunningNumber, nil
}

func (r *hospitalRepository) IncrementStaffRunningNumber(id uint) (int, error) {
	var hospital domain.Hospital
	if err := r.db.First(&hospital, id).Error; err != nil {
		return 0, err
	}

	hospital.StaffRunningNumber++
	if err := r.db.Save(&hospital).Error; err != nil {
		return 0, err
	}

	return hospital.StaffRunningNumber, nil
}
