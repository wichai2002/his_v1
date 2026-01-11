package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/wichai2002/his_v1/internal/domain"
	"github.com/wichai2002/his_v1/internal/infrastructure/database"
	"github.com/wichai2002/his_v1/pkg/utils"
)

const (
	// TenantContextKey is the key for tenant context in gin.Context
	TenantContextKey = "tenant_context"
	// TenantSchemaKey is the key for schema name
	TenantSchemaKey = "tenant_schema"
	// TenantIDKey is the key for tenant ID
	TenantIDKey = "tenant_id"
	// TenantSubdomainKey is the key for subdomain
	TenantSubdomainKey = "tenant_subdomain"
)

// TenantMiddleware extracts tenant information from subdomain and sets up database context
func TenantMiddleware(tenantService domain.TenantService, dbManager *database.TenantDBManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract subdomain from host
		subdomain := extractSubdomain(c.Request.Host)

		// Allow requests without subdomain to access public routes (e.g., hospital list)
		if subdomain == "" || subdomain == "www" || subdomain == "api" {
			// No tenant context - continue with public schema
			c.Set(TenantSchemaKey, "public")
			c.Next()
			return
		}

		// Look up tenant by subdomain
		tenant, err := tenantService.GetBySubdomain(subdomain)
		if err != nil {
			utils.ErrorResponse(c, http.StatusNotFound, "tenant not found for subdomain: "+subdomain)
			c.Abort()
			return
		}

		if !tenant.IsActive {
			utils.ErrorResponse(c, http.StatusForbidden, "tenant is inactive")
			c.Abort()
			return
		}

		// Create tenant context
		tenantCtx := &database.TenantContext{
			SchemaName: tenant.SchemaName,
			TenantID:   tenant.ID,
		}

		// Set tenant information in gin context
		c.Set(TenantContextKey, tenantCtx)
		c.Set(TenantSchemaKey, tenant.SchemaName)
		c.Set(TenantIDKey, tenant.ID)
		c.Set(TenantSubdomainKey, subdomain)

		// Add tenant context to request context
		ctx := context.WithValue(c.Request.Context(), database.TenantContextKey{}, tenantCtx)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

// TenantRequiredMiddleware ensures a tenant context is present
func TenantRequiredMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		schema, exists := c.Get(TenantSchemaKey)
		if !exists || schema == "" || schema == "public" {
			utils.ErrorResponse(c, http.StatusBadRequest, "tenant context required - please use subdomain to access this resource")
			c.Abort()
			return
		}
		c.Next()
	}
}

// OptionalTenantMiddleware sets tenant context if subdomain is present, but doesn't require it
func OptionalTenantMiddleware(tenantService domain.TenantService, dbManager *database.TenantDBManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		subdomain := extractSubdomain(c.Request.Host)

		if subdomain == "" || subdomain == "www" || subdomain == "api" {
			c.Set(TenantSchemaKey, "public")
			c.Next()
			return
		}

		tenant, err := tenantService.GetBySubdomain(subdomain)
		if err != nil {
			// Continue without tenant context
			c.Set(TenantSchemaKey, "public")
			c.Next()
			return
		}

		if tenant.IsActive {
			tenantCtx := &database.TenantContext{
				SchemaName: tenant.SchemaName,
				TenantID:   tenant.ID,
			}
			c.Set(TenantContextKey, tenantCtx)
			c.Set(TenantSchemaKey, tenant.SchemaName)
			c.Set(TenantIDKey, tenant.ID)
			c.Set(TenantSubdomainKey, subdomain)

			ctx := context.WithValue(c.Request.Context(), database.TenantContextKey{}, tenantCtx)
			c.Request = c.Request.WithContext(ctx)
		}

		c.Next()
	}
}

// extractSubdomain extracts the subdomain from the host header
// Supports formats like: tenant.example.com, tenant.localhost:8080
func extractSubdomain(host string) string {
	// Remove port if present
	if idx := strings.LastIndex(host, ":"); idx != -1 {
		host = host[:idx]
	}

	// Handle localhost specially
	if host == "localhost" || host == "127.0.0.1" {
		return ""
	}

	parts := strings.Split(host, ".")

	// Need at least 3 parts for subdomain.domain.tld
	// Or 2 parts for subdomain.localhost
	if len(parts) >= 3 {
		// Handle cases like tenant.example.com
		return parts[0]
	} else if len(parts) == 2 && parts[1] == "localhost" {
		// Handle cases like tenant.localhost
		return parts[0]
	}

	return ""
}

// GetTenantContext retrieves tenant context from gin context
func GetTenantContext(c *gin.Context) *database.TenantContext {
	if ctx, exists := c.Get(TenantContextKey); exists {
		if tenantCtx, ok := ctx.(*database.TenantContext); ok {
			return tenantCtx
		}
	}
	return nil
}

// GetTenantSchema retrieves the current tenant schema name
func GetTenantSchema(c *gin.Context) string {
	if schema, exists := c.Get(TenantSchemaKey); exists {
		if s, ok := schema.(string); ok {
			return s
		}
	}
	return "public"
}

// GetTenantID retrieves the current tenant ID
func GetTenantID(c *gin.Context) uint {
	if id, exists := c.Get(TenantIDKey); exists {
		if tenantID, ok := id.(uint); ok {
			return tenantID
		}
	}
	return 0
}

