package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/wichai2002/his_v1/internal/delivery/http/handler"
	"github.com/wichai2002/his_v1/internal/delivery/http/middleware"
	"github.com/wichai2002/his_v1/pkg/jwt"
)

// RegisterStaffRoutes registers all staff-related routes
func RegisterStaffRoutes(router *gin.RouterGroup, staffHandler *handler.StaffHandler, jwtService jwt.JWTService) {
	staffGroup := router.Group("/staff")
	{
		// Public routes
		staffGroup.POST("/login", staffHandler.Login)

		// Protected routes
		protected := staffGroup.Group("")
		protected.Use(middleware.AuthMiddleware(jwtService))
		{
			protected.POST("/logout", staffHandler.Logout)
			protected.GET("/", staffHandler.GetAll)
			protected.GET("/:id", staffHandler.GetByID)
			protected.PUT("/update/:id", staffHandler.Update)
		}

		// Admin only routes
		admin := staffGroup.Group("")
		admin.Use(middleware.AuthMiddleware(jwtService), middleware.AdminMiddleware())
		{
			admin.POST("/create", staffHandler.Create)
			admin.DELETE("/delete/:id", staffHandler.Delete)
		}
	}
}
