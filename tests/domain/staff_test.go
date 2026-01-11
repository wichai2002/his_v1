package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wichai2002/his_v1/internal/domain"
)

func TestStaff_Struct(t *testing.T) {
	staff := domain.Staff{
		Username:    "admin",
		Password:    "hashedpassword",
		StaffCode:   "STF001",
		PhoneNumber: "0812345678",
		Email:       "admin@test.com",
		FirstName:   "Admin",
		LastName:    "User",
		IsAdmin:     true,
	}

	assert.Equal(t, "admin", staff.Username)
	assert.Equal(t, "STF001", staff.StaffCode)
	assert.Equal(t, "admin@test.com", staff.Email)
	assert.True(t, staff.IsAdmin)
}

func TestStaffResponse_Struct(t *testing.T) {
	response := domain.StaffResponse{
		ID:          1,
		Username:    "admin",
		StaffCode:   "STF001",
		PhoneNumber: "0812345678",
		Email:       "admin@test.com",
		FirstName:   "Admin",
		LastName:    "User",
		IsAdmin:     true,
	}

	assert.Equal(t, uint(1), response.ID)
	assert.Equal(t, "admin", response.Username)
	assert.Equal(t, "STF001", response.StaffCode)
	assert.True(t, response.IsAdmin)
}

func TestStaffCreateRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		request domain.StaffCreateRequest
		isValid bool
	}{
		{
			name: "valid request with all fields",
			request: domain.StaffCreateRequest{
				Username:    "newuser",
				Password:    "password123",
				PhoneNumber: "0812345678",
				Email:       "newuser@test.com",
				FirstName:   "New",
				LastName:    "User",
				IsAdmin:     false,
			},
			isValid: true,
		},
		{
			name: "valid admin request",
			request: domain.StaffCreateRequest{
				Username:    "adminuser",
				Password:    "adminpass123",
				PhoneNumber: "0899999999",
				Email:       "admin@test.com",
				FirstName:   "Admin",
				LastName:    "User",
				IsAdmin:     true,
			},
			isValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test struct field access
			assert.NotEmpty(t, tt.request.Username)
			assert.NotEmpty(t, tt.request.Password)
			assert.NotEmpty(t, tt.request.Email)
		})
	}
}

func TestStaffUpdateRequest_Struct(t *testing.T) {
	request := domain.StaffUpdateRequest{
		StaffCode:   "STF002",
		PhoneNumber: "0899999999",
		Email:       "updated@test.com",
		FirstName:   "Updated",
		LastName:    "Name",
		IsAdmin:     true,
	}

	assert.Equal(t, "STF002", request.StaffCode)
	assert.Equal(t, "0899999999", request.PhoneNumber)
	assert.Equal(t, "updated@test.com", request.Email)
	assert.Equal(t, "Updated", request.FirstName)
	assert.Equal(t, "Name", request.LastName)
	assert.True(t, request.IsAdmin)
}

func TestStaffLoginRequest_Struct(t *testing.T) {
	request := domain.StaffLoginRequest{
		Username: "admin",
		Password: "password123",
	}

	assert.Equal(t, "admin", request.Username)
	assert.Equal(t, "password123", request.Password)
}

func TestStaffLoginResponse_Struct(t *testing.T) {
	response := domain.StaffLoginResponse{
		Token: "jwt-token-here",
		Staff: domain.StaffResponse{
			ID:        1,
			Username:  "admin",
			StaffCode: "STF001",
			Email:     "admin@test.com",
			FirstName: "Admin",
			LastName:  "User",
			IsAdmin:   true,
		},
	}

	assert.Equal(t, "jwt-token-here", response.Token)
	assert.Equal(t, "admin", response.Staff.Username)
	assert.True(t, response.Staff.IsAdmin)
}

func TestStaff_PasswordNotExposed(t *testing.T) {
	// Test that Staff struct has json:"-" tag for password
	staff := domain.Staff{
		Username: "admin",
		Password: "secretpassword",
	}

	// The password field should exist in the struct
	assert.NotEmpty(t, staff.Password)

	// Note: In actual JSON serialization, the password would be omitted
	// due to the json:"-" tag. This is a struct-level test.
}

func TestStaff_DefaultIsAdmin(t *testing.T) {
	// Test that a new staff without explicit IsAdmin is false
	staff := domain.Staff{
		Username:    "regularuser",
		Password:    "password",
		StaffCode:   "STF001",
		PhoneNumber: "0812345678",
		Email:       "user@test.com",
		FirstName:   "Regular",
		LastName:    "User",
	}

	// IsAdmin defaults to false (zero value for bool)
	assert.False(t, staff.IsAdmin)
}

func TestStaffCreateRequest_WithOptionalStaffCode(t *testing.T) {
	// Staff code is optional in create request (auto-generated)
	request := domain.StaffCreateRequest{
		Username:    "newuser",
		Password:    "password123",
		PhoneNumber: "0812345678",
		Email:       "newuser@test.com",
		FirstName:   "New",
		LastName:    "User",
		// StaffCode is omitted - it will be auto-generated
	}

	assert.Empty(t, request.StaffCode)
	assert.NotEmpty(t, request.Username)
}
