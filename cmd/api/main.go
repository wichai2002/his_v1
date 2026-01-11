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

	// Run migrations
	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize JWT service
	jwtService := jwt.NewJWTService(cfg.JWT.SecretKey, cfg.JWT.ExpiresIn)

	// Initialize repositories
	hospitalRepo := repository.NewHospitalRepository(db)
	staffRepo := repository.NewStaffRepository(db)
	patientRepo := repository.NewPatientRepository(db)

	// Initialize services
	hospitalService := services.NewHospitalService(hospitalRepo)
	staffService := services.NewStaffService(staffRepo, jwtService, hospitalService)
	patientService := services.NewPatientService(patientRepo, hospitalService)

	// Initialize handlers
	hospitalHandler := handler.NewHospitalHandler(hospitalService)
	staffHandler := handler.NewStaffHandler(staffService)
	patientHandler := handler.NewPatientHandler(patientService)

	// Setup router
	router := http.NewRouter(
		staffHandler,
		hospitalHandler,
		patientHandler,
		jwtService,
	)
	engine := router.Setup()

	// Start server
	log.Printf("Server starting on port %s", cfg.Server.Port)
	if err := engine.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
