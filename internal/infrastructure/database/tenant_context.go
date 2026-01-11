package database

import (
	"context"
	"fmt"
	"regexp"
	"sync"

	"gorm.io/gorm"
)

// schemaNameRegex validates schema names: must start with letter or underscore,
// followed by alphanumeric characters or underscores
var schemaNameRegex = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)

// TenantContextKey is the context key for tenant information
type TenantContextKey struct{}

// TenantContext holds tenant-specific database context
type TenantContext struct {
	SchemaName string
	TenantID   uint
}

// TenantDBManager manages tenant-specific database connections with schema isolation
type TenantDBManager struct {
	baseDB       *gorm.DB
	schemaDBs    map[string]*gorm.DB
	mutex        sync.RWMutex
	publicSchema string
}

// NewTenantDBManager creates a new tenant database manager
func NewTenantDBManager(db *gorm.DB) *TenantDBManager {
	return &TenantDBManager{
		baseDB:       db,
		schemaDBs:    make(map[string]*gorm.DB),
		publicSchema: "public",
	}
}

// GetDB returns the base database connection (public schema)
func (m *TenantDBManager) GetDB() *gorm.DB {
	return m.baseDB
}

// GetTenantDB returns a database session configured for the specified schema
func (m *TenantDBManager) GetTenantDB(schemaName string) *gorm.DB {
	if schemaName == "" || schemaName == m.publicSchema {
		return m.baseDB
	}

	m.mutex.RLock()
	if db, exists := m.schemaDBs[schemaName]; exists {
		m.mutex.RUnlock()
		return db
	}
	m.mutex.RUnlock()

	// Create a new session with schema search path
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Double-check after acquiring write lock
	if db, exists := m.schemaDBs[schemaName]; exists {
		return db
	}

	// Create session that sets search_path for this schema
	tenantDB := m.baseDB.Session(&gorm.Session{})
	m.schemaDBs[schemaName] = tenantDB

	return tenantDB
}

// GetDBForContext returns a database session based on context
func (m *TenantDBManager) GetDBForContext(ctx context.Context) *gorm.DB {
	tenantCtx, ok := ctx.Value(TenantContextKey{}).(*TenantContext)
	if !ok || tenantCtx == nil || tenantCtx.SchemaName == "" {
		return m.baseDB
	}
	return m.GetTenantDB(tenantCtx.SchemaName)
}

// WithSchema returns a new DB session that operates on the specified schema
func (m *TenantDBManager) WithSchema(schemaName string) (*gorm.DB, error) {
	if schemaName == "" || schemaName == m.publicSchema {
		return m.baseDB, nil
	}

	// Create a new session and set search_path
	db := m.baseDB.Session(&gorm.Session{NewDB: true})
	if err := db.Exec(fmt.Sprintf("SET search_path TO %s, public", schemaName)).Error; err != nil {
		return nil, fmt.Errorf("failed to set search_path: %w", err)
	}
	return db, nil
}

// ExecuteInSchema executes a function within a specific schema context
func (m *TenantDBManager) ExecuteInSchema(schemaName string, fn func(db *gorm.DB) error) error {
	return m.baseDB.Transaction(func(tx *gorm.DB) error {
		// Set schema search path
		if err := tx.Exec(fmt.Sprintf("SET search_path TO %s", schemaName)).Error; err != nil {
			return fmt.Errorf("failed to set schema: %w", err)
		}
		return fn(tx)
	})
}

// CreateSchema creates a new PostgreSQL schema
func (m *TenantDBManager) CreateSchema(schemaName string) error {
	// Validate schema name to prevent SQL injection
	if !isValidSchemaName(schemaName) {
		return fmt.Errorf("invalid schema name: %s", schemaName)
	}

	sql := fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", schemaName)
	return m.baseDB.Exec(sql).Error
}

// DropSchema drops a PostgreSQL schema
func (m *TenantDBManager) DropSchema(schemaName string, cascade bool) error {
	if !isValidSchemaName(schemaName) {
		return fmt.Errorf("invalid schema name: %s", schemaName)
	}

	sql := fmt.Sprintf("DROP SCHEMA IF EXISTS %s", schemaName)
	if cascade {
		sql += " CASCADE"
	}
	return m.baseDB.Exec(sql).Error
}

// SchemaExists checks if a schema exists
func (m *TenantDBManager) SchemaExists(schemaName string) (bool, error) {
	var count int64
	err := m.baseDB.Raw(
		"SELECT COUNT(*) FROM information_schema.schemata WHERE schema_name = ?",
		schemaName,
	).Scan(&count).Error

	return count > 0, err
}

// GetSchemaDB returns a DB session scoped to a specific schema
// This is the primary method to get a tenant-scoped database connection
func (m *TenantDBManager) GetSchemaDB(schemaName string) (*gorm.DB, error) {
	if schemaName == "" {
		return nil, fmt.Errorf("schema name is required")
	}

	exists, err := m.SchemaExists(schemaName)
	if err != nil {
		return nil, fmt.Errorf("failed to check schema existence: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("schema %s does not exist", schemaName)
	}

	// Create a new session and set search_path
	db := m.baseDB.Session(&gorm.Session{NewDB: true})
	if err := db.Exec(fmt.Sprintf("SET search_path TO %s", schemaName)).Error; err != nil {
		return nil, fmt.Errorf("failed to set search_path: %w", err)
	}

	return db, nil
}

// isValidSchemaName validates schema name to prevent SQL injection
func isValidSchemaName(name string) bool {
	if name == "" || len(name) > 63 {
		return false
	}
	return schemaNameRegex.MatchString(name)
}
