package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/wichai2002/his_v1/internal/delivery/http/handler"
	"github.com/wichai2002/his_v1/internal/delivery/http/middleware"
	"github.com/wichai2002/his_v1/internal/domain"
	"github.com/wichai2002/his_v1/internal/mocks"
	"github.com/wichai2002/his_v1/pkg/utils"
)

// setupStaffRouter creates a test router with tenant middleware
func setupStaffRouter(mockService *mocks.MockStaffService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Add tenant schema to context for testing
	router.Use(func(c *gin.Context) {
		c.Set(middleware.TenantSchemaKey, testSchemaName)
		c.Next()
	})

	staffHandler := handler.NewStaffHandler(mockService)

	// Setup routes
	staff := router.Group("/staff")
	{
		staff.POST("/login", staffHandler.Login)
		staff.POST("/logout", staffHandler.Logout)
		staff.GET("", staffHandler.GetAll)
		staff.GET("/:id", staffHandler.GetByID)
		staff.POST("/create", staffHandler.Create)
		staff.PUT("/update/:id", staffHandler.Update)
		staff.DELETE("/delete/:id", staffHandler.Delete)
	}

	return router
}

// ==================== LOGIN TESTS ====================

func TestStaffHandler_Login_Success(t *testing.T) {
	mockService := mocks.NewMockStaffService()
	router := setupStaffRouter(mockService)

	loginRequest := domain.StaffLoginRequest{
		Username: "admin",
		Password: "password123",
	}

	expectedResponse := &domain.StaffLoginResponse{
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

	mockService.On("Login", mock.AnythingOfType("*domain.StaffLoginRequest"), testSchemaName).Return(expectedResponse, nil)

	body, _ := json.Marshal(loginRequest)
	req, _ := http.NewRequest("POST", "/staff/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var response utils.Response
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, "login successful", response.Message)

	mockService.AssertExpectations(t)
}

func TestStaffHandler_Login_InvalidCredentials(t *testing.T) {
	mockService := mocks.NewMockStaffService()
	router := setupStaffRouter(mockService)

	loginRequest := domain.StaffLoginRequest{
		Username: "admin",
		Password: "wrongpassword",
	}

	mockService.On("Login", mock.AnythingOfType("*domain.StaffLoginRequest"), testSchemaName).Return(nil, errors.New("invalid credentials"))

	body, _ := json.Marshal(loginRequest)
	req, _ := http.NewRequest("POST", "/staff/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusUnauthorized, resp.Code)

	var response utils.Response
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Equal(t, "invalid credentials", response.Error)

	mockService.AssertExpectations(t)
}

func TestStaffHandler_Login_UserNotFound(t *testing.T) {
	mockService := mocks.NewMockStaffService()
	router := setupStaffRouter(mockService)

	loginRequest := domain.StaffLoginRequest{
		Username: "nonexistent",
		Password: "password123",
	}

	mockService.On("Login", mock.AnythingOfType("*domain.StaffLoginRequest"), testSchemaName).Return(nil, errors.New("invalid credentials"))

	body, _ := json.Marshal(loginRequest)
	req, _ := http.NewRequest("POST", "/staff/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusUnauthorized, resp.Code)
	mockService.AssertExpectations(t)
}

func TestStaffHandler_Login_InvalidJSON(t *testing.T) {
	mockService := mocks.NewMockStaffService()
	router := setupStaffRouter(mockService)

	req, _ := http.NewRequest("POST", "/staff/login", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestStaffHandler_Login_MissingUsername(t *testing.T) {
	mockService := mocks.NewMockStaffService()
	router := setupStaffRouter(mockService)

	loginRequest := map[string]interface{}{
		"password": "password123",
		// Missing username
	}

	body, _ := json.Marshal(loginRequest)
	req, _ := http.NewRequest("POST", "/staff/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestStaffHandler_Login_MissingPassword(t *testing.T) {
	mockService := mocks.NewMockStaffService()
	router := setupStaffRouter(mockService)

	loginRequest := map[string]interface{}{
		"username": "admin",
		// Missing password
	}

	body, _ := json.Marshal(loginRequest)
	req, _ := http.NewRequest("POST", "/staff/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestStaffHandler_Login_ShortUsername(t *testing.T) {
	mockService := mocks.NewMockStaffService()
	router := setupStaffRouter(mockService)

	loginRequest := map[string]interface{}{
		"username": "ab",   // Less than min=5
		"password": "password123",
	}

	body, _ := json.Marshal(loginRequest)
	req, _ := http.NewRequest("POST", "/staff/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestStaffHandler_Login_ShortPassword(t *testing.T) {
	mockService := mocks.NewMockStaffService()
	router := setupStaffRouter(mockService)

	loginRequest := map[string]interface{}{
		"username": "admin",
		"password": "12345", // Less than min=6
	}

	body, _ := json.Marshal(loginRequest)
	req, _ := http.NewRequest("POST", "/staff/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

// ==================== LOGOUT TESTS ====================

func TestStaffHandler_Logout_Success(t *testing.T) {
	mockService := mocks.NewMockStaffService()
	router := setupStaffRouter(mockService)

	req, _ := http.NewRequest("POST", "/staff/logout", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var response utils.Response
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, "logout successful", response.Message)
}

// ==================== GET ALL TESTS ====================

func TestStaffHandler_GetAll_Success(t *testing.T) {
	mockService := mocks.NewMockStaffService()
	router := setupStaffRouter(mockService)

	expectedStaffs := []domain.Staff{
		{Username: "admin", StaffCode: "STF001", Email: "admin@test.com", FirstName: "Admin", LastName: "User", IsAdmin: true},
		{Username: "user1", StaffCode: "STF002", Email: "user1@test.com", FirstName: "User", LastName: "One", IsAdmin: false},
	}

	mockService.On("GetAll", testSchemaName).Return(expectedStaffs, nil)

	req, _ := http.NewRequest("GET", "/staff", nil)
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

func TestStaffHandler_GetAll_EmptyList(t *testing.T) {
	mockService := mocks.NewMockStaffService()
	router := setupStaffRouter(mockService)

	mockService.On("GetAll", testSchemaName).Return([]domain.Staff{}, nil)

	req, _ := http.NewRequest("GET", "/staff", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var response utils.Response
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)

	mockService.AssertExpectations(t)
}

func TestStaffHandler_GetAll_ServiceError(t *testing.T) {
	mockService := mocks.NewMockStaffService()
	router := setupStaffRouter(mockService)

	mockService.On("GetAll", testSchemaName).Return(nil, errors.New("database error"))

	req, _ := http.NewRequest("GET", "/staff", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)

	var response utils.Response
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Equal(t, "database error", response.Error)

	mockService.AssertExpectations(t)
}

// ==================== GET BY ID TESTS ====================

func TestStaffHandler_GetByID_Success(t *testing.T) {
	mockService := mocks.NewMockStaffService()
	router := setupStaffRouter(mockService)

	expectedStaff := &domain.Staff{
		Username:  "admin",
		StaffCode: "STF001",
		Email:     "admin@test.com",
		FirstName: "Admin",
		LastName:  "User",
		IsAdmin:   true,
	}

	mockService.On("GetByID", uint(1), testSchemaName).Return(expectedStaff, nil)

	req, _ := http.NewRequest("GET", "/staff/1", nil)
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

func TestStaffHandler_GetByID_InvalidID(t *testing.T) {
	mockService := mocks.NewMockStaffService()
	router := setupStaffRouter(mockService)

	req, _ := http.NewRequest("GET", "/staff/invalid", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)

	var response utils.Response
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Equal(t, "invalid id", response.Error)
}

func TestStaffHandler_GetByID_NotFound(t *testing.T) {
	mockService := mocks.NewMockStaffService()
	router := setupStaffRouter(mockService)

	mockService.On("GetByID", uint(999), testSchemaName).Return(nil, errors.New("record not found"))

	req, _ := http.NewRequest("GET", "/staff/999", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)

	var response utils.Response
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Equal(t, "staff not found", response.Error)

	mockService.AssertExpectations(t)
}

// ==================== CREATE TESTS ====================

func TestStaffHandler_Create_Success(t *testing.T) {
	mockService := mocks.NewMockStaffService()
	router := setupStaffRouter(mockService)

	createRequest := domain.StaffCreateRequest{
		Username:    "newuser",
		Password:    "password123",
		PhoneNumber: "0812345678",
		Email:       "newuser@test.com",
		FirstName:   "New",
		LastName:    "User",
		IsAdmin:     false,
	}

	expectedStaff := &domain.Staff{
		Username:    "newuser",
		StaffCode:   "STF003",
		PhoneNumber: "0812345678",
		Email:       "newuser@test.com",
		FirstName:   "New",
		LastName:    "User",
		IsAdmin:     false,
	}

	mockService.On("Create", mock.AnythingOfType("*domain.StaffCreateRequest"), testSchemaName).Return(expectedStaff, nil)

	body, _ := json.Marshal(createRequest)
	req, _ := http.NewRequest("POST", "/staff/create", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)

	var response utils.Response
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, "staff created successfully", response.Message)

	mockService.AssertExpectations(t)
}

func TestStaffHandler_Create_InvalidJSON(t *testing.T) {
	mockService := mocks.NewMockStaffService()
	router := setupStaffRouter(mockService)

	req, _ := http.NewRequest("POST", "/staff/create", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestStaffHandler_Create_MissingRequiredFields(t *testing.T) {
	mockService := mocks.NewMockStaffService()
	router := setupStaffRouter(mockService)

	createRequest := map[string]interface{}{
		"username": "newuser",
		// Missing password, email, phone_number, etc.
	}

	body, _ := json.Marshal(createRequest)
	req, _ := http.NewRequest("POST", "/staff/create", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestStaffHandler_Create_InvalidEmail(t *testing.T) {
	mockService := mocks.NewMockStaffService()
	router := setupStaffRouter(mockService)

	createRequest := map[string]interface{}{
		"username":     "newuser",
		"password":     "password123",
		"phone_number": "0812345678",
		"email":        "invalid-email", // Invalid email
		"first_name":   "New",
		"last_name":    "User",
	}

	body, _ := json.Marshal(createRequest)
	req, _ := http.NewRequest("POST", "/staff/create", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestStaffHandler_Create_ShortPassword(t *testing.T) {
	mockService := mocks.NewMockStaffService()
	router := setupStaffRouter(mockService)

	createRequest := map[string]interface{}{
		"username":     "newuser",
		"password":     "12345", // Less than min=6
		"phone_number": "0812345678",
		"email":        "newuser@test.com",
		"first_name":   "New",
		"last_name":    "User",
	}

	body, _ := json.Marshal(createRequest)
	req, _ := http.NewRequest("POST", "/staff/create", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestStaffHandler_Create_DuplicateUsername(t *testing.T) {
	mockService := mocks.NewMockStaffService()
	router := setupStaffRouter(mockService)

	createRequest := domain.StaffCreateRequest{
		Username:    "existinguser",
		Password:    "password123",
		PhoneNumber: "0812345678",
		Email:       "existing@test.com",
		FirstName:   "Existing",
		LastName:    "User",
	}

	mockService.On("Create", mock.AnythingOfType("*domain.StaffCreateRequest"), testSchemaName).Return(nil, errors.New("duplicate key"))

	body, _ := json.Marshal(createRequest)
	req, _ := http.NewRequest("POST", "/staff/create", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)

	var response utils.Response
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Equal(t, "duplicate key", response.Error)

	mockService.AssertExpectations(t)
}

// ==================== UPDATE TESTS ====================

func TestStaffHandler_Update_Success(t *testing.T) {
	mockService := mocks.NewMockStaffService()
	router := setupStaffRouter(mockService)

	updateRequest := domain.StaffUpdateRequest{
		StaffCode:   "STF001",
		PhoneNumber: "0899999999",
		Email:       "updated@test.com",
		FirstName:   "Updated",
		LastName:    "Name",
		IsAdmin:     true,
	}

	expectedStaff := &domain.Staff{
		Username:    "admin",
		StaffCode:   "STF001",
		PhoneNumber: "0899999999",
		Email:       "updated@test.com",
		FirstName:   "Updated",
		LastName:    "Name",
		IsAdmin:     true,
	}

	mockService.On("Update", uint(1), mock.AnythingOfType("*domain.StaffUpdateRequest"), testSchemaName).Return(expectedStaff, nil)

	body, _ := json.Marshal(updateRequest)
	req, _ := http.NewRequest("PUT", "/staff/update/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var response utils.Response
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, "staff updated successfully", response.Message)

	mockService.AssertExpectations(t)
}

func TestStaffHandler_Update_InvalidID(t *testing.T) {
	mockService := mocks.NewMockStaffService()
	router := setupStaffRouter(mockService)

	updateRequest := domain.StaffUpdateRequest{
		PhoneNumber: "0899999999",
		Email:       "updated@test.com",
		FirstName:   "Updated",
		LastName:    "Name",
	}

	body, _ := json.Marshal(updateRequest)
	req, _ := http.NewRequest("PUT", "/staff/update/invalid", bytes.NewBuffer(body))
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

func TestStaffHandler_Update_NotFound(t *testing.T) {
	mockService := mocks.NewMockStaffService()
	router := setupStaffRouter(mockService)

	updateRequest := domain.StaffUpdateRequest{
		StaffCode:   "STF001",
		PhoneNumber: "0899999999",
		Email:       "updated@test.com",
		FirstName:   "Updated",
		LastName:    "Name",
	}

	mockService.On("Update", uint(999), mock.AnythingOfType("*domain.StaffUpdateRequest"), testSchemaName).Return(nil, errors.New("record not found"))

	body, _ := json.Marshal(updateRequest)
	req, _ := http.NewRequest("PUT", "/staff/update/999", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	mockService.AssertExpectations(t)
}

func TestStaffHandler_Update_InvalidJSON(t *testing.T) {
	mockService := mocks.NewMockStaffService()
	router := setupStaffRouter(mockService)

	req, _ := http.NewRequest("PUT", "/staff/update/1", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestStaffHandler_Update_InvalidEmail(t *testing.T) {
	mockService := mocks.NewMockStaffService()
	router := setupStaffRouter(mockService)

	updateRequest := map[string]interface{}{
		"phone_number": "0899999999",
		"email":        "invalid-email", // Invalid email
		"first_name":   "Updated",
		"last_name":    "Name",
	}

	body, _ := json.Marshal(updateRequest)
	req, _ := http.NewRequest("PUT", "/staff/update/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestStaffHandler_Update_DuplicateEmail(t *testing.T) {
	mockService := mocks.NewMockStaffService()
	router := setupStaffRouter(mockService)

	updateRequest := domain.StaffUpdateRequest{
		StaffCode:   "STF001",
		PhoneNumber: "0899999999",
		Email:       "duplicate@test.com",
		FirstName:   "Updated",
		LastName:    "Name",
	}

	mockService.On("Update", uint(1), mock.AnythingOfType("*domain.StaffUpdateRequest"), testSchemaName).Return(nil, errors.New("duplicate key"))

	body, _ := json.Marshal(updateRequest)
	req, _ := http.NewRequest("PUT", "/staff/update/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	mockService.AssertExpectations(t)
}

// ==================== DELETE TESTS ====================

func TestStaffHandler_Delete_Success(t *testing.T) {
	mockService := mocks.NewMockStaffService()
	router := setupStaffRouter(mockService)

	mockService.On("Delete", uint(1), testSchemaName).Return(nil)

	req, _ := http.NewRequest("DELETE", "/staff/delete/1", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var response utils.Response
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, "staff deleted successfully", response.Message)

	mockService.AssertExpectations(t)
}

func TestStaffHandler_Delete_InvalidID(t *testing.T) {
	mockService := mocks.NewMockStaffService()
	router := setupStaffRouter(mockService)

	req, _ := http.NewRequest("DELETE", "/staff/delete/invalid", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)

	var response utils.Response
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Equal(t, "invalid id", response.Error)
}

func TestStaffHandler_Delete_NotFound(t *testing.T) {
	mockService := mocks.NewMockStaffService()
	router := setupStaffRouter(mockService)

	mockService.On("Delete", uint(999), testSchemaName).Return(errors.New("record not found"))

	req, _ := http.NewRequest("DELETE", "/staff/delete/999", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)

	var response utils.Response
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Equal(t, "record not found", response.Error)

	mockService.AssertExpectations(t)
}

func TestStaffHandler_Delete_ZeroID(t *testing.T) {
	mockService := mocks.NewMockStaffService()
	router := setupStaffRouter(mockService)

	mockService.On("Delete", uint(0), testSchemaName).Return(errors.New("invalid id"))

	req, _ := http.NewRequest("DELETE", "/staff/delete/0", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	mockService.AssertExpectations(t)
}

// ==================== EDGE CASE TESTS ====================

func TestStaffHandler_Create_EmptyBody(t *testing.T) {
	mockService := mocks.NewMockStaffService()
	router := setupStaffRouter(mockService)

	req, _ := http.NewRequest("POST", "/staff/create", bytes.NewBuffer([]byte("{}")))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestStaffHandler_Update_EmptyBody(t *testing.T) {
	mockService := mocks.NewMockStaffService()
	router := setupStaffRouter(mockService)

	req, _ := http.NewRequest("PUT", "/staff/update/1", bytes.NewBuffer([]byte("{}")))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	// Should fail validation due to missing required fields
	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestStaffHandler_Login_EmptyBody(t *testing.T) {
	mockService := mocks.NewMockStaffService()
	router := setupStaffRouter(mockService)

	req, _ := http.NewRequest("POST", "/staff/login", bytes.NewBuffer([]byte("{}")))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestStaffHandler_Create_AdminUser(t *testing.T) {
	mockService := mocks.NewMockStaffService()
	router := setupStaffRouter(mockService)

	createRequest := domain.StaffCreateRequest{
		Username:    "newadmin",
		Password:    "password123",
		PhoneNumber: "0812345678",
		Email:       "newadmin@test.com",
		FirstName:   "New",
		LastName:    "Admin",
		IsAdmin:     true,
	}

	expectedStaff := &domain.Staff{
		Username:    "newadmin",
		StaffCode:   "STF003",
		PhoneNumber: "0812345678",
		Email:       "newadmin@test.com",
		FirstName:   "New",
		LastName:    "Admin",
		IsAdmin:     true,
	}

	mockService.On("Create", mock.AnythingOfType("*domain.StaffCreateRequest"), testSchemaName).Return(expectedStaff, nil)

	body, _ := json.Marshal(createRequest)
	req, _ := http.NewRequest("POST", "/staff/create", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)
	mockService.AssertExpectations(t)
}

func TestStaffHandler_GetByID_NegativeID(t *testing.T) {
	mockService := mocks.NewMockStaffService()
	router := setupStaffRouter(mockService)

	req, _ := http.NewRequest("GET", "/staff/-1", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	// Negative ID should fail parsing
	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestStaffHandler_Delete_NegativeID(t *testing.T) {
	mockService := mocks.NewMockStaffService()
	router := setupStaffRouter(mockService)

	req, _ := http.NewRequest("DELETE", "/staff/delete/-1", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	// Negative ID should fail parsing
	assert.Equal(t, http.StatusBadRequest, resp.Code)
}
