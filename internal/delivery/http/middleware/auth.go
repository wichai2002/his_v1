package middleware

import (
	"net/http"
	"strings"

	"github.com/wichai2002/his_v1/pkg/jwt"
	"github.com/wichai2002/his_v1/pkg/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates the JWT token and sets user context
func AuthMiddleware(jwtService jwt.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// It checks if the Authorization header is present
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "authorization header required")
			c.Abort()
			return
		}

		// splits the Authorization header into the token type and the token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "invalid authorization header format")
			c.Abort()
			return
		}

		// validates token
		claims, err := jwtService.ValidateToken(parts[1])
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, "invalid or expired token")
			c.Abort()
			return
		}

		// Validate that JWT schema matches tenant schema from subdomain (if tenant context exists)
		if tenantSchema, exists := c.Get(TenantSchemaKey); exists {
			if ts, ok := tenantSchema.(string); ok && ts != "public" && ts != claims.SchemaName {
				utils.ErrorResponse(c, http.StatusForbidden, "token not valid for this tenant")
				c.Abort()
				return
			}
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("is_admin", claims.IsAdmin)
		c.Set("tenant_id", claims.TenantID)
		c.Set("jwt_schema", claims.SchemaName) // Schema from JWT for tenant context

		// If no tenant context from middleware, use JWT schema
		if _, exists := c.Get(TenantSchemaKey); !exists {
			c.Set(TenantSchemaKey, claims.SchemaName)
		}

		c.Next()
	}
}

// AdminMiddleware checks if the user is an admin
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// checks if the user is an admin
		isAdmin, exists := c.Get("is_admin")
		if !exists || !isAdmin.(bool) {
			utils.ErrorResponse(c, http.StatusForbidden, "admin access required")
			c.Abort()
			return
		}

		c.Next()
	}
}

// GetUserID retrieves the current user ID from context
func GetUserID(c *gin.Context) uint {
	if id, exists := c.Get("user_id"); exists {
		if userID, ok := id.(uint); ok {
			return userID
		}
	}
	return 0
}

// IsAdmin checks if the current user is an admin
func IsAdmin(c *gin.Context) bool {
	if isAdmin, exists := c.Get("is_admin"); exists {
		if admin, ok := isAdmin.(bool); ok {
			return admin
		}
	}
	return false
}
