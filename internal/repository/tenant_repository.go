package repository

import (
	"github.com/wichai2002/his_v1/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type tenantRepository struct {
	db *gorm.DB
}

// NewTenantRepository creates a new tenant repository
func NewTenantRepository(db *gorm.DB) domain.TenantRepository {
	return &tenantRepository{db: db}
}

// GetBySubdomain retrieves a tenant by its subdomain
func (r *tenantRepository) GetBySubdomain(subdomain string) (*domain.Tenant, error) {
	var tenant domain.Tenant
	if err := r.db.Where("subdomain = ? AND is_active = ?", subdomain, true).First(&tenant).Error; err != nil {
		return nil, err
	}
	return &tenant, nil
}

// GetBySchemaName retrieves a tenant by its schema name
func (r *tenantRepository) GetBySchemaName(schemaName string) (*domain.Tenant, error) {
	var tenant domain.Tenant
	if err := r.db.Where("schema_name = ?", schemaName).First(&tenant).Error; err != nil {
		return nil, err
	}
	return &tenant, nil
}

// Create creates a new tenant
func (r *tenantRepository) Create(tenant *domain.Tenant) error {
	return r.db.Create(tenant).Error
}

// Update updates an existing tenant
func (r *tenantRepository) Update(tenant *domain.Tenant) error {
	return r.db.Save(tenant).Error
}

// Delete deletes a tenant by ID
func (r *tenantRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Tenant{}, id).Error
}

// SchemaExists checks if a schema name already exists
func (r *tenantRepository) SchemaExists(schemaName string) (bool, error) {
	var count int64
	if err := r.db.Model(&domain.Tenant{}).Where("schema_name = ?", schemaName).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// IncrementHNRunning atomically increments HNRunning and returns the new value
func (r *tenantRepository) IncrementHNRunning(tenantID uint) (uint64, error) {
	var tenant domain.Tenant

	// Use transaction with row-level locking for atomic increment
	err := r.db.Transaction(func(tx *gorm.DB) error {
		// Lock the row for update
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&tenant, tenantID).Error; err != nil {
			return err
		}

		// Increment the HN running number
		tenant.HNRunning++

		// Save the updated value
		return tx.Save(&tenant).Error
	})

	if err != nil {
		return 0, err
	}

	return tenant.HNRunning, nil
}
