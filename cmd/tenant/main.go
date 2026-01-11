package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"

	"github.com/wichai2002/his_v1/config"
	"github.com/wichai2002/his_v1/internal/domain"
	"github.com/wichai2002/his_v1/internal/infrastructure/database"
	"github.com/wichai2002/his_v1/internal/repository"
	"github.com/wichai2002/his_v1/internal/services"
	"golang.org/x/term"
)

func main() {
	// Define subcommands
	createCmd := flag.NewFlagSet("create", flag.ExitOnError)
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)

	// Create command flags
	tenantCode := createCmd.String("code", "", "Tenant code (required)")
	tenantName := createCmd.String("name", "", "Tenant name (required)")
	subdomain := createCmd.String("subdomain", "", "Subdomain for the tenant (required)")
	hospitalName := createCmd.String("hospital-name", "", "Hospital name (required, max 150 chars)")
	hospitalCode := createCmd.String("hospital-code", "", "Hospital code (required, exactly 8 chars)")
	address := createCmd.String("address", "", "Hospital address (optional)")
	adminUsername := createCmd.String("admin-user", "", "Admin username (required, min 5 chars)")
	adminPassword := createCmd.String("admin-pass", "", "Admin password (optional, will prompt if not provided)")
	adminEmail := createCmd.String("admin-email", "", "Admin email (required)")

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "create":
		createCmd.Parse(os.Args[2:])
		runCreate(
			*tenantCode,
			*tenantName,
			*subdomain,
			*hospitalName,
			*hospitalCode,
			*address,
			*adminUsername,
			*adminPassword,
			*adminEmail,
		)

	case "list":
		listCmd.Parse(os.Args[2:])
		runList()

	default:
		printUsage()
		os.Exit(1)
	}
}

func runCreate(
	tenantCode, tenantName, subdomain,
	hospitalName, hospitalCode, address,
	adminUsername, adminPassword, adminEmail string,
) {
	// Validate required fields
	var errors []string

	if tenantCode == "" {
		errors = append(errors, "tenant code is required (-code)")
	}
	if tenantName == "" {
		errors = append(errors, "tenant name is required (-name)")
	}
	if subdomain == "" {
		errors = append(errors, "subdomain is required (-subdomain)")
	}
	if hospitalName == "" {
		errors = append(errors, "hospital name is required (-hospital-name)")
	}
	if hospitalCode == "" {
		errors = append(errors, "hospital code is required (-hospital-code)")
	} else if len(hospitalCode) != 8 {
		errors = append(errors, "hospital code must be exactly 8 characters")
	}
	if adminUsername == "" {
		errors = append(errors, "admin username is required (-admin-user)")
	} else if len(adminUsername) < 5 {
		errors = append(errors, "admin username must be at least 5 characters")
	}
	if adminEmail == "" {
		errors = append(errors, "admin email is required (-admin-email)")
	}

	if len(errors) > 0 {
		fmt.Println("Validation errors:")
		for _, err := range errors {
			fmt.Printf("  - %s\n", err)
		}
		os.Exit(1)
	}

	// Prompt for password if not provided
	if adminPassword == "" {
		fmt.Print("Enter admin password (min 6 chars): ")
		passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			log.Fatalf("Failed to read password: %v", err)
		}
		fmt.Println()
		adminPassword = string(passwordBytes)

		if len(adminPassword) < 6 {
			log.Fatal("Admin password must be at least 6 characters")
		}

		// Confirm password
		fmt.Print("Confirm admin password: ")
		confirmBytes, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			log.Fatalf("Failed to read password confirmation: %v", err)
		}
		fmt.Println()

		if adminPassword != string(confirmBytes) {
			log.Fatal("Passwords do not match")
		}
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

	// Initialize dependencies
	dbManager := database.NewTenantDBManager(db)
	tenantRepo := repository.NewTenantRepository(db)
	tenantService := services.NewTenantService(tenantRepo, dbManager, db)

	// Handle optional address
	var addressPtr *string
	if address != "" {
		addressPtr = &address
	}

	// Display summary
	fmt.Println("\n=== Tenant Creation Summary ===")
	fmt.Printf("Tenant Code:    %s\n", tenantCode)
	fmt.Printf("Tenant Name:    %s\n", tenantName)
	fmt.Printf("Subdomain:      %s\n", subdomain)
	fmt.Printf("Hospital Name:  %s\n", hospitalName)
	fmt.Printf("Hospital Code:  %s\n", hospitalCode)
	if address != "" {
		fmt.Printf("Address:        %s\n", address)
	}
	fmt.Printf("Admin Username: %s\n", adminUsername)
	fmt.Printf("Admin Email:    %s\n", adminEmail)
	fmt.Println("================================")

	// Confirm creation
	fmt.Print("\nProceed with tenant creation? (y/N): ")
	reader := bufio.NewReader(os.Stdin)
	confirm, _ := reader.ReadString('\n')
	confirm = strings.TrimSpace(strings.ToLower(confirm))

	if confirm != "y" && confirm != "yes" {
		fmt.Println("Tenant creation cancelled.")
		os.Exit(0)
	}

	// Create tenant with admin
	fmt.Println("\nCreating tenant...")
	tenant, err := tenantService.SetupTenantWithAdmin(
		tenantCode,
		tenantName,
		subdomain,
		hospitalName,
		hospitalCode,
		addressPtr,
		adminUsername,
		adminPassword,
		adminEmail,
	)
	if err != nil {
		log.Fatalf("Failed to create tenant: %v", err)
	}

	fmt.Println("\n✓ Tenant created successfully!")
	fmt.Println("\n=== Tenant Details ===")
	fmt.Printf("ID:           %d\n", tenant.ID)
	fmt.Printf("Tenant Code:  %s\n", tenant.TenantCode)
	fmt.Printf("Schema Name:  %s\n", tenant.SchemaName)
	fmt.Printf("Subdomain:    %s\n", tenant.Subdomain)
	fmt.Printf("Hospital:     %s (%s)\n", tenant.HospitalName, tenant.HospitalCode)
	fmt.Println("======================")
	fmt.Println("\nAdmin user has been created. You can now login with the provided credentials.")
}

func runList() {
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

	// Get all tenants directly from database
	var tenants []domain.Tenant
	if err := db.Find(&tenants).Error; err != nil {
		log.Fatalf("Failed to get tenants: %v", err)
	}

	if len(tenants) == 0 {
		fmt.Println("No tenants found.")
		return
	}

	fmt.Println("\n=== Registered Tenants ===")
	fmt.Printf("%-5s %-15s %-25s %-20s %-10s\n", "ID", "Code", "Name", "Schema", "Active")
	fmt.Println(strings.Repeat("-", 80))

	for _, t := range tenants {
		activeStatus := "Yes"
		if !t.IsActive {
			activeStatus = "No"
		}
		fmt.Printf("%-5d %-15s %-25s %-20s %-10s\n",
			t.ID,
			truncate(t.TenantCode, 15),
			truncate(t.Name, 25),
			truncate(t.SchemaName, 20),
			activeStatus,
		)
	}
	fmt.Println(strings.Repeat("-", 80))
	fmt.Printf("Total: %d tenant(s)\n", len(tenants))
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func printUsage() {
	fmt.Println(`
HIS Tenant Management Tool

Usage:
  go run cmd/tenant/main.go <command> [options]

Commands:
  create    Create a new tenant with admin user
  list      List all registered tenants

Create Options:
  -code           Tenant code (required)
  -name           Tenant name (required)
  -subdomain      Subdomain for the tenant (required)
  -hospital-name  Hospital name (required, max 150 chars)
  -hospital-code  Hospital code (required, exactly 8 chars)
  -address        Hospital address (optional)
  -admin-user     Admin username (required, min 5 chars)
  -admin-pass     Admin password (optional, will prompt if not provided)
  -admin-email    Admin email (required)

Examples:
  # Create a new tenant (will prompt for password)
  go run cmd/tenant/main.go create \
    -code=HOSP001 \
    -name="Bangkok Hospital" \
    -subdomain=bangkok \
    -hospital-name="Bangkok General Hospital" \
    -hospital-code=BKGH0001 \
    -admin-user=admin \
    -admin-email=admin@bangkok-hospital.com

  # Create with password inline (for scripts)
  go run cmd/tenant/main.go create \
    -code=HOSP001 \
    -name="Bangkok Hospital" \
    -subdomain=bangkok \
    -hospital-name="Bangkok General Hospital" \
    -hospital-code=BKGH0001 \
    -admin-user=admin \
    -admin-pass=secret123 \
    -admin-email=admin@bangkok-hospital.com

  # List all tenants
  go run cmd/tenant/main.go list
`)
}
