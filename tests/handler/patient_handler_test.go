package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/wichai2002/his_v1/internal/delivery/http/handler"
	"github.com/wichai2002/his_v1/internal/delivery/http/middleware"
	"github.com/wichai2002/his_v1/internal/domain"
	"github.com/wichai2002/his_v1/internal/mocks"
	"github.com/wichai2002/his_v1/pkg/utils"
)

const testSchemaName = "tenant_test"

// setupPatientRouter creates a test router with tenant middleware
func setupPatientRouter(mockService *mocks.MockPatientService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Add tenant schema to context for testing
	router.Use(func(c *gin.Context) {
		c.Set(middleware.TenantSchemaKey, testSchemaName)
		c.Next()
	})

	patientHandler := handler.NewPatientHandler(mockService)

	// Setup routes
	patients := router.Group("/patients")
	{
		patients.GET("/search", patientHandler.Search)
		patients.POST("", patientHandler.Create)
		patients.PUT("/:id", patientHandler.Update)
		patients.PATCH("/:id", patientHandler.PartialUpdate)
		patients.DELETE("/:id", patientHandler.Delete)
	}

	return router
}

// ==================== SEARCH TESTS ====================

func TestPatientHandler_Search_Success(t *testing.T) {
	mockService := mocks.NewMockPatientService()
	router := setupPatientRouter(mockService)

	expectedPatients := []domain.Patient{
		{
			FirstNameEN: "John",
			LastNameEN:  "Doe",
			PatientHN:   "HOSP0001-00000001",
			DateOfBirth: time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC),
			Gender:      domain.Male,
			BloodGrp:    domain.A,
		},
		{
			FirstNameEN: "Jane",
			LastNameEN:  "Smith",
			PatientHN:   "HOSP0001-00000002",
			DateOfBirth: time.Date(1985, 5, 20, 0, 0, 0, 0, time.UTC),
			Gender:      domain.Female,
			BloodGrp:    domain.B,
		},
	}

	mockService.On("Search", "John", testSchemaName).Return(expectedPatients, nil)

	req, _ := http.NewRequest("GET", "/patients/search?query=John", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var response utils.Response
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, "success", response.Message)

	mockService.AssertExpectations(t)
}

func TestPatientHandler_Search_EmptyQuery(t *testing.T) {
	mockService := mocks.NewMockPatientService()
	router := setupPatientRouter(mockService)

	expectedPatients := []domain.Patient{
		{FirstNameEN: "John", LastNameEN: "Doe", PatientHN: "HOSP0001-00000001"},
	}

	mockService.On("Search", "", testSchemaName).Return(expectedPatients, nil)

	req, _ := http.NewRequest("GET", "/patients/search", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	mockService.AssertExpectations(t)
}

func TestPatientHandler_Search_NotFound(t *testing.T) {
	mockService := mocks.NewMockPatientService()
	router := setupPatientRouter(mockService)

	// Use domain.ErrNotFound for proper error handling
	mockService.On("Search", "NonExistent", testSchemaName).Return(nil, domain.ErrNotFound)

	req, _ := http.NewRequest("GET", "/patients/search?query=NonExistent", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)

	var response utils.Response
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Equal(t, "patient not found", response.Error)

	mockService.AssertExpectations(t)
}

// ==================== CREATE TESTS ====================

func TestPatientHandler_Create_Success(t *testing.T) {
	mockService := mocks.NewMockPatientService()
	router := setupPatientRouter(mockService)

	createRequest := domain.PatientCreateRequest{
		FirstNameTH: "สมชาย",
		LastNameTH:  "ใจดี",
		FirstNameEN: "Somchai",
		LastNameEN:  "Jaidee",
		DateOfBirth: "1990-01-15",
		NickNameTH:  "ชาย",
		NickNameEN:  "Chai",
		NationalID:  "1234567890123",
		PassportID:  "AB1234567",
		PhoneNumber: "0812345678",
		Email:       "somchai@example.com",
		Gender:      domain.Male,
		Nationality: "Thai",
		BloodGrp:    "A",
	}

	expectedPatient := &domain.Patient{
		FirstNameTH: "สมชาย",
		LastNameTH:  "ใจดี",
		FirstNameEN: "Somchai",
		LastNameEN:  "Jaidee",
		DateOfBirth: time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC),
		NationalID:  "1234567890123",
		PhoneNumber: "0812345678",
		Email:       "somchai@example.com",
		Gender:      domain.Male,
		Nationality: "Thai",
		BloodGrp:    domain.A,
		PatientHN:   "HOSP0001-00000001",
	}

	mockService.On("Create", mock.AnythingOfType("*domain.PatientCreateRequest"), testSchemaName).Return(expectedPatient, nil)

	body, _ := json.Marshal(createRequest)
	req, _ := http.NewRequest("POST", "/patients", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)

	var response utils.Response
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, "patient created successfully", response.Message)

	mockService.AssertExpectations(t)
}

func TestPatientHandler_Create_InvalidJSON(t *testing.T) {
	mockService := mocks.NewMockPatientService()
	router := setupPatientRouter(mockService)

	req, _ := http.NewRequest("POST", "/patients", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)

	var response utils.Response
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
}

func TestPatientHandler_Create_MissingRequiredFields(t *testing.T) {
	mockService := mocks.NewMockPatientService()
	router := setupPatientRouter(mockService)

	// Missing required fields
	createRequest := map[string]interface{}{
		"first_name_en": "John",
		// Missing other required fields
	}

	body, _ := json.Marshal(createRequest)
	req, _ := http.NewRequest("POST", "/patients", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)

	var response utils.Response
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
}

func TestPatientHandler_Create_InvalidEmail(t *testing.T) {
	mockService := mocks.NewMockPatientService()
	router := setupPatientRouter(mockService)

	createRequest := map[string]interface{}{
		"first_name_th": "สมชาย",
		"last_name_th":  "ใจดี",
		"first_name_en": "Somchai",
		"last_name_en":  "Jaidee",
		"date_of_birth": "1990-01-15",
		"nick_name_th":  "ชาย",
		"nick_name_en":  "Chai",
		"national_id":   "1234567890123",
		"passport_id":   "AB1234567",
		"phone_number":  "0812345678",
		"email":         "invalid-email", // Invalid email
		"gender":        "M",
		"nationality":   "Thai",
		"blood_grp":     "A",
	}

	body, _ := json.Marshal(createRequest)
	req, _ := http.NewRequest("POST", "/patients", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestPatientHandler_Create_InvalidGender(t *testing.T) {
	mockService := mocks.NewMockPatientService()
	router := setupPatientRouter(mockService)

	createRequest := map[string]interface{}{
		"first_name_th": "สมชาย",
		"last_name_th":  "ใจดี",
		"first_name_en": "Somchai",
		"last_name_en":  "Jaidee",
		"date_of_birth": "1990-01-15",
		"nick_name_th":  "ชาย",
		"nick_name_en":  "Chai",
		"national_id":   "1234567890123",
		"passport_id":   "AB1234567",
		"phone_number":  "0812345678",
		"email":         "test@example.com",
		"gender":        "INVALID", // Invalid gender
		"nationality":   "Thai",
		"blood_grp":     "A",
	}

	body, _ := json.Marshal(createRequest)
	req, _ := http.NewRequest("POST", "/patients", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestPatientHandler_Create_InvalidBloodGroup(t *testing.T) {
	mockService := mocks.NewMockPatientService()
	router := setupPatientRouter(mockService)

	createRequest := map[string]interface{}{
		"first_name_th": "สมชาย",
		"last_name_th":  "ใจดี",
		"first_name_en": "Somchai",
		"last_name_en":  "Jaidee",
		"date_of_birth": "1990-01-15",
		"nick_name_th":  "ชาย",
		"nick_name_en":  "Chai",
		"national_id":   "1234567890123",
		"passport_id":   "AB1234567",
		"phone_number":  "0812345678",
		"email":         "test@example.com",
		"gender":        "M",
		"nationality":   "Thai",
		"blood_grp":     "X", // Invalid blood group
	}

	body, _ := json.Marshal(createRequest)
	req, _ := http.NewRequest("POST", "/patients", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestPatientHandler_Create_ServiceError(t *testing.T) {
	mockService := mocks.NewMockPatientService()
	router := setupPatientRouter(mockService)

	createRequest := domain.PatientCreateRequest{
		FirstNameTH: "สมชาย",
		LastNameTH:  "ใจดี",
		FirstNameEN: "Somchai",
		LastNameEN:  "Jaidee",
		DateOfBirth: "1990-01-15",
		NickNameTH:  "ชาย",
		NickNameEN:  "Chai",
		NationalID:  "1234567890123",
		PassportID:  "AB1234567",
		PhoneNumber: "0812345678",
		Email:       "somchai@example.com",
		Gender:      domain.Male,
		Nationality: "Thai",
		BloodGrp:    "A",
	}

	// Use domain.ErrDuplicateEntry for proper error handling
	mockService.On("Create", mock.AnythingOfType("*domain.PatientCreateRequest"), testSchemaName).Return(nil, domain.ErrDuplicateEntry)

	body, _ := json.Marshal(createRequest)
	req, _ := http.NewRequest("POST", "/patients", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusConflict, resp.Code)

	var response utils.Response
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Equal(t, "duplicate patient entry", response.Error)

	mockService.AssertExpectations(t)
}

// ==================== UPDATE TESTS ====================

func TestPatientHandler_Update_Success(t *testing.T) {
	mockService := mocks.NewMockPatientService()
	router := setupPatientRouter(mockService)

	updateRequest := domain.PatientUpdateRequest{
		FirstNameTH: "สมหญิง",
		LastNameTH:  "ใจดี",
		FirstNameEN: "Jane",
		LastNameEN:  "Jaidee",
		DateOfBirth: "1990-01-15",
		Gender:      "F",
		Nationality: "Thai",
		BloodGrp:    "B",
	}

	expectedPatient := &domain.Patient{
		FirstNameTH: "สมหญิง",
		LastNameTH:  "ใจดี",
		FirstNameEN: "Jane",
		LastNameEN:  "Jaidee",
		DateOfBirth: time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC),
		Gender:      domain.Female,
		Nationality: "Thai",
		BloodGrp:    domain.B,
		PatientHN:   "HOSP0001-00000001",
	}

	mockService.On("Update", uint(1), mock.AnythingOfType("*domain.PatientUpdateRequest"), testSchemaName).Return(expectedPatient, nil)

	body, _ := json.Marshal(updateRequest)
	req, _ := http.NewRequest("PUT", "/patients/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var response utils.Response
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, "patient updated successfully", response.Message)

	mockService.AssertExpectations(t)
}

func TestPatientHandler_Update_InvalidID(t *testing.T) {
	mockService := mocks.NewMockPatientService()
	router := setupPatientRouter(mockService)

	updateRequest := domain.PatientUpdateRequest{
		FirstNameEN: "Jane",
		LastNameEN:  "Doe",
		DateOfBirth: "1990-01-15",
		Gender:      "F",
		Nationality: "Thai",
		BloodGrp:    "B",
	}

	body, _ := json.Marshal(updateRequest)
	req, _ := http.NewRequest("PUT", "/patients/invalid", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)

	var response utils.Response
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Equal(t, "invalid id", response.Error)
}

func TestPatientHandler_Update_NotFound(t *testing.T) {
	mockService := mocks.NewMockPatientService()
	router := setupPatientRouter(mockService)

	updateRequest := domain.PatientUpdateRequest{
		FirstNameTH: "สมหญิง",
		LastNameTH:  "ใจดี",
		FirstNameEN: "Jane",
		LastNameEN:  "Jaidee",
		DateOfBirth: "1990-01-15",
		Gender:      "F",
		Nationality: "Thai",
		BloodGrp:    "B",
	}

	// Use domain.ErrNotFound for proper error handling
	mockService.On("Update", uint(999), mock.AnythingOfType("*domain.PatientUpdateRequest"), testSchemaName).Return(nil, domain.ErrNotFound)

	body, _ := json.Marshal(updateRequest)
	req, _ := http.NewRequest("PUT", "/patients/999", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)

	var response utils.Response
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Equal(t, "patient not found", response.Error)

	mockService.AssertExpectations(t)
}

func TestPatientHandler_Update_InvalidJSON(t *testing.T) {
	mockService := mocks.NewMockPatientService()
	router := setupPatientRouter(mockService)

	req, _ := http.NewRequest("PUT", "/patients/1", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestPatientHandler_Update_MissingRequiredFields(t *testing.T) {
	mockService := mocks.NewMockPatientService()
	router := setupPatientRouter(mockService)

	// Missing required fields for PUT
	updateRequest := map[string]interface{}{
		"first_name_en": "Jane",
		// Missing other required fields like gender, nationality, blood_grp
	}

	body, _ := json.Marshal(updateRequest)
	req, _ := http.NewRequest("PUT", "/patients/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

// ==================== PARTIAL UPDATE TESTS ====================

func TestPatientHandler_PartialUpdate_Success(t *testing.T) {
	mockService := mocks.NewMockPatientService()
	router := setupPatientRouter(mockService)

	firstName := "UpdatedName"
	partialUpdateRequest := domain.PatientPartialUpdateRequest{
		FirstNameEN: &firstName,
	}

	expectedPatient := &domain.Patient{
		FirstNameEN: "UpdatedName",
		LastNameEN:  "Doe",
		PatientHN:   "HOSP0001-00000001",
	}

	mockService.On("PartialUpdate", uint(1), mock.AnythingOfType("*domain.PatientPartialUpdateRequest"), testSchemaName).Return(expectedPatient, nil)

	body, _ := json.Marshal(partialUpdateRequest)
	req, _ := http.NewRequest("PATCH", "/patients/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var response utils.Response
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, "patient updated successfully", response.Message)

	mockService.AssertExpectations(t)
}

func TestPatientHandler_PartialUpdate_MultipleFields(t *testing.T) {
	mockService := mocks.NewMockPatientService()
	router := setupPatientRouter(mockService)

	firstName := "UpdatedFirstName"
	lastName := "UpdatedLastName"
	email := "updated@example.com"
	partialUpdateRequest := domain.PatientPartialUpdateRequest{
		FirstNameEN: &firstName,
		LastNameEN:  &lastName,
		Email:       &email,
	}

	expectedPatient := &domain.Patient{
		FirstNameEN: "UpdatedFirstName",
		LastNameEN:  "UpdatedLastName",
		Email:       "updated@example.com",
		PatientHN:   "HOSP0001-00000001",
	}

	mockService.On("PartialUpdate", uint(1), mock.AnythingOfType("*domain.PatientPartialUpdateRequest"), testSchemaName).Return(expectedPatient, nil)

	body, _ := json.Marshal(partialUpdateRequest)
	req, _ := http.NewRequest("PATCH", "/patients/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	mockService.AssertExpectations(t)
}

func TestPatientHandler_PartialUpdate_InvalidID(t *testing.T) {
	mockService := mocks.NewMockPatientService()
	router := setupPatientRouter(mockService)

	firstName := "UpdatedName"
	partialUpdateRequest := domain.PatientPartialUpdateRequest{
		FirstNameEN: &firstName,
	}

	body, _ := json.Marshal(partialUpdateRequest)
	req, _ := http.NewRequest("PATCH", "/patients/invalid", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)

	var response utils.Response
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Equal(t, "invalid id", response.Error)
}

func TestPatientHandler_PartialUpdate_NotFound(t *testing.T) {
	mockService := mocks.NewMockPatientService()
	router := setupPatientRouter(mockService)

	firstName := "UpdatedName"
	partialUpdateRequest := domain.PatientPartialUpdateRequest{
		FirstNameEN: &firstName,
	}

	// Use domain.ErrNotFound for proper error handling
	mockService.On("PartialUpdate", uint(999), mock.AnythingOfType("*domain.PatientPartialUpdateRequest"), testSchemaName).Return(nil, domain.ErrNotFound)

	body, _ := json.Marshal(partialUpdateRequest)
	req, _ := http.NewRequest("PATCH", "/patients/999", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)

	var response utils.Response
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Equal(t, "patient not found", response.Error)

	mockService.AssertExpectations(t)
}

func TestPatientHandler_PartialUpdate_InvalidJSON(t *testing.T) {
	mockService := mocks.NewMockPatientService()
	router := setupPatientRouter(mockService)

	req, _ := http.NewRequest("PATCH", "/patients/1", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestPatientHandler_PartialUpdate_InvalidEmail(t *testing.T) {
	mockService := mocks.NewMockPatientService()
	router := setupPatientRouter(mockService)

	invalidEmail := "invalid-email"
	partialUpdateRequest := domain.PatientPartialUpdateRequest{
		Email: &invalidEmail,
	}

	body, _ := json.Marshal(partialUpdateRequest)
	req, _ := http.NewRequest("PATCH", "/patients/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestPatientHandler_PartialUpdate_InvalidGender(t *testing.T) {
	mockService := mocks.NewMockPatientService()
	router := setupPatientRouter(mockService)

	invalidGender := "INVALID"
	partialUpdateRequest := domain.PatientPartialUpdateRequest{
		Gender: &invalidGender,
	}

	body, _ := json.Marshal(partialUpdateRequest)
	req, _ := http.NewRequest("PATCH", "/patients/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

// ==================== DELETE TESTS ====================

func TestPatientHandler_Delete_Success(t *testing.T) {
	mockService := mocks.NewMockPatientService()
	router := setupPatientRouter(mockService)

	mockService.On("Delete", uint(1), testSchemaName).Return(nil)

	req, _ := http.NewRequest("DELETE", "/patients/1", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var response utils.Response
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, "patient deleted successfully", response.Message)

	mockService.AssertExpectations(t)
}

func TestPatientHandler_Delete_InvalidID(t *testing.T) {
	mockService := mocks.NewMockPatientService()
	router := setupPatientRouter(mockService)

	req, _ := http.NewRequest("DELETE", "/patients/invalid", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)

	var response utils.Response
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Equal(t, "invalid id", response.Error)
}

func TestPatientHandler_Delete_NotFound(t *testing.T) {
	mockService := mocks.NewMockPatientService()
	router := setupPatientRouter(mockService)

	// Use domain.ErrNotFound for proper error handling
	mockService.On("Delete", uint(999), testSchemaName).Return(domain.ErrNotFound)

	req, _ := http.NewRequest("DELETE", "/patients/999", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)

	var response utils.Response
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Equal(t, "patient not found", response.Error)

	mockService.AssertExpectations(t)
}

func TestPatientHandler_Delete_ZeroID(t *testing.T) {
	mockService := mocks.NewMockPatientService()
	router := setupPatientRouter(mockService)

	// Use domain.ErrNotFound as ID 0 won't match any record
	mockService.On("Delete", uint(0), testSchemaName).Return(domain.ErrNotFound)

	req, _ := http.NewRequest("DELETE", "/patients/0", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)
	mockService.AssertExpectations(t)
}

// ==================== EDGE CASE TESTS ====================

func TestPatientHandler_Create_MaxLengthFields(t *testing.T) {
	mockService := mocks.NewMockPatientService()
	router := setupPatientRouter(mockService)

	// Create string longer than max allowed (255 chars)
	longString := string(make([]byte, 300))
	for i := range longString {
		longString = longString[:i] + "a" + longString[i+1:]
	}

	createRequest := map[string]interface{}{
		"first_name_th": longString, // Exceeds max=255
		"last_name_th":  "ใจดี",
		"first_name_en": "Somchai",
		"last_name_en":  "Jaidee",
		"date_of_birth": "1990-01-15",
		"nick_name_th":  "ชาย",
		"nick_name_en":  "Chai",
		"national_id":   "1234567890123",
		"passport_id":   "AB1234567",
		"phone_number":  "0812345678",
		"email":         "test@example.com",
		"gender":        "M",
		"nationality":   "Thai",
		"blood_grp":     "A",
	}

	body, _ := json.Marshal(createRequest)
	req, _ := http.NewRequest("POST", "/patients", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestPatientHandler_Create_EmptyBody(t *testing.T) {
	mockService := mocks.NewMockPatientService()
	router := setupPatientRouter(mockService)

	req, _ := http.NewRequest("POST", "/patients", bytes.NewBuffer([]byte("{}")))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestPatientHandler_Update_EmptyBody(t *testing.T) {
	mockService := mocks.NewMockPatientService()
	router := setupPatientRouter(mockService)

	req, _ := http.NewRequest("PUT", "/patients/1", bytes.NewBuffer([]byte("{}")))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	// Should fail validation due to missing required fields
	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestPatientHandler_PartialUpdate_EmptyBody(t *testing.T) {
	mockService := mocks.NewMockPatientService()
	router := setupPatientRouter(mockService)

	expectedPatient := &domain.Patient{
		FirstNameEN: "John",
		LastNameEN:  "Doe",
		PatientHN:   "HOSP0001-00000001",
	}

	// Empty body is valid for PATCH - no fields to update
	mockService.On("PartialUpdate", uint(1), mock.AnythingOfType("*domain.PatientPartialUpdateRequest"), testSchemaName).Return(expectedPatient, nil)

	req, _ := http.NewRequest("PATCH", "/patients/1", bytes.NewBuffer([]byte("{}")))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	mockService.AssertExpectations(t)
}
