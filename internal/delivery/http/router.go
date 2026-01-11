package http

import (
	"github.com/gin-gonic/gin"
	"github.com/wichai2002/his_v1/internal/delivery/http/handler"
	"github.com/wichai2002/his_v1/internal/delivery/http/middleware"
	"github.com/wichai2002/his_v1/internal/delivery/http/routes"
	"github.com/wichai2002/his_v1/internal/domain"
	"github.com/wichai2002/his_v1/internal/infrastructure/database"
	"github.com/wichai2002/his_v1/pkg/jwt"
)

type Router struct {
	staffHandler   *handler.StaffHandler
	patientHandler *handler.PatientHandler
	jwtService     jwt.JWTService
	tenantService  domain.TenantService
	dbManager      *database.TenantDBManager
}

func NewRouter(
	staffHandler *handler.StaffHandler,
	patientHandler *handler.PatientHandler,
	jwtService jwt.JWTService,
	tenantService domain.TenantService,
	dbManager *database.TenantDBManager,
) *Router {
	return &Router{
		staffHandler:   staffHandler,
		patientHandler: patientHandler,
		jwtService:     jwtService,
		tenantService:  tenantService,
		dbManager:      dbManager,
	}
}

func (r *Router) Setup() *gin.Engine {
	router := gin.Default()

	// Health check - no tenant required
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API v1 routes
	routerV1 := router.Group("/api/v1")

	// Apply tenant middleware to all API routes
	// This will extract tenant from subdomain and set context
	routerV1.Use(middleware.TenantMiddleware(r.tenantService, r.dbManager))

	// Staff routes - some public (login), some require auth
	routes.RegisterStaffRoutes(routerV1, r.staffHandler, r.jwtService)

	// Protected routes (require tenant context and auth)
	routes.RegisterPatientRoutes(routerV1, r.patientHandler, r.jwtService)

	return router
}
