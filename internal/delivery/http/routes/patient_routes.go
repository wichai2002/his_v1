package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/wichai2002/his_v1/internal/delivery/http/handler"
	"github.com/wichai2002/his_v1/internal/delivery/http/middleware"
	"github.com/wichai2002/his_v1/pkg/jwt"
)

// RegisterPatientRoutes registers all patient-related routes
func RegisterPatientRoutes(router *gin.RouterGroup, patientHandler *handler.PatientHandler, jwtService jwt.JWTService) {
	patientGroup := router.Group("/patient")
	patientGroup.Use(middleware.AuthMiddleware(jwtService))
	{
		patientGroup.GET("/search/:id", patientHandler.Search)
		patientGroup.POST("/create", patientHandler.Create)
		patientGroup.PUT("/update/:id", patientHandler.Update)
		patientGroup.PATCH("/update/:id", patientHandler.PartialUpdate)
		patientGroup.DELETE("/delete/:id", patientHandler.Delete)
	}
}
