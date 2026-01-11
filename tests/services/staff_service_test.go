package services_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/wichai2002/his_v1/internal/domain"
	"github.com/wichai2002/his_v1/internal/mocks"
	"github.com/wichai2002/his_v1/internal/services"
	"golang.org/x/crypto/bcrypt"
)

func TestStaffService_Login(t *testing.T) {
	// Create a hashed password for testing
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	tests := []struct {
		name        string
		request     *domain.StaffLoginRequest
		schemaName  string
		mockStaff   *domain.Staff
		mockError   error
		mockToken   string
		tokenError  error
		expectError bool
		errorMsg    string
	}{
		{
			name: "successful login",
			request: &domain.StaffLoginRequest{
				Username: "admin",
				Password: "password123",
			},
			schemaName: "tenant_test",
			mockStaff: &domain.Staff{
				Username:  "admin",
				Password:  string(hashedPassword),
				StaffCode: "STF001",
				Email:     "admin@test.com",
				FirstName: "Admin",
				LastName:  "User",
				IsAdmin:   true,
			},
			mockError:   nil,
			mockToken:   "jwt-token-here",
			tokenError:  nil,
			expectError: false,
		},
		{
			name: "user not found",
			request: &domain.StaffLoginRequest{
				Username: "nonexistent",
				Password: "password123",
			},
			schemaName:  "tenant_test",
			mockStaff:   nil,
			mockError:   errors.New("record not found"),
			expectError: true,
			errorMsg:    "invalid credentials",
		},
		{
			name: "wrong password",
			request: &domain.StaffLoginRequest{
				Username: "admin",
				Password: "wrongpassword",
			},
			schemaName: "tenant_test",
			mockStaff: &domain.Staff{
				Username:  "admin",
				Password:  string(hashedPassword),
				StaffCode: "STF001",
			},
			mockError:   nil,
			expectError: true,
			errorMsg:    "invalid credentials",
		},
		{
			name: "token generation error",
			request: &domain.StaffLoginRequest{
				Username: "admin",
				Password: "password123",
			},
			schemaName: "tenant_test",
			mockStaff: &domain.Staff{
				Username:  "admin",
				Password:  string(hashedPassword),
				StaffCode: "STF001",
			},
			mockError:   nil,
			mockToken:   "",
			tokenError:  errors.New("token generation failed"),
			expectError: true,
			errorMsg:    "token generation failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewMockStaffRepository()
			mockJWT := mocks.NewMockJWTService()

			mockRepo.On("GetByUsername", tt.request.Username, tt.schemaName).Return(tt.mockStaff, tt.mockError)

			if tt.mockError == nil && tt.mockStaff != nil {
				// Only expect token generation if password verification will pass
				err := bcrypt.CompareHashAndPassword([]byte(tt.mockStaff.Password), []byte(tt.request.Password))
				if err == nil {
					mockJWT.On("GenerateToken", tt.mockStaff.ID, tt.mockStaff.Username, tt.mockStaff.IsAdmin, tt.schemaName).Return(tt.mockToken, tt.tokenError)
				}
			}

			service := services.NewStaffService(mockRepo, mockJWT)
			result, err := service.Login(tt.request, tt.schemaName)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.mockToken, result.Token)
				assert.Equal(t, tt.mockStaff.Username, result.Staff.Username)
			}

			mockRepo.AssertExpectations(t)
			mockJWT.AssertExpectations(t)
		})
	}
}

func TestStaffService_GetAll(t *testing.T) {
	tests := []struct {
		name          string
		schemaName    string
		mockStaffs    []domain.Staff
		mockError     error
		expectedCount int
		expectError   bool
	}{
		{
			name:       "successful get all",
			schemaName: "tenant_test",
			mockStaffs: []domain.Staff{
				{Username: "admin", StaffCode: "STF001"},
				{Username: "user1", StaffCode: "STF002"},
			},
			mockError:     nil,
			expectedCount: 2,
			expectError:   false,
		},
		{
			name:          "empty list",
			schemaName:    "tenant_test",
			mockStaffs:    []domain.Staff{},
			mockError:     nil,
			expectedCount: 0,
			expectError:   false,
		},
		{
			name:          "database error",
			schemaName:    "tenant_test",
			mockStaffs:    nil,
			mockError:     errors.New("database error"),
			expectedCount: 0,
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewMockStaffRepository()
			mockJWT := mocks.NewMockJWTService()

			mockRepo.On("GetAll", tt.schemaName).Return(tt.mockStaffs, tt.mockError)

			service := services.NewStaffService(mockRepo, mockJWT)
			result, err := service.GetAll(tt.schemaName)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, tt.expectedCount)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestStaffService_GetByID(t *testing.T) {
	tests := []struct {
		name        string
		id          uint
		schemaName  string
		mockStaff   *domain.Staff
		mockError   error
		expectError bool
	}{
		{
			name:       "successful get by ID",
			id:         1,
			schemaName: "tenant_test",
			mockStaff: &domain.Staff{
				Username:  "admin",
				StaffCode: "STF001",
				Email:     "admin@test.com",
			},
			mockError:   nil,
			expectError: false,
		},
		{
			name:        "staff not found",
			id:          999,
			schemaName:  "tenant_test",
			mockStaff:   nil,
			mockError:   errors.New("record not found"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewMockStaffRepository()
			mockJWT := mocks.NewMockJWTService()

			mockRepo.On("GetByID", tt.id, tt.schemaName).Return(tt.mockStaff, tt.mockError)

			service := services.NewStaffService(mockRepo, mockJWT)
			result, err := service.GetByID(tt.id, tt.schemaName)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.mockStaff.Username, result.Username)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestStaffService_Create(t *testing.T) {
	tests := []struct {
		name        string
		request     *domain.StaffCreateRequest
		schemaName  string
		createError error
		expectError bool
	}{
		{
			name: "successful create",
			request: &domain.StaffCreateRequest{
				Username:    "newuser",
				Password:    "password123",
				PhoneNumber: "0812345678",
				Email:       "newuser@test.com",
				FirstName:   "New",
				LastName:    "User",
				IsAdmin:     false,
			},
			schemaName:  "tenant_test",
			createError: nil,
			expectError: false,
		},
		{
			name: "duplicate username error",
			request: &domain.StaffCreateRequest{
				Username:    "existinguser",
				Password:    "password123",
				PhoneNumber: "0812345678",
				Email:       "existing@test.com",
				FirstName:   "Existing",
				LastName:    "User",
			},
			schemaName:  "tenant_test",
			createError: errors.New("duplicate key"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewMockStaffRepository()
			mockJWT := mocks.NewMockJWTService()

			mockRepo.On("Create", mock.AnythingOfType("*domain.Staff"), tt.schemaName).Return(tt.createError)

			service := services.NewStaffService(mockRepo, mockJWT)
			result, err := service.Create(tt.request, tt.schemaName)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.request.Username, result.Username)
				assert.Equal(t, tt.request.Email, result.Email)
				// Verify password is hashed
				assert.NotEqual(t, tt.request.Password, result.Password)
				// Verify staff code is generated
				assert.NotEmpty(t, result.StaffCode)
				assert.Contains(t, result.StaffCode, "STF")
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestStaffService_Update(t *testing.T) {
	existingStaff := &domain.Staff{
		Username:    "admin",
		StaffCode:   "STF001",
		PhoneNumber: "0812345678",
		Email:       "admin@test.com",
		FirstName:   "Admin",
		LastName:    "User",
	}
	existingStaff.ID = 1

	tests := []struct {
		name        string
		id          uint
		request     *domain.StaffUpdateRequest
		schemaName  string
		mockStaff   *domain.Staff
		getError    error
		updateError error
		expectError bool
	}{
		{
			name: "successful update",
			id:   1,
			request: &domain.StaffUpdateRequest{
				StaffCode:   "STF001",
				PhoneNumber: "0899999999",
				Email:       "updated@test.com",
				FirstName:   "Updated",
				LastName:    "Name",
				IsAdmin:     true,
			},
			schemaName:  "tenant_test",
			mockStaff:   existingStaff,
			getError:    nil,
			updateError: nil,
			expectError: false,
		},
		{
			name: "staff not found",
			id:   999,
			request: &domain.StaffUpdateRequest{
				StaffCode:   "STF001",
				PhoneNumber: "0899999999",
				Email:       "updated@test.com",
				FirstName:   "Updated",
				LastName:    "Name",
			},
			schemaName:  "tenant_test",
			mockStaff:   nil,
			getError:    errors.New("record not found"),
			updateError: nil,
			expectError: true,
		},
		{
			name: "update error",
			id:   1,
			request: &domain.StaffUpdateRequest{
				StaffCode:   "STF001",
				PhoneNumber: "0899999999",
				Email:       "duplicate@test.com",
				FirstName:   "Updated",
				LastName:    "Name",
			},
			schemaName:  "tenant_test",
			mockStaff:   existingStaff,
			getError:    nil,
			updateError: errors.New("duplicate key"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewMockStaffRepository()
			mockJWT := mocks.NewMockJWTService()

			mockRepo.On("GetByID", tt.id, tt.schemaName).Return(tt.mockStaff, tt.getError)

			if tt.getError == nil {
				mockRepo.On("Update", mock.AnythingOfType("*domain.Staff"), tt.schemaName).Return(tt.updateError)
			}

			service := services.NewStaffService(mockRepo, mockJWT)
			result, err := service.Update(tt.id, tt.request, tt.schemaName)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.request.FirstName, result.FirstName)
				assert.Equal(t, tt.request.Email, result.Email)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestStaffService_Delete(t *testing.T) {
	tests := []struct {
		name        string
		id          uint
		schemaName  string
		deleteError error
		expectError bool
	}{
		{
			name:        "successful delete",
			id:          1,
			schemaName:  "tenant_test",
			deleteError: nil,
			expectError: false,
		},
		{
			name:        "staff not found",
			id:          999,
			schemaName:  "tenant_test",
			deleteError: errors.New("record not found"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewMockStaffRepository()
			mockJWT := mocks.NewMockJWTService()

			mockRepo.On("Delete", tt.id, tt.schemaName).Return(tt.deleteError)

			service := services.NewStaffService(mockRepo, mockJWT)
			err := service.Delete(tt.id, tt.schemaName)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
