package repository

import (
	"github.com/wichai2002/his_v1/internal/domain"

	"gorm.io/gorm"
)

type staffRepository struct {
	db *gorm.DB
}

func NewStaffRepository(db *gorm.DB) domain.StaffRepository {
	return &staffRepository{db: db}
}

func (r *staffRepository) GetAll() ([]domain.Staff, error) {
	var staffs []domain.Staff
	if err := r.db.Preload("Hospital").Find(&staffs).Error; err != nil {
		return nil, err
	}
	return staffs, nil
}

func (r *staffRepository) GetByID(id uint) (*domain.Staff, error) {
	var staff domain.Staff
	if err := r.db.Preload("Hospital").First(&staff, id).Error; err != nil {
		return nil, err
	}
	return &staff, nil
}

func (r *staffRepository) GetByUsernameAndHospitalID(username string, hospitalID uint) (*domain.Staff, error) {
	var staff domain.Staff
	if err := r.db.Preload("Hospital").Where("username = ? AND hospital_id = ?", username, hospitalID).First(&staff).Error; err != nil {
		return nil, err
	}
	return &staff, nil
}

func (r *staffRepository) Create(staff *domain.Staff) error {
	return r.db.Create(staff).Error
}

func (r *staffRepository) Update(staff *domain.Staff) error {
	return r.db.Save(staff).Error
}

func (r *staffRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Staff{}, id).Error
}
