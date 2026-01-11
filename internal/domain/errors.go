package domain

import "errors"

// Domain errors for consistent error handling across the application
var (
	// ErrNotFound is returned when a requested resource does not exist
	ErrNotFound = errors.New("record not found")

	// ErrDuplicateEntry is returned when attempting to create a duplicate record
	ErrDuplicateEntry = errors.New("duplicate entry")

	// ErrInvalidInput is returned when input validation fails
	ErrInvalidInput = errors.New("invalid input")

	// ErrInvalidDateFormat is returned when date parsing fails
	ErrInvalidDateFormat = errors.New("invalid date format")

	// ErrFutureDateOfBirth is returned when date of birth is in the future
	ErrFutureDateOfBirth = errors.New("date of birth cannot be in the future")

	// ErrDateOfBirthTooOld is returned when date of birth is too far in the past
	ErrDateOfBirthTooOld = errors.New("date of birth is too old")

	// ErrTenantRequired is returned when tenant context is required but not provided
	ErrTenantRequired = errors.New("tenant context required")

	// ErrInvalidTenant is returned when tenant is not valid or inactive
	ErrInvalidTenant = errors.New("invalid or inactive tenant")

	// ErrInvalidSchemaName is returned when schema name contains invalid characters
	ErrInvalidSchemaName = errors.New("invalid schema name")

	// ErrDatabaseConnection is returned when database connection fails
	ErrDatabaseConnection = errors.New("database connection error")
)

// IsNotFoundError checks if the error is a not found error
func IsNotFoundError(err error) bool {
	return errors.Is(err, ErrNotFound)
}

// IsDuplicateError checks if the error is a duplicate entry error
func IsDuplicateError(err error) bool {
	return errors.Is(err, ErrDuplicateEntry)
}
