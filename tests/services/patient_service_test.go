package services_test

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/wichai2002/his_v1/internal/domain"
	"github.com/wichai2002/his_v1/internal/mocks"
	"github.com/wichai2002/his_v1/internal/services"
)

func TestPatientService_Search(t *testing.T) {
	tests := []struct {
		name          string
		query         string
		schemaName    string
		mockPatients  []domain.Patient
		mockError     error
		expectedCount int
		expectError   bool
	}{
		{
			name:       "successful search with results",
			query:      "John",
			schemaName: "tenant_test",
			mockPatients: []domain.Patient{
				{FirstNameEN: "John", LastNameEN: "Doe", PatientHN: "HOSP0001-00000001"},
				{FirstNameEN: "Johnny", LastNameEN: "Smith", PatientHN: "HOSP0001-00000002"},
			},
			mockError:     nil,
			expectedCount: 2,
			expectError:   false,
		},
		{
			name:          "successful search with no results",
			query:         "NonExistent",
			schemaName:    "tenant_test",
			mockPatients:  []domain.Patient{},
			mockError:     nil,
			expectedCount: 0,
			expectError:   false,
		},
		{
			name:          "search with empty query returns all",
			query:         "",
			schemaName:    "tenant_test",
			mockPatients:  []domain.Patient{{FirstNameEN: "John"}},
			mockError:     nil,
			expectedCount: 1,
			expectError:   false,
		},
		{
			name:          "search with database error",
			query:         "test",
			schemaName:    "tenant_test",
			mockPatients:  nil,
			mockError:     errors.New("database error"),
			expectedCount: 0,
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewMockPatientRepository()
			mockTenantService := mocks.NewMockTenantService()

			mockRepo.On("Search", tt.query, tt.schemaName).Return(tt.mockPatients, tt.mockError)

			service := services.NewPatientService(mockRepo, mockTenantService)
			result, err := service.Search(tt.query, tt.schemaName)

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

func TestPatientService_SearchByID(t *testing.T) {
	tests := []struct {
		name        string
		id          uint
		schemaName  string
		mockPatient *domain.Patient
		mockError   error
		expectError bool
	}{
		{
			name:       "successful search by ID",
			id:         1,
			schemaName: "tenant_test",
			mockPatient: &domain.Patient{
				FirstNameEN: "John",
				LastNameEN:  "Doe",
				PatientHN:   "HOSP0001-00000001",
			},
			mockError:   nil,
			expectError: false,
		},
		{
			name:        "patient not found",
			id:          999,
			schemaName:  "tenant_test",
			mockPatient: nil,
			mockError:   errors.New("record not found"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewMockPatientRepository()
			mockTenantService := mocks.NewMockTenantService()

			mockRepo.On("SearchByID", tt.id, tt.schemaName).Return(tt.mockPatient, tt.mockError)

			service := services.NewPatientService(mockRepo, mockTenantService)
			result, err := service.SearchByID(tt.id, tt.schemaName)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.mockPatient.FirstNameEN, result.FirstNameEN)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestPatientService_Create(t *testing.T) {
	tests := []struct {
		name         string
		request      *domain.PatientCreateRequest
		schemaName   string
		generatedHN  string
		hnError      error
		createError  error
		expectError  bool
		errorMessage string
	}{
		{
			name: "successful patient creation",
			request: &domain.PatientCreateRequest{
				FirstNameTH:  "สมชาย",
				LastNameTH:   "ใจดี",
				FirstNameEN:  "Somchai",
				LastNameEN:   "Jaidee",
				DateOfBirth:  "1990-01-15",
				NationalID:   "1234567890123",
				PhoneNumber:  "0812345678",
				Email:        "somchai@example.com",
				Gender:       domain.Male,
				Nationality:  "Thai",
				BloodGrp:     "A",
			},
			schemaName:  "tenant_test",
			generatedHN: "HOSP0001-00000001",
			hnError:     nil,
			createError: nil,
			expectError: false,
		},
		{
			name: "failed HN generation",
			request: &domain.PatientCreateRequest{
				FirstNameEN: "John",
				LastNameEN:  "Doe",
				DateOfBirth: "1990-01-15",
				Gender:      domain.Male,
				Nationality: "Thai",
				BloodGrp:    "A",
			},
			schemaName:   "tenant_test",
			generatedHN:  "",
			hnError:      errors.New("failed to generate HN"),
			createError:  nil,
			expectError:  true,
			errorMessage: "failed to generate HN",
		},
		{
			name: "invalid date of birth format",
			request: &domain.PatientCreateRequest{
				FirstNameEN: "John",
				LastNameEN:  "Doe",
				DateOfBirth: "invalid-date",
				Gender:      domain.Male,
				Nationality: "Thai",
				BloodGrp:    "A",
			},
			schemaName:  "tenant_test",
			generatedHN: "HOSP0001-00000001",
			hnError:     nil,
			createError: nil,
			expectError: true,
		},
		{
			name: "database error on create",
			request: &domain.PatientCreateRequest{
				FirstNameEN: "John",
				LastNameEN:  "Doe",
				DateOfBirth: "1990-01-15",
				Gender:      domain.Male,
				Nationality: "Thai",
				BloodGrp:    "A",
			},
			schemaName:   "tenant_test",
			generatedHN:  "HOSP0001-00000001",
			hnError:      nil,
			createError:  errors.New("duplicate key"),
			expectError:  true,
			errorMessage: "duplicate key",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewMockPatientRepository()
			mockTenantService := mocks.NewMockTenantService()

			// Date validation happens first, so GenerateHN is only called for valid dates
			if tt.request.DateOfBirth != "invalid-date" {
				mockTenantService.On("GenerateHN", tt.schemaName).Return(tt.generatedHN, tt.hnError)

				if tt.hnError == nil {
					mockRepo.On("Create", mock.AnythingOfType("*domain.Patient"), tt.schemaName).Return(tt.createError)
				}
			}

			service := services.NewPatientService(mockRepo, mockTenantService)
			result, err := service.Create(tt.request, tt.schemaName)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
				if tt.errorMessage != "" {
					assert.Contains(t, err.Error(), tt.errorMessage)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.generatedHN, result.PatientHN)
				assert.Equal(t, tt.request.FirstNameEN, result.FirstNameEN)
			}

			mockTenantService.AssertExpectations(t)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestPatientService_Update(t *testing.T) {
	existingPatient := &domain.Patient{
		FirstNameEN: "John",
		LastNameEN:  "Doe",
		PatientHN:   "HOSP0001-00000001",
		DateOfBirth: time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC),
	}
	existingPatient.ID = 1

	tests := []struct {
		name        string
		id          uint
		request     *domain.PatientUpdateRequest
		schemaName  string
		mockPatient *domain.Patient
		getError    error
		updateError error
		expectError bool
	}{
		{
			name: "successful update",
			id:   1,
			request: &domain.PatientUpdateRequest{
				FirstNameTH: "สมหญิง",
				LastNameTH:  "ใจดี",
				FirstNameEN: "Jane",
				LastNameEN:  "Doe",
				DateOfBirth: "1990-01-15",
				Gender:      "F",
				Nationality: "Thai",
				BloodGrp:    "B",
			},
			schemaName:  "tenant_test",
			mockPatient: existingPatient,
			getError:    nil,
			updateError: nil,
			expectError: false,
		},
		{
			name: "patient not found",
			id:   999,
			request: &domain.PatientUpdateRequest{
				FirstNameEN: "Jane",
				LastNameEN:  "Doe",
				DateOfBirth: "1990-01-15",
				Gender:      "F",
				Nationality: "Thai",
				BloodGrp:    "B",
			},
			schemaName:  "tenant_test",
			mockPatient: nil,
			getError:    errors.New("record not found"),
			updateError: nil,
			expectError: true,
		},
		{
			name: "invalid date format",
			id:   1,
			request: &domain.PatientUpdateRequest{
				FirstNameEN: "Jane",
				LastNameEN:  "Doe",
				DateOfBirth: "invalid-date",
				Gender:      "F",
				Nationality: "Thai",
				BloodGrp:    "B",
			},
			schemaName:  "tenant_test",
			mockPatient: existingPatient,
			getError:    nil,
			updateError: nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewMockPatientRepository()
			mockTenantService := mocks.NewMockTenantService()

			mockRepo.On("GetByID", tt.id, tt.schemaName).Return(tt.mockPatient, tt.getError)

			if tt.getError == nil && tt.request.DateOfBirth != "invalid-date" {
				mockRepo.On("Update", mock.AnythingOfType("*domain.Patient"), tt.schemaName).Return(tt.updateError)
			}

			service := services.NewPatientService(mockRepo, mockTenantService)
			result, err := service.Update(tt.id, tt.request, tt.schemaName)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.request.FirstNameEN, result.FirstNameEN)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestPatientService_PartialUpdate(t *testing.T) {
	existingPatient := &domain.Patient{
		FirstNameEN: "John",
		LastNameEN:  "Doe",
		PatientHN:   "HOSP0001-00000001",
	}
	existingPatient.ID = 1

	newFirstName := "Jane"

	tests := []struct {
		name           string
		id             uint
		request        *domain.PatientPartialUpdateRequest
		schemaName     string
		getPatient     *domain.Patient
		getError       error
		updateError    error
		getAfterUpdate *domain.Patient
		getAfterError  error
		expectError    bool
	}{
		{
			name: "successful partial update",
			id:   1,
			request: &domain.PatientPartialUpdateRequest{
				FirstNameEN: &newFirstName,
			},
			schemaName:     "tenant_test",
			getPatient:     existingPatient,
			getError:       nil,
			updateError:    nil,
			getAfterUpdate: &domain.Patient{FirstNameEN: "Jane", LastNameEN: "Doe", PatientHN: "HOSP0001-00000001"},
			getAfterError:  nil,
			expectError:    false,
		},
		{
			name: "patient not found",
			id:   999,
			request: &domain.PatientPartialUpdateRequest{
				FirstNameEN: &newFirstName,
			},
			schemaName:  "tenant_test",
			getPatient:  nil,
			getError:    errors.New("record not found"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewMockPatientRepository()
			mockTenantService := mocks.NewMockTenantService()

			// New logic: PartialUpdate is called first, then GetByID is called to fetch the updated record
			if tt.getError != nil {
				// Simulate not found error from PartialUpdate (repository returns error)
				mockRepo.On("PartialUpdate", tt.id, mock.Anything, tt.schemaName).Return(tt.getError)
			} else {
				mockRepo.On("PartialUpdate", tt.id, mock.Anything, tt.schemaName).Return(tt.updateError)
				if tt.updateError == nil {
					mockRepo.On("GetByID", tt.id, tt.schemaName).Return(tt.getAfterUpdate, tt.getAfterError)
				}
			}

			service := services.NewPatientService(mockRepo, mockTenantService)
			result, err := service.PartialUpdate(tt.id, tt.request, tt.schemaName)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestPatientService_Delete(t *testing.T) {
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
			name:        "patient not found",
			id:          999,
			schemaName:  "tenant_test",
			deleteError: errors.New("record not found"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewMockPatientRepository()
			mockTenantService := mocks.NewMockTenantService()

			// New logic: GetByID is called first to check if patient exists
			if tt.deleteError != nil {
				// Simulate not found error from GetByID
				mockRepo.On("GetByID", tt.id, tt.schemaName).Return(nil, tt.deleteError)
			} else {
				mockRepo.On("GetByID", tt.id, tt.schemaName).Return(&domain.Patient{}, nil)
				mockRepo.On("Delete", tt.id, tt.schemaName).Return(tt.deleteError)
			}

			service := services.NewPatientService(mockRepo, mockTenantService)
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
