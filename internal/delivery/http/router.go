package http

import (
	"github.com/gin-gonic/gin"
	"github.com/wichai2002/his_v1/internal/delivery/http/handler"
	"github.com/wichai2002/his_v1/internal/delivery/http/routes"
	"github.com/wichai2002/his_v1/pkg/jwt"
)

type Router struct {
	staffHandler    *handler.StaffHandler
	hospitalHandler *handler.HospitalHandler
	patientHandler  *handler.PatientHandler
	jwtService      jwt.JWTService
}

func NewRouter(
	staffHandler *handler.StaffHandler,
	hospitalHandler *handler.HospitalHandler,
	patientHandler *handler.PatientHandler,
	jwtService jwt.JWTService,
) *Router {
	return &Router{
		staffHandler:    staffHandler,
		hospitalHandler: hospitalHandler,
		patientHandler:  patientHandler,
		jwtService:      jwtService,
	}
}

func (r *Router) Setup() *gin.Engine {
	router := gin.Default()

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API v1 routes
	routerV1 := router.Group("/api/v1")
	routes.RegisterStaffRoutes(routerV1, r.staffHandler, r.jwtService)
	routes.RegisterHospitalPublicRoutes(routerV1, r.hospitalHandler)
	routes.RegisterHospitalRoutes(routerV1, r.hospitalHandler, r.jwtService)
	routes.RegisterPatientRoutes(routerV1, r.patientHandler, r.jwtService)

	return router
}
