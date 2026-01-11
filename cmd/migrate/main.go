package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/wichai2002/his_v1/config"
	"github.com/wichai2002/his_v1/internal/infrastructure/database"
)

func main() {
	// Define flags
	upCmd := flag.NewFlagSet("up", flag.ExitOnError)
	downCmd := flag.NewFlagSet("down", flag.ExitOnError)
	statusCmd := flag.NewFlagSet("status", flag.ExitOnError)
	resetCmd := flag.NewFlagSet("reset", flag.ExitOnError)

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to database
	db, err := database.NewPostgresDB(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	switch os.Args[1] {
	case "up":
		upCmd.Parse(os.Args[2:])
		fmt.Println("Running migrations...")
		if err := database.RunMigrations(db); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
		fmt.Println("Migrations completed successfully!")

	case "down":
		downCmd.Parse(os.Args[2:])
		fmt.Println("Rolling back last migration...")
		if err := database.RollbackMigration(db); err != nil {
			log.Fatalf("Rollback failed: %v", err)
		}
		fmt.Println("Rollback completed successfully!")

	case "reset":
		resetCmd.Parse(os.Args[2:])
		fmt.Println("Rolling back all migrations...")
		if err := database.RollbackAllMigrations(db); err != nil {
			log.Fatalf("Reset failed: %v", err)
		}
		fmt.Println("All migrations rolled back successfully!")

	case "status":
		statusCmd.Parse(os.Args[2:])
		database.MigrationStatus(db)

	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`
		HIS Migration Tool

		Usage:
		go run cmd/migrate/main.go <command>

		Commands:
		up      Run all pending migrations
		down    Rollback the last migration
		reset   Rollback all migrations
		status  Show migration status

		Examples:
		go run cmd/migrate/main.go up
		go run cmd/migrate/main.go down
		go run cmd/migrate/main.go status
		go run cmd/migrate/main.go reset
	`)
}
