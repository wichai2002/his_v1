package repository

import (
	"fmt"
	"regexp"

	"github.com/wichai2002/his_v1/internal/domain"
	"github.com/wichai2002/his_v1/internal/infrastructure/database"
	"gorm.io/gorm"
)

// schemaNameRegex validates schema names to prevent SQL injection
// Only allows alphanumeric characters and underscores, must start with a letter or underscore
var schemaNameRegex = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)

// TenantAwareRepository provides base functionality for tenant-scoped repositories
type TenantAwareRepository struct {
	baseDB    *gorm.DB
	dbManager *database.TenantDBManager
}

// NewTenantAwareRepository creates a new tenant-aware repository base
func NewTenantAwareRepository(db *gorm.DB, dbManager *database.TenantDBManager) *TenantAwareRepository {
	return &TenantAwareRepository{
		baseDB:    db,
		dbManager: dbManager,
	}
}

// GetDB returns the base database (public schema)
func (r *TenantAwareRepository) GetDB() *gorm.DB {
	return r.baseDB
}

// isValidSchemaName validates the schema name to prevent SQL injection
func isValidSchemaName(schemaName string) bool {
	if len(schemaName) == 0 || len(schemaName) > 63 { // PostgreSQL identifier limit
		return false
	}
	return schemaNameRegex.MatchString(schemaName)
}

// GetTenantDB returns a database session scoped to the specified schema
// It uses SET LOCAL to ensure the search_path only affects the current transaction
func (r *TenantAwareRepository) GetTenantDB(schemaName string) (*gorm.DB, error) {
	if schemaName == "" || schemaName == "public" {
		return r.baseDB, nil
	}

	// Validate schema name to prevent SQL injection
	if !isValidSchemaName(schemaName) {
		return nil, fmt.Errorf("%w: %s", domain.ErrInvalidSchemaName, schemaName)
	}

	// Create new session and set search_path using SET LOCAL
	// SET LOCAL ensures the setting only affects the current transaction,
	// preventing potential data leakage in connection pooling scenarios
	db := r.baseDB.Session(&gorm.Session{NewDB: true})
	if err := db.Exec(fmt.Sprintf("SET LOCAL search_path TO %s, public", schemaName)).Error; err != nil {
		return nil, fmt.Errorf("failed to set search_path: %w", err)
	}

	return db, nil
}

// ExecuteInSchema executes a function within a specific schema context using a transaction
// This ensures proper isolation and automatic cleanup of search_path
func (r *TenantAwareRepository) ExecuteInSchema(schemaName string, fn func(db *gorm.DB) error) error {
	if schemaName == "" || schemaName == "public" {
		return fn(r.baseDB)
	}

	// Validate schema name
	if !isValidSchemaName(schemaName) {
		return fmt.Errorf("%w: %s", domain.ErrInvalidSchemaName, schemaName)
	}

	// Use transaction to ensure search_path is properly scoped and cleaned up
	return r.baseDB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(fmt.Sprintf("SET LOCAL search_path TO %s, public", schemaName)).Error; err != nil {
			return fmt.Errorf("failed to set search_path: %w", err)
		}
		return fn(tx)
	})
}
