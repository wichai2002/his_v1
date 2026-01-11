package services

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/wichai2002/his_v1/internal/domain"
	"github.com/wichai2002/his_v1/internal/infrastructure/database"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type tenantService struct {
	tenantRepo domain.TenantRepository
	dbManager  *database.TenantDBManager
	db         *gorm.DB
}

// NewTenantService creates a new tenant service
func NewTenantService(
	tenantRepo domain.TenantRepository,
	dbManager *database.TenantDBManager,
	db *gorm.DB,
) domain.TenantService {
	return &tenantService{
		tenantRepo: tenantRepo,
		dbManager:  dbManager,
		db:         db,
	}
}

// GetBySubdomain retrieves a tenant by subdomain
func (s *tenantService) GetBySubdomain(subdomain string) (*domain.Tenant, error) {
	return s.tenantRepo.GetBySubdomain(subdomain)
}

// GetBySchemaName retrieves a tenant by schema name
func (s *tenantService) GetBySchemaName(schemaName string) (*domain.Tenant, error) {
	return s.tenantRepo.GetBySchemaName(schemaName)
}

// CreateTenant creates a new tenant record
func (s *tenantService) CreateTenant(req *domain.TenantCreateRequest) (*domain.Tenant, error) {
	// Validate and sanitize schema name
	schemaName := sanitizeSchemaName(req.SchemaName)
	if schemaName == "" {
		return nil, fmt.Errorf("invalid schema name")
	}

	// Check if schema name already exists
	exists, err := s.tenantRepo.SchemaExists(schemaName)
	if err != nil {
		return nil, fmt.Errorf("failed to check schema existence: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("schema name already exists: %s", schemaName)
	}

	// Validate hospital code length (must be exactly 8 characters)
	if len(req.HospitalCode) != 8 {
		return nil, fmt.Errorf("hospital code must be exactly 8 characters")
	}

	// Validate hospital name length
	if len(req.HospitalName) < 1 || len(req.HospitalName) > 150 {
		return nil, fmt.Errorf("hospital name must be between 1 and 150 characters")
	}

	tenant := &domain.Tenant{
		TenantCode:   req.TenantCode,
		Name:         req.Name,
		SchemaName:   schemaName,
		Subdomain:    strings.ToLower(req.Subdomain),
		IsActive:     true,
		HospitalName: req.HospitalName,
		HospitalCode: req.HospitalCode,
		Address:      req.Address,
		HNRunning:    0,
	}

	if err := s.tenantRepo.Create(tenant); err != nil {
		return nil, fmt.Errorf("failed to create tenant: %w", err)
	}

	return tenant, nil
}

// CreateTenantSchema creates the PostgreSQL schema for a tenant
func (s *tenantService) CreateTenantSchema(schemaName string) error {
	schemaName = sanitizeSchemaName(schemaName)
	if schemaName == "" {
		return fmt.Errorf("invalid schema name")
	}
	return s.dbManager.CreateSchema(schemaName)
}

// MigrateTenantSchema runs migrations for a specific tenant schema
func (s *tenantService) MigrateTenantSchema(schemaName string) error {
	schemaName = sanitizeSchemaName(schemaName)
	if schemaName == "" {
		return fmt.Errorf("invalid schema name")
	}

	// Execute migrations within the tenant schema
	return s.dbManager.ExecuteInSchema(schemaName, func(tx *gorm.DB) error {
		// Create tables in tenant schema
		// Note: We create tenant-specific tables (Staff, Patient)
		// Tenant table remains in public schema
		if err := tx.AutoMigrate(&domain.Staff{}, &domain.Patient{}); err != nil {
			return fmt.Errorf("failed to migrate tenant schema: %w", err)
		}
		return nil
	})
}

// SetupTenantWithAdmin creates a complete tenant setup including schema, migrations, and admin user
func (s *tenantService) SetupTenantWithAdmin(
	tenantCode string,
	name string,
	subdomain string,
	hospitalName string,
	hospitalCode string,
	address *string,
	adminUsername string,
	adminPassword string,
	adminEmail string,
) (*domain.Tenant, error) {
	// Generate schema name from tenant code
	schemaName := generateSchemaName(tenantCode)

	// Validate hospital code length (must be exactly 8 characters)
	if len(hospitalCode) != 8 {
		return nil, fmt.Errorf("hospital code must be exactly 8 characters")
	}

	// Validate hospital name length
	if len(hospitalName) < 1 || len(hospitalName) > 150 {
		return nil, fmt.Errorf("hospital name must be between 1 and 150 characters")
	}

	// Start transaction
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. Create tenant record in public schema
	tenant := &domain.Tenant{
		TenantCode:   tenantCode,
		Name:         name,
		SchemaName:   schemaName,
		Subdomain:    strings.ToLower(subdomain),
		IsActive:     true,
		HospitalName: hospitalName,
		HospitalCode: hospitalCode,
		Address:      address,
		HNRunning:    0,
	}

	if err := tx.Create(tenant).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create tenant record: %w", err)
	}

	// 2. Create PostgreSQL schema
	if err := s.dbManager.CreateSchema(schemaName); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create schema: %w", err)
	}

	// 3. Migrate tenant schema (create tables)
	if err := s.migrateTenantSchemaWithDB(tx, schemaName); err != nil {
		tx.Rollback()
		// Try to clean up the schema
		_ = s.dbManager.DropSchema(schemaName, true)
		return nil, fmt.Errorf("failed to migrate tenant schema: %w", err)
	}

	// 4. Create admin user in tenant schema
	if err := s.createAdminInSchema(tx, schemaName, tenant, adminUsername, adminPassword, adminEmail); err != nil {
		tx.Rollback()
		// Try to clean up the schema
		_ = s.dbManager.DropSchema(schemaName, true)
		return nil, fmt.Errorf("failed to create admin user: %w", err)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		_ = s.dbManager.DropSchema(schemaName, true)
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return tenant, nil
}

// migrateTenantSchemaWithDB runs migrations for tenant schema using provided transaction
func (s *tenantService) migrateTenantSchemaWithDB(tx *gorm.DB, schemaName string) error {
	// Set search path to tenant schema
	if err := tx.Exec(fmt.Sprintf("SET search_path TO %s", schemaName)).Error; err != nil {
		return err
	}

	// Create Staff table in tenant schema
	staffTable := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s.staffs (
			id SERIAL PRIMARY KEY,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			deleted_at TIMESTAMP WITH TIME ZONE,
			username VARCHAR(100) UNIQUE NOT NULL,
			password VARCHAR(255) NOT NULL,
			staff_code VARCHAR(50) UNIQUE NOT NULL,
			phone_number VARCHAR(20) UNIQUE NOT NULL,
			email VARCHAR(255) UNIQUE NOT NULL,
			first_name VARCHAR(255) NOT NULL,
			last_name VARCHAR(255) NOT NULL,
			is_admin BOOLEAN DEFAULT FALSE
		)
	`, schemaName)
	if err := tx.Exec(staffTable).Error; err != nil {
		return fmt.Errorf("failed to create staffs table: %w", err)
	}

	// Create index on deleted_at
	staffIndex := fmt.Sprintf("CREATE INDEX IF NOT EXISTS idx_%s_staffs_deleted_at ON %s.staffs(deleted_at)", schemaName, schemaName)
	if err := tx.Exec(staffIndex).Error; err != nil {
		return fmt.Errorf("failed to create staffs index: %w", err)
	}

	// Create Patient table in tenant schema
	patientTable := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s.patients (
			id SERIAL PRIMARY KEY,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			deleted_at TIMESTAMP WITH TIME ZONE,
			first_name_th VARCHAR(255) NOT NULL,
			last_name_th VARCHAR(255) NOT NULL,
			middle_name_th VARCHAR(255),
			first_name_en VARCHAR(255) NOT NULL,
			last_name_en VARCHAR(255) NOT NULL,
			middle_name_en VARCHAR(255),
			date_of_birth DATE NOT NULL,
			nick_name_th VARCHAR(50),
			nick_name_en VARCHAR(50),
			patient_hn VARCHAR(50) UNIQUE NOT NULL,
			national_id VARCHAR(20) UNIQUE,
			passport_id VARCHAR(50) UNIQUE,
			phone_number VARCHAR(20) UNIQUE,
			email VARCHAR(255) UNIQUE,
			gender VARCHAR(10) NOT NULL,
			nationality VARCHAR(100) NOT NULL,
			blood_grp VARCHAR(5) NOT NULL
		)
	`, schemaName)
	if err := tx.Exec(patientTable).Error; err != nil {
		return fmt.Errorf("failed to create patients table: %w", err)
	}

	// Create index on deleted_at for patients
	patientIndex := fmt.Sprintf("CREATE INDEX IF NOT EXISTS idx_%s_patients_deleted_at ON %s.patients(deleted_at)", schemaName, schemaName)
	if err := tx.Exec(patientIndex).Error; err != nil {
		return fmt.Errorf("failed to create patients index: %w", err)
	}

	// Reset search path
	if err := tx.Exec("SET search_path TO public").Error; err != nil {
		return err
	}

	return nil
}

// createAdminInSchema creates an admin user in the specified tenant schema
func (s *tenantService) createAdminInSchema(
	tx *gorm.DB,
	schemaName string,
	tenant *domain.Tenant,
	adminUsername string,
	adminPassword string,
	adminEmail string,
) error {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Generate staff code from tenant code
	staffCode := "ADM" + tenant.TenantCode

	// Insert admin user into tenant schema
	insertSQL := fmt.Sprintf(`
		INSERT INTO %s.staffs (
			username, password, staff_code, phone_number, email,
			first_name, last_name, is_admin
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, schemaName)

	return tx.Exec(
		insertSQL,
		adminUsername,
		string(hashedPassword),
		staffCode,
		"",
		adminEmail,
		"Admin",
		tenant.Name,
		true,
	).Error
}

// generateSchemaName generates a valid schema name from tenant code
func generateSchemaName(tenantCode string) string {
	// Convert to lowercase and prefix with tenant_
	name := "tenant_" + strings.ToLower(tenantCode)
	return sanitizeSchemaName(name)
}

// sanitizeSchemaName ensures the schema name is valid for PostgreSQL
func sanitizeSchemaName(name string) string {
	// Convert to lowercase
	name = strings.ToLower(name)

	// Replace invalid characters with underscore
	reg := regexp.MustCompile(`[^a-z0-9_]`)
	name = reg.ReplaceAllString(name, "_")

	// Ensure it starts with a letter or underscore
	if len(name) > 0 && !((name[0] >= 'a' && name[0] <= 'z') || name[0] == '_') {
		name = "_" + name
	}

	// Truncate to max 63 characters (PostgreSQL limit)
	if len(name) > 63 {
		name = name[:63]
	}

	return name
}

// GenerateHN generates a new HN in format 'hospitalCode-HNRunning'
// Example: "HOSP0001-00000001"
func (s *tenantService) GenerateHN(schemaName string) (string, error) {
	// Get tenant by schema name
	tenant, err := s.tenantRepo.GetBySchemaName(schemaName)
	if err != nil {
		return "", fmt.Errorf("failed to get tenant: %w", err)
	}

	// Atomically increment HN running number
	newHN, err := s.tenantRepo.IncrementHNRunning(tenant.ID)
	if err != nil {
		return "", fmt.Errorf("failed to increment HN running: %w", err)
	}

	// Format: hospitalCode-HNRunning (8 digits padded)
	hn := fmt.Sprintf("%s-%08d", tenant.HospitalCode, newHN)

	return hn, nil
}
