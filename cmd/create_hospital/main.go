package main

import (
	"crypto/rand"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/wichai2002/his_v1/config"
	"github.com/wichai2002/his_v1/internal/domain"
	"github.com/wichai2002/his_v1/internal/infrastructure/database"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func main() {
	// Define command line flags
	name := flag.String("name", "", "Hospital name (required)")
	code := flag.String("code", "", "Hospital code (required, e.g., HOS002)")
	phone := flag.String("phone", "", "Hospital phone number (required)")
	email := flag.String("email", "", "Hospital email (required)")
	address := flag.String("address", "", "Hospital address (required)")

	// Admin user flags (optional - will generate defaults if not provided)
	adminUsername := flag.String("admin-username", "", "Admin username (optional, will be auto-generated)")
	adminPassword := flag.String("admin-password", "", "Admin password (optional, will be auto-generated)")
	adminEmail := flag.String("admin-email", "", "Admin email (optional, will use hospital email)")
	adminFirstName := flag.String("admin-firstname", "Admin", "Admin first name")
	adminLastName := flag.String("admin-lastname", "", "Admin last name (optional, will use hospital name)")

	flag.Parse()

	// Validate required fields
	if *name == "" || *code == "" || *phone == "" || *email == "" || *address == "" {
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

	// Create hospital and admin
	result, err := createHospitalWithAdmin(db, HospitalInput{
		Name:           *name,
		Code:           *code,
		Phone:          *phone,
		Email:          *email,
		Address:        *address,
		AdminUsername:  *adminUsername,
		AdminPassword:  *adminPassword,
		AdminEmail:     *adminEmail,
		AdminFirstName: *adminFirstName,
		AdminLastName:  *adminLastName,
	})

	if err != nil {
		log.Fatalf("Failed to create hospital: %v", err)
	}

	// Display results
	printResult(result)
}

type HospitalInput struct {
	Name           string
	Code           string
	Phone          string
	Email          string
	Address        string
	AdminUsername  string
	AdminPassword  string
	AdminEmail     string
	AdminFirstName string
	AdminLastName  string
}

type CreateResult struct {
	Hospital      *domain.Hospital
	AdminUsername string
	AdminPassword string
	AdminEmail    string
}

func createHospitalWithAdmin(db *gorm.DB, input HospitalInput) (*CreateResult, error) {
	// Generate defaults for admin if not provided
	adminUsername := input.AdminUsername
	if adminUsername == "" {
		// Generate username from hospital code, e.g., HOS002 -> admin_hos002
		adminUsername = "admin_" + strings.ToLower(input.Code)
	}

	adminPassword := input.AdminPassword
	if adminPassword == "" {
		// Generate a random password
		adminPassword = generateRandomPassword(12)
	}

	adminEmail := input.AdminEmail
	if adminEmail == "" {
		adminEmail = "admin_" + strings.ToLower(input.Code) + "@his.com"
	}

	adminLastName := input.AdminLastName
	if adminLastName == "" {
		adminLastName = input.Name
	}

	// Start transaction
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create hospital
	hospital := &domain.Hospital{
		Name:            input.Name,
		HospitalCode:    input.Code,
		PhoneNumber:     input.Phone,
		Email:           input.Email,
		Address:         input.Address,
		HNRunningNumber: 0,
	}

	if err := tx.Create(hospital).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create hospital: %w", err)
	}

	// Hash admin password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Generate staff code for admin
	staffCode := "ADM" + input.Code[len(input.Code)-3:]

	// Create admin user
	admin := &domain.Staff{
		Username:    adminUsername,
		Password:    string(hashedPassword),
		StaffCode:   staffCode,
		PhoneNumber: input.Phone,
		Email:       adminEmail,
		FirstName:   input.AdminFirstName,
		LastName:    adminLastName,
		HospitalID:  hospital.ID,
		IsAdmin:     true,
	}

	if err := tx.Create(admin).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create admin user: %w", err)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &CreateResult{
		Hospital:      hospital,
		AdminUsername: adminUsername,
		AdminPassword: adminPassword,
		AdminEmail:    adminEmail,
	}, nil
}

func generateRandomPassword(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback to a default if random generation fails
		return "Admin@" + hex.EncodeToString(bytes)[:8]
	}
	// Create a password with alphanumeric characters
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$"
	password := make([]byte, length)
	for i := range password {
		password[i] = charset[int(bytes[i])%len(charset)]
	}
	return string(password)
}

func printResult(result *CreateResult) {
	fmt.Println()
	fmt.Println("╔══════════════════════════════════════════════════════════════════╗")
	fmt.Println("║              HOSPITAL CREATED SUCCESSFULLY                       ║")
	fmt.Println("╠══════════════════════════════════════════════════════════════════╣")
	fmt.Println("║ HOSPITAL DETAILS                                                 ║")
	fmt.Println("╠══════════════════════════════════════════════════════════════════╣")
	fmt.Printf("║  ID:           %-50d ║\n", result.Hospital.ID)
	fmt.Printf("║  Name:         %-50s ║\n", truncateString(result.Hospital.Name, 50))
	fmt.Printf("║  Code:         %-50s ║\n", result.Hospital.HospitalCode)
	fmt.Printf("║  Phone:        %-50s ║\n", result.Hospital.PhoneNumber)
	fmt.Printf("║  Email:        %-50s ║\n", truncateString(result.Hospital.Email, 50))
	fmt.Printf("║  Address:      %-50s ║\n", truncateString(result.Hospital.Address, 50))
	fmt.Println("╠══════════════════════════════════════════════════════════════════╣")
	fmt.Println("║ ADMIN CREDENTIALS                                                ║")
	fmt.Println("╠══════════════════════════════════════════════════════════════════╣")
	fmt.Printf("║  Username:     %-50s ║\n", result.AdminUsername)
	fmt.Printf("║  Password:     %-50s ║\n", result.AdminPassword)
	fmt.Printf("║  Email:        %-50s ║\n", truncateString(result.AdminEmail, 50))
	fmt.Println("╠══════════════════════════════════════════════════════════════════╣")
	fmt.Println("║  ⚠️  IMPORTANT: Save these credentials securely!                  ║")
	fmt.Println("║  The password cannot be retrieved after this.                    ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════════╝")
	fmt.Println()
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func printUsage() {
	fmt.Println(`
		HIS Create Hospital Tool

		Usage:
		go run cmd/create_hospital/main.go [flags]

		Required Flags:
		-name string      Hospital name
		-code string      Hospital code (e.g., HOS002)
		-phone string     Hospital phone number
		-email string     Hospital email
		-address string   Hospital address

		Optional Flags (Admin User):
		-admin-username string    Admin username (default: auto-generated from hospital code)
		-admin-password string    Admin password (default: auto-generated random password)
		-admin-email string       Admin email (default: admin_<code>@his.com)
		-admin-firstname string   Admin first name (default: "Admin")
		-admin-lastname string    Admin last name (default: hospital name)

		Examples:
		# Create hospital with auto-generated admin credentials
		go run cmd/create_hospital/main.go \
			-name "Bangkok General Hospital" \
			-code "HOS002" \
			-phone "0812345678" \
			-email "contact@bgh.com" \
			-address "123 Sukhumvit Road, Bangkok"

		# Create hospital with custom admin credentials
		go run cmd/create_hospital/main.go \
			-name "Chiang Mai Hospital" \
			-code "HOS003" \
			-phone "0823456789" \
			-email "contact@cmh.com" \
			-address "456 Nimman Road, Chiang Mai" \
			-admin-username "cmh_admin" \
			-admin-password "SecurePass123"
	`)
}
