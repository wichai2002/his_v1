.PHONY: run build test clean migrate-up migrate-down migrate-status migrate-reset create-hospital

# Build the application
build:
	go build -o bin/api cmd/api/main.go
	go build -o bin/migrate cmd/migrate/main.go
	go build -o bin/create_hospital cmd/create_hospital/main.go

# Run the application
run:
	go run cmd/api/main.go

# Run tests
test:
	go test -v ./...

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

# Create a new hospital with admin user
# Usage: make create-hospital NAME="Hospital Name" CODE="HOS002" PHONE="0812345678" EMAIL="email@hospital.com" ADDRESS="Address"
create-hospital:
	@if [ -z "$(NAME)" ] || [ -z "$(CODE)" ] || [ -z "$(PHONE)" ] || [ -z "$(EMAIL)" ] || [ -z "$(ADDRESS)" ]; then \
		echo "Usage: make create-hospital NAME=\"Hospital Name\" CODE=\"HOS002\" PHONE=\"0812345678\" EMAIL=\"email@hospital.com\" ADDRESS=\"Address\""; \
		exit 1; \
	fi
	go run cmd/create_hospital/main.go -name "$(NAME)" -code "$(CODE)" -phone "$(PHONE)" -email "$(EMAIL)" -address "$(ADDRESS)"
