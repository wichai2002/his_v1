# Hospital Information System (HIS) API

A RESTful API built with Go, Gin, GORM, and PostgreSQL following Clean Architecture principles.

## Features

- **Clean Architecture**: Separation of concerns with Domain, Repository, Usecase, and Delivery layers
- **JWT Authentication**: Secure authentication with role-based access control
- **PostgreSQL Database**: Robust data persistence with GORM ORM
- **Version-controlled Migrations**: GORM-based migration system with version tracking
- **RESTful API**: Well-structured API endpoints

## Project Structure

```
.
├── cmd/
│   ├── api/
│   │   └── main.go              # API server entry point
│   └── migrate/
│       └── main.go              # Migration CLI tool
├── config/
│   └── config.go                # Configuration management
├── internal/
│   ├── domain/                  # Business entities and interfaces
│   │   ├── hospital.go
│   │   ├── staff.go
│   │   └── patient.go
│   ├── repository/              # Data access layer
│   │   ├── hospital_repository.go
│   │   ├── staff_repository.go
│   │   └── patient_repository.go
│   ├── usecase/                 # Business logic layer
│   │   ├── hospital_usecase.go
│   │   ├── staff_usecase.go
│   │   └── patient_usecase.go
│   ├── delivery/
│   │   └── http/
│   │       ├── handler/         # HTTP handlers
│   │       ├── middleware/      # Auth middleware
│   │       └── router.go        # Route definitions
│   └── infrastructure/
│       └── database/
│           ├── postgres.go      # Database connection
│           ├── migration.go     # Migration engine
│           └── migrations.go    # Migration definitions
├── pkg/
│   ├── jwt/                     # JWT utilities
│   └── utils/                   # Response helpers
├── env.example
├── go.mod
├── Makefile
└── README.md
```

## Prerequisites

- Go 1.21 or higher
- PostgreSQL 14 or higher

## Getting Started

### 1. Clone the repository

```bash
cd HIS_v1
```

### 2. Set up environment variables

```bash
cp env.example .env
# Edit .env with your configuration
```

### 3. Create PostgreSQL database

```bash
createdb his_db
# or using make
make create-db
```

### 4. Install dependencies

```bash
go mod tidy
# or using make
make deps
```

### 5. Run migrations

```bash
go run cmd/migrate/main.go up
# or using make
make migrate-up
```

### 6. Run the application

```bash
go run cmd/api/main.go
# or using make
make run
```

The server will start on `http://localhost:8080`

## Database Migrations

This project uses a **version-controlled GORM migration system**. Each migration has a unique version and is tracked in a `migrations` table.

### Migration Commands

```bash
# Run all pending migrations
make migrate-up
# or: go run cmd/migrate/main.go up

# Rollback the last migration
make migrate-down
# or: go run cmd/migrate/main.go down

# Show migration status
make migrate-status
# or: go run cmd/migrate/main.go status

# Rollback all migrations
make migrate-reset
# or: go run cmd/migrate/main.go reset
```

### Migration Status Example

```
=== Migration Status ===
VERSION              NAME                                     STATUS     APPLIED AT
--------------------------------------------------------------------------------
20240101_001         create_hospitals_table                   Applied    2024-01-15 10:30:00
20240101_002         create_staffs_table                      Applied    2024-01-15 10:30:00
20240101_003         create_patients_table                    Applied    2024-01-15 10:30:00
20240101_004         seed_default_data                        Applied    2024-01-15 10:30:00
```

### Creating New Migrations

Add new migrations in `internal/infrastructure/database/migrations.go`:

```go
// 1. Create the migration function
func migration_20240615_005_add_department_table() MigrationDefinition {
    return MigrationDefinition{
        Version: "20240615_005",
        Name:    "add_department_table",
        Up: func(db *gorm.DB) error {
            type Department struct {
                ID   uint   `gorm:"primaryKey"`
                Name string `gorm:"not null"`
            }
            return db.AutoMigrate(&Department{})
        },
        Down: func(db *gorm.DB) error {
            return db.Migrator().DropTable("departments")
        },
    }
}

// 2. Add it to GetMigrations() slice
func GetMigrations() []MigrationDefinition {
    return []MigrationDefinition{
        migration_20240101_001_create_hospitals_table(),
        migration_20240101_002_create_staffs_table(),
        migration_20240101_003_create_patients_table(),
        migration_20240101_004_seed_default_data(),
        migration_20240615_005_add_department_table(), // Add here
    }
}
```

### Version Naming Convention

```
YYYYMMDD_XXX_description
│        │   │
│        │   └── Brief description (snake_case)
│        └────── Sequence number (001, 002, ...)
└─────────────── Date (Year/Month/Day)
```

## API Endpoints

### Health Check
- `GET /health` - Check API status

### Staff APIs
| Method | Endpoint | Description | Auth Required | Admin Only |
|--------|----------|-------------|---------------|------------|
| POST | `/staff/login` | Staff login | No | No |
| POST | `/staff/logout` | Staff logout | Yes | No |
| GET | `/staff/` | Get all staff | Yes | No |
| GET | `/staff/:id` | Get staff by ID | Yes | No |
| POST | `/staff/create` | Create new staff | Yes | Yes |
| PUT | `/staff/update/:id` | Update staff | Yes | No |
| DELETE | `/staff/delete/:id` | Delete staff | Yes | Yes |

### Hospital APIs
| Method | Endpoint | Description | Auth Required | Admin Only |
|--------|----------|-------------|---------------|------------|
| GET | `/hospital/` | Get all hospitals | Yes | No |
| GET | `/hospital/:id` | Get hospital by ID | Yes | No |
| POST | `/hospital/create` | Create new hospital | Yes | Yes |
| PUT | `/hospital/update/:id` | Update hospital | Yes | Yes |
| DELETE | `/hospital/delete/:id` | Delete hospital | Yes | Yes |

### Patient APIs
| Method | Endpoint | Description | Auth Required | Admin Only |
|--------|----------|-------------|---------------|------------|
| GET | `/patient/search/:id` | Search patient by ID | Yes | No |
| POST | `/patient/create` | Create new patient | Yes | No |
| PUT | `/patient/update/:id` | Update patient | Yes | No |
| DELETE | `/patient/delete/:id` | Delete patient | Yes | No |

## Authentication

### Login

```bash
curl -X POST http://localhost:8080/staff/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "admin123"}'
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
      ...
    }
  }
}
```

### Using the Token

Include the JWT token in the Authorization header:

```bash
curl -X GET http://localhost:8080/staff/ \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..."
```

## Default Admin Credentials

- **Username**: admin
- **Password**: admin123

⚠️ **Important**: Change the default credentials in production!

## Models

### Hospital
| Field | Type | Constraints |
|-------|------|-------------|
| name | string | unique, not null |
| hospital_code | string | unique, not null |
| phone_number | string | unique, not null |
| email | string | - |
| address | string | - |
| hn_running_number | int | default: 0 |

### Staff
| Field | Type | Constraints |
|-------|------|-------------|
| username | string | unique, not null |
| password | string | hashed, not null |
| staff_code | string | - |
| phone_number | string | unique, not null |
| email | string | unique, not null |
| first_name | string | - |
| last_name | string | - |
| hospital_id | uint | foreign key |
| is_admin | bool | default: false |

### Patient
| Field | Type | Constraints |
|-------|------|-------------|
| first_name_th | string | - |
| last_name_th | string | - |
| middle_name_th | string | - |
| first_name_en | string | - |
| last_name_en | string | - |
| middle_name_en | string | - |
| date_of_birth | time.Time | - |
| nick_name_th | string | - |
| nick_name_en | string | - |
| patient_hn | string | unique, not null, auto-generated |
| national_id | string | unique |
| passport_id | string | unique |
| phone_number | string | unique |
| email | string | unique |
| gender | string | - |
| nationality | string | - |
| blood_grp | string | - |
| hospital_id | uint | foreign key |

## License

MIT License
