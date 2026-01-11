package repository

import (
	"fmt"

	"github.com/wichai2002/his_v1/internal/domain"
	"github.com/wichai2002/his_v1/internal/infrastructure/database"
	"gorm.io/gorm"
)

type staffRepository struct {
	*TenantAwareRepository
}

// NewStaffRepository creates a new staff repository
func NewStaffRepository(db *gorm.DB, dbManager *database.TenantDBManager) domain.StaffRepository {
	return &staffRepository{
		TenantAwareRepository: NewTenantAwareRepository(db, dbManager),
	}
}

// getDB returns the appropriate database based on schema
func (r *staffRepository) getDB(schemaName string) (*gorm.DB, error) {
	if schemaName == "" || schemaName == "public" {
		return r.GetDB(), nil
	}
	return r.GetTenantDB(schemaName)
}

func (r *staffRepository) GetAll(schemaName string) ([]domain.Staff, error) {
	db, err := r.getDB(schemaName)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant db: %w", err)
	}

	var staffs []domain.Staff
	if err := db.Find(&staffs).Error; err != nil {
		return nil, err
	}
	return staffs, nil
}

func (r *staffRepository) GetByID(id uint, schemaName string) (*domain.Staff, error) {
	db, err := r.getDB(schemaName)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant db: %w", err)
	}

	var staff domain.Staff
	if err := db.First(&staff, id).Error; err != nil {
		return nil, err
	}
	return &staff, nil
}

// GetByUsername finds staff by username (tenant schema provides isolation)
func (r *staffRepository) GetByUsername(username string, schemaName string) (*domain.Staff, error) {
	db, err := r.getDB(schemaName)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant db: %w", err)
	}

	var staff domain.Staff
	if err := db.Where("username = ?", username).First(&staff).Error; err != nil {
		return nil, err
	}
	return &staff, nil
}

func (r *staffRepository) Create(staff *domain.Staff, schemaName string) error {
	db, err := r.getDB(schemaName)
	if err != nil {
		return fmt.Errorf("failed to get tenant db: %w", err)
	}
	return db.Create(staff).Error
}

func (r *staffRepository) Update(staff *domain.Staff, schemaName string) error {
	db, err := r.getDB(schemaName)
	if err != nil {
		return fmt.Errorf("failed to get tenant db: %w", err)
	}
	return db.Save(staff).Error
}

func (r *staffRepository) Delete(id uint, schemaName string) error {
	db, err := r.getDB(schemaName)
	if err != nil {
		return fmt.Errorf("failed to get tenant db: %w", err)
	}
	return db.Delete(&domain.Staff{}, id).Error
}
