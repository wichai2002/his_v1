.PHONY: run build test test-unit test-cover test-verbose clean migrate-up migrate-down migrate-status migrate-reset tenant-create tenant-list

# Build the application
build:
	go build -o bin/api cmd/api/main.go
	go build -o bin/migrate cmd/migrate/main.go
	go build -o bin/tenant cmd/tenant/main.go

# Run the application
run:
	go run cmd/api/main.go

# Run all tests
test:
	go test ./tests/...

# Run unit tests with verbose output
test-verbose:
	go test -v ./tests/...

# Run tests with coverage report
test-cover:
	go test -cover ./tests/...

# Run tests with detailed coverage report
test-cover-html:
	go test -coverprofile=coverage.out ./tests/...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Clean build artifacts
clean:
	rm -rf bin/

# Download dependencies
deps:
	go mod tidy
	go mod download

# Migration commands
migrate-up:
	go run cmd/migrate/main.go up

migrate-down:
	go run cmd/migrate/main.go down

migrate-status:
	go run cmd/migrate/main.go status

migrate-reset:
	go run cmd/migrate/main.go reset

# Create database
create-db:
	createdb -U postgres his_db

# Drop database
drop-db:
	dropdb -U postgres his_db

# Docker compose up (if using docker)
docker-up:
	docker-compose up -d

# Docker compose down
docker-down:
	docker-compose down

# Tenant Management Commands
# Create a new tenant with admin user
# Usage: make tenant-create CODE=HOSP001 NAME="Hospital Name" SUBDOMAIN=hospital HOSPITAL_NAME="Hospital" HOSPITAL_CODE=HOSP0001 ADMIN_USER=admin ADMIN_EMAIL=admin@hospital.com
tenant-create:
	@if [ -z "$(CODE)" ] || [ -z "$(NAME)" ] || [ -z "$(SUBDOMAIN)" ] || [ -z "$(HOSPITAL_NAME)" ] || [ -z "$(HOSPITAL_CODE)" ] || [ -z "$(ADMIN_USER)" ] || [ -z "$(ADMIN_EMAIL)" ]; then \
		echo "Usage: make tenant-create CODE=HOSP001 NAME=\"Hospital Name\" SUBDOMAIN=hospital HOSPITAL_NAME=\"Hospital\" HOSPITAL_CODE=HOSP0001 ADMIN_USER=admin ADMIN_EMAIL=admin@hospital.com [ADDRESS=\"...\"] [ADMIN_PASS=\"...\"]"; \
		exit 1; \
	fi
	go run cmd/tenant/main.go create \
		-code="$(CODE)" \
		-name="$(NAME)" \
		-subdomain="$(SUBDOMAIN)" \
		-hospital-name="$(HOSPITAL_NAME)" \
		-hospital-code="$(HOSPITAL_CODE)" \
		-admin-user="$(ADMIN_USER)" \
		-admin-email="$(ADMIN_EMAIL)" \
		$(if $(ADDRESS),-address="$(ADDRESS)") \
		$(if $(ADMIN_PASS),-admin-pass="$(ADMIN_PASS)")

# List all tenants
tenant-list:
	go run cmd/tenant/main.go list
