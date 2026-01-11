package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/wichai2002/his_v1/internal/delivery/http/handler"
	"github.com/wichai2002/his_v1/internal/delivery/http/middleware"
	"github.com/wichai2002/his_v1/pkg/jwt"
)

// RegisterHospitalRoutes registers all hospital-related routes
func RegisterHospitalRoutes(router *gin.RouterGroup, hospitalHandler *handler.HospitalHandler, jwtService jwt.JWTService) {
	hospitalGroup := router.Group("/hospital")
	hospitalGroup.Use(middleware.AuthMiddleware(jwtService))
	{
		// Protected routes (all authenticated users)
		hospitalGroup.GET("/", hospitalHandler.GetAll)
		hospitalGroup.GET("/:id", hospitalHandler.GetByID)

		// Admin only routes
		admin := hospitalGroup.Group("")
		admin.Use(middleware.AdminMiddleware())
		{
			admin.POST("/create", hospitalHandler.Create)
			admin.PUT("/update/:id", hospitalHandler.Update)
			admin.DELETE("/delete/:id", hospitalHandler.Delete)
		}
	}
}

func RegisterHospitalPublicRoutes(router *gin.RouterGroup, hospitalHandler *handler.HospitalHandler) {
	hospitalPublicGroup := router.Group("/public/hospital")
	{
		hospitalPublicGroup.GET("/", hospitalHandler.GetAllPublic)
	}
}
