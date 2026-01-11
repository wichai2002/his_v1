package middleware

import (
	"net/http"
	"strings"

	"github.com/wichai2002/his_v1/pkg/jwt"
	"github.com/wichai2002/his_v1/pkg/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates the JWT token
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

		// sets the user_id, username, is_admin, and hospital_id in the context for later use
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("is_admin", claims.IsAdmin)
		c.Set("hospital_id", claims.HospitalID)

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
