package main

import (
	"log"

	"github.com/wichai2002/his_v1/config"
	"github.com/wichai2002/his_v1/internal/delivery/http"
	"github.com/wichai2002/his_v1/internal/delivery/http/handler"
	"github.com/wichai2002/his_v1/internal/infrastructure/database"
	"github.com/wichai2002/his_v1/internal/repository"
	"github.com/wichai2002/his_v1/internal/services"
	"github.com/wichai2002/his_v1/pkg/jwt"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	db, err := database.NewPostgresDB(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run migrations (for public schema)
	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize tenant database manager
	dbManager := database.NewTenantDBManager(db)

	// Initialize JWT service
	jwtService := jwt.NewJWTService(cfg.JWT.SecretKey, cfg.JWT.ExpiresIn)

	// Initialize repositories
	tenantRepo := repository.NewTenantRepository(db)
	staffRepo := repository.NewStaffRepository(db, dbManager)
	patientRepo := repository.NewPatientRepository(db, dbManager)

	// Initialize services
	tenantService := services.NewTenantService(tenantRepo, dbManager, db)
	staffService := services.NewStaffService(staffRepo, jwtService)
	patientService := services.NewPatientService(patientRepo, tenantService)

	// Initialize handlers
	staffHandler := handler.NewStaffHandler(staffService)
	patientHandler := handler.NewPatientHandler(patientService)

	// Setup router with tenant support
	router := http.NewRouter(
		staffHandler,
		patientHandler,
		jwtService,
		tenantService,
		dbManager,
	)
	engine := router.Setup()

	// Start server
	log.Printf("Server starting on port %s", cfg.Server.Port)
	log.Println("Multi-tenant mode enabled - use subdomain to access tenant data")
	if err := engine.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
