# Hospital Information System (HIS) API

A multi-tenant RESTful API built with Go, Gin, GORM, and PostgreSQL following Clean Architecture principles. Each hospital (tenant) has complete data isolation through PostgreSQL schemas.

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────────────┐
│                              NGINX                                       │
│                    (Rate Limiting, SSL, Load Balancing)                 │
└─────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                           Go API Server                                  │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐         │
│  │  Tenant         │  │  Auth           │  │  Handlers       │         │
│  │  Middleware     │──│  Middleware     │──│  (Staff/Patient)│         │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘         │
└─────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                         PostgreSQL Database                              │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐                  │
│  │public schema │  │tenant_hosp01 │  │tenant_hosp02 │  ...             │
│  │  - tenants   │  │  - staffs    │  │  - staffs    │                  │
│  │  - migrations│  │  - patients  │  │  - patients  │                  │
│  └──────────────┘  └──────────────┘  └──────────────┘                  │
└─────────────────────────────────────────────────────────────────────────┘
```

## Features

- **Multi-Tenant Architecture**: Schema-based tenant isolation for complete data separation
- **Subdomain Routing**: Access tenant data via `{subdomain}.yourdomain.com`
- **Clean Architecture**: Separation of concerns with Domain, Repository, Service, and Delivery layers
- **JWT Authentication**: Secure authentication with role-based access control (Admin/Staff)
- **PostgreSQL Database**: Robust data persistence with GORM ORM
- **Version-controlled Migrations**: GORM-based migration system with version tracking
- **Docker Ready**: Production-ready Docker and Docker Compose configuration
- **Rate Limiting**: NGINX-based rate limiting to prevent abuse
- **Auto HN Generation**: Automatic Hospital Number generation per tenant

## ER Diagram

![ER Diagram](docs/er-diagram.svg)

## Project Structure

```
.
├── cmd/
│   ├── api/
│   │   └── main.go              # API server entry point
│   ├── migrate/
│   │   └── main.go              # Migration CLI tool
│   └── tenant/
│       └── main.go              # Tenant management CLI
├── config/
│   └── config.go                # Configuration management
├── internal/
│   ├── domain/                  # Business entities and interfaces
│   │   ├── errors.go
│   │   ├── tenant.go            # Tenant entity
│   │   ├── staff.go             # Staff entity
│   │   └── patient.go           # Patient entity
│   ├── repository/              # Data access layer
│   │   ├── tenant_repository.go
│   │   ├── staff_repository.go
│   │   └── patient_repository.go
│   ├── services/                # Business logic layer
│   │   ├── tenant_service.go
│   │   ├── staff_service.go
│   │   └── patient_service.go
│   ├── delivery/
│   │   └── http/
│   │       ├── handler/         # HTTP handlers
│   │       ├── middleware/      # Auth & Tenant middleware
│   │       ├── routes/          # Route definitions
│   │       └── router.go        # Router setup
│   └── infrastructure/
│       └── database/
│           ├── postgres.go      # Database connection
│           ├── migrator.go      # Migration engine
│           ├── tenant_context.go# Tenant DB manager
│           └── migrations/      # Migration definitions
├── pkg/
│   ├── jwt/                     # JWT utilities
│   └── utils/                   # Response helpers
├── nginx/
│   ├── nginx.conf               # NGINX configuration
│   └── ssl/                     # SSL certificates
├── postman/                     # Postman collection & environment
├── tests/                       # Test files
├── docker-compose.yml
├── Dockerfile
├── Makefile
└── README.md
```

## Prerequisites

- Go 1.23 or higher
- PostgreSQL 16 or higher
- Docker & Docker Compose (optional)

## Quick Start with Docker

### 1. Clone and configure

```bash
cd HIS_v1

# Create environment file (optional, has defaults)
cat > .env << EOF
DB_USER=postgres
DB_PASSWORD=your_secure_password
DB_NAME=his_db
JWT_SECRET_KEY=your-super-secret-key-change-in-production
EOF
```

### 2. Start services

```bash
# Build and start all services
docker compose up -d --build

# View logs
docker compose logs -f api
```

### 3. Create your first tenant

```bash
# Using docker exec
docker exec -it his_api ./tenant create \
  -code=HOSP001 \
  -name="Bangkok Hospital" \
  -subdomain=bangkok \
  -hospital-name="Bangkok General Hospital" \
  -hospital-code=BKGH0001 \
  -admin-user=admin \
  -admin-email=admin@bangkok-hospital.com
```

### 4. Access the API

```bash
# Health check
curl http://localhost/health

# Login (use subdomain for tenant isolation)
curl -X POST http://bangkok.localhost/api/v1/staff/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "your_password"}'
```

## Local Development Setup

### 1. Set up environment

```bash
# Copy example environment
cp .env.example .env

# Install dependencies
go mod tidy
```

### 2. Start PostgreSQL

```bash
# Using docker
docker run -d --name his_postgres \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=his_db \
  -p 5432:5432 \
  postgres:16-alpine

# Or create database manually
createdb -U postgres his_db
```

### 3. Run migrations

```bash
make migrate-up
```

### 4. Create a tenant

```bash
make tenant-create \
  CODE=HOSP001 \
  NAME="Test Hospital" \
  SUBDOMAIN=test \
  HOSPITAL_NAME="Test General Hospital" \
  HOSPITAL_CODE=TEST0001 \
  ADMIN_USER=admin \
  ADMIN_EMAIL=admin@test.com
```

### 5. Run the server

```bash
make run
```

## Multi-Tenant Usage

### Subdomain-based Access

Each tenant is accessed via their subdomain:

```bash
# Tenant 1: Bangkok Hospital
curl http://bangkok.localhost/api/v1/staff/login

# Tenant 2: Chiang Mai Hospital  
curl http://chiangmai.localhost/api/v1/staff/login
```

### Local Development Subdomains

Add to `/etc/hosts` for local testing:

```
127.0.0.1 bangkok.localhost
127.0.0.1 chiangmai.localhost
```

## Database Migrations

### Migration Commands

```bash
# Run all pending migrations
make migrate-up

# Rollback the last migration
make migrate-down

# Show migration status
make migrate-status

# Rollback all migrations
make migrate-reset
```

### Creating New Migrations

Add new migrations in `internal/infrastructure/database/migrations/`:

```go
func Migration_YYYYMMDD_XXX_description() MigrationDefinition {
    return MigrationDefinition{
        Version: "YYYYMMDD_XXX",
        Name:    "description",
        Up: func(db *gorm.DB) error {
            // Migration logic
            return db.AutoMigrate(&YourModel{})
        },
        Down: func(db *gorm.DB) error {
            return db.Migrator().DropTable("your_table")
        },
    }
}
```

## Tenant Management

### Create a Tenant

```bash
# Interactive (prompts for password)
make tenant-create \
  CODE=HOSP001 \
  NAME="Hospital Name" \
  SUBDOMAIN=hospital \
  HOSPITAL_NAME="Full Hospital Name" \
  HOSPITAL_CODE=HOSP0001 \
  ADMIN_USER=admin \
  ADMIN_EMAIL=admin@hospital.com

# Non-interactive (for scripts)
make tenant-create \
  CODE=HOSP001 \
  NAME="Hospital Name" \
  SUBDOMAIN=hospital \
  HOSPITAL_NAME="Full Hospital Name" \
  HOSPITAL_CODE=HOSP0001 \
  ADMIN_USER=admin \
  ADMIN_EMAIL=admin@hospital.com \
  ADMIN_PASS=secretpassword
```

### List All Tenants

```bash
make tenant-list
```

## API Endpoints

### Health Check

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Check API status |

### Staff APIs

| Method | Endpoint | Description | Auth | Admin |
|--------|----------|-------------|------|-------|
| POST | `/api/v1/staff/login` | Staff login | ❌ | ❌ |
| GET | `/api/v1/staff/` | Get all staff | ✅ | ❌ |
| GET | `/api/v1/staff/:id` | Get staff by ID | ✅ | ❌ |
| POST | `/api/v1/staff/create` | Create new staff | ✅ | ✅ |
| PUT | `/api/v1/staff/update/:id` | Update staff | ✅ | ❌ |
| DELETE | `/api/v1/staff/delete/:id` | Delete staff | ✅ | ✅ |

### Patient APIs

| Method | Endpoint | Description | Auth | Admin |
|--------|----------|-------------|------|-------|
| GET | `/api/v1/patient/search` | Search patients | ✅ | ❌ |
| GET | `/api/v1/patient/search/:id` | Get patient by ID | ✅ | ❌ |
| POST | `/api/v1/patient/create` | Create patient | ✅ | ❌ |
| PUT | `/api/v1/patient/update/:id` | Full update | ✅ | ❌ |
| PATCH | `/api/v1/patient/update/:id` | Partial update | ✅ | ❌ |
| DELETE | `/api/v1/patient/delete/:id` | Delete patient | ✅ | ❌ |

## Authentication

### Login

```bash
curl -X POST http://bangkok.localhost/api/v1/staff/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin123"
  }'
```

Response:
```json
{
  "success": true,
  "message": "login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "staff": {
      "id": 1,
      "username": "admin",
      "staff_code": "STF00001",
      "email": "admin@hospital.com",
      "is_admin": true
    }
  }
}
```

### Using the Token

```bash
curl -X GET http://bangkok.localhost/api/v1/staff/ \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..."
```

## Data Models

### Tenant (Public Schema)

| Field | Type | Description |
|-------|------|-------------|
| id | uint | Primary key |
| tenant_code | string | Unique tenant code |
| name | string | Tenant display name |
| schema_name | string | PostgreSQL schema name |
| subdomain | string | URL subdomain |
| is_active | bool | Tenant active status |
| hospital_name | string | Hospital full name |
| hospital_code | string | 8-char hospital code |
| address | string | Hospital address |
| hn_running | uint64 | HN sequence counter |

### Staff (Tenant Schema)

| Field | Type | Description |
|-------|------|-------------|
| id | uint | Primary key |
| username | string | Unique login username |
| password | string | Hashed password |
| staff_code | string | Auto-generated staff code |
| phone_number | string | Contact phone |
| email | string | Email address |
| first_name | string | First name |
| last_name | string | Last name |
| is_admin | bool | Admin privileges |

### Patient (Tenant Schema)

| Field | Type | Description |
|-------|------|-------------|
| id | uint | Primary key |
| patient_hn | string | Auto-generated HN |
| first_name_th | string | Thai first name |
| last_name_th | string | Thai last name |
| first_name_en | string | English first name |
| last_name_en | string | English last name |
| date_of_birth | date | Birth date |
| national_id | string | National ID (13 digits) |
| passport_id | string | Passport number |
| phone_number | string | Contact phone |
| email | string | Email address |
| gender | enum | M, F, OTHER |
| nationality | string | Nationality |
| blood_grp | enum | A, B, O, AB |

## Docker Commands

```bash
# Start all services
docker compose up -d

# Stop all services
docker compose down

# View logs
docker compose logs -f

# Rebuild and restart
docker compose up -d --build

# Access API container
docker exec -it his_api sh

# Access PostgreSQL
docker exec -it his_postgres psql -U postgres -d his_db

# Run migrations in container
docker exec -it his_api ./migrate up

# Create tenant in container
docker exec -it his_api ./tenant create ...
```

## Testing

```bash
# Run all tests
make test

# Run with verbose output
make test-verbose

# Run with coverage
make test-cover

# Generate HTML coverage report
make test-cover-html
```

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| SERVER_PORT | 8080 | API server port |
| DB_HOST | localhost | PostgreSQL host |
| DB_PORT | 5432 | PostgreSQL port |
| DB_USER | postgres | Database user |
| DB_PASSWORD | postgres | Database password |
| DB_NAME | his_db | Database name |
| DB_SSLMODE | disable | SSL mode |
| JWT_SECRET_KEY | - | JWT signing key |
| JWT_EXPIRES_IN_HOURS | 24 | Token expiry |

## Security Features

- **Schema Isolation**: Each tenant's data is in a separate PostgreSQL schema
- **JWT Authentication**: Stateless authentication with configurable expiry
- **Password Hashing**: bcrypt hashing for all passwords
- **Rate Limiting**: NGINX rate limits (30 req/s API, 5 req/min login)
- **Security Headers**: X-Frame-Options, X-Content-Type-Options, etc.
- **Non-root Container**: API runs as non-privileged user

## License

MIT License
