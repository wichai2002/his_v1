package domain

import (
	"gorm.io/gorm"
)

// Tenant represents a tenant with their schema mapping
type Tenant struct {
	gorm.Model
	TenantCode   string `json:"tenant_code" gorm:"uniqueIndex;not null"`
	Name         string `json:"name" gorm:"not null"`
	SchemaName   string `json:"schema_name" gorm:"uniqueIndex;not null"`
	Subdomain    string `json:"subdomain" gorm:"uniqueIndex;not null"`
	IsActive     bool   `json:"is_active" gorm:"default:true"`
	DatabaseHost string `json:"database_host" gorm:"default:''"` // For future multi-database support

	// Hospital settings
	HospitalName string  `json:"hospital_name" gorm:"not null;size:150" binding:"required,min=1,max=150"`
	HospitalCode string  `json:"hospital_code" gorm:"not null;size:8" binding:"required,len=8"`
	Address      *string `json:"address" gorm:"type:text"` // Can be null
	HNRunning    uint64  `json:"hn_running" gorm:"not null;default:0"`
}

// TenantInfo holds the tenant context information for requests
type TenantInfo struct {
	TenantID     uint
	TenantCode   string
	SchemaName   string
	Subdomain    string
	HospitalCode string
}

// TenantCreateRequest represents the request to create a new tenant
type TenantCreateRequest struct {
	TenantCode   string  `json:"tenant_code" binding:"required"`
	Name         string  `json:"name" binding:"required"`
	SchemaName   string  `json:"schema_name" binding:"required"`
	Subdomain    string  `json:"subdomain" binding:"required"`
	HospitalName string  `json:"hospital_name" binding:"required,min=1,max=150"`
	HospitalCode string  `json:"hospital_code" binding:"required,len=8"`
	Address      *string `json:"address"`
}

// TenantRepository interface for tenant operations
type TenantRepository interface {
	GetBySubdomain(subdomain string) (*Tenant, error)
	GetBySchemaName(schemaName string) (*Tenant, error)
	Create(tenant *Tenant) error
	Update(tenant *Tenant) error
	Delete(id uint) error
	SchemaExists(schemaName string) (bool, error)
	// IncrementHNRunning atomically increments HNRunning and returns the new value
	IncrementHNRunning(tenantID uint) (uint64, error)
}

// TenantService interface for tenant business logic
type TenantService interface {
	GetBySubdomain(subdomain string) (*Tenant, error)
	GetBySchemaName(schemaName string) (*Tenant, error)
	CreateTenant(req *TenantCreateRequest) (*Tenant, error)
	CreateTenantSchema(schemaName string) error
	MigrateTenantSchema(schemaName string) error
	SetupTenantWithAdmin(tenantCode, name, subdomain, hospitalName, hospitalCode string, address *string, adminUsername, adminPassword, adminEmail string) (*Tenant, error)
	// GenerateHN generates a new HN in format 'hospitalCode-HNRunning'
	GenerateHN(schemaName string) (string, error)
}
