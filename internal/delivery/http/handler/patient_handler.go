package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/wichai2002/his_v1/internal/delivery/http/middleware"
	"github.com/wichai2002/his_v1/internal/domain"
	"github.com/wichai2002/his_v1/pkg/utils"

	"github.com/gin-gonic/gin"
)

type PatientHandler struct {
	patientService domain.PatientService
}

func NewPatientHandler(patientService domain.PatientService) *PatientHandler {
	return &PatientHandler{
		patientService: patientService,
	}
}

// handleServiceError maps domain errors to appropriate HTTP status codes
func (h *PatientHandler) handleServiceError(c *gin.Context, err error, resourceName string) {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		utils.ErrorResponse(c, http.StatusNotFound, resourceName+" not found")
	case errors.Is(err, domain.ErrDuplicateEntry):
		utils.ErrorResponse(c, http.StatusConflict, "duplicate "+resourceName+" entry")
	case errors.Is(err, domain.ErrInvalidInput):
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
	case errors.Is(err, domain.ErrInvalidDateFormat):
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid date format")
	case errors.Is(err, domain.ErrFutureDateOfBirth):
		utils.ErrorResponse(c, http.StatusBadRequest, "date of birth cannot be in the future")
	case errors.Is(err, domain.ErrDateOfBirthTooOld):
		utils.ErrorResponse(c, http.StatusBadRequest, "date of birth is too old")
	case errors.Is(err, domain.ErrTenantRequired):
		utils.ErrorResponse(c, http.StatusBadRequest, "tenant context required")
	case errors.Is(err, domain.ErrInvalidSchemaName):
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid tenant schema")
	default:
		// Log the actual error for debugging (in production, use proper logging)
		// log.Printf("Internal error: %v", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "internal server error")
	}
}

func (h *PatientHandler) Search(c *gin.Context) {
	schemaName := middleware.GetTenantSchema(c)

	patients, err := h.patientService.Search(c.Query("query"), schemaName)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, "patient not found")
			return
		}
		h.handleServiceError(c, err, "patient")
		return
	}

	// Return empty array instead of error when no results found
	if len(patients) == 0 {
		utils.SuccessResponse(c, http.StatusOK, "success", []domain.Patient{})
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "success", patients)
}

func (h *PatientHandler) Create(c *gin.Context) {
	var req domain.PatientCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	schemaName := middleware.GetTenantSchema(c)

	patient, err := h.patientService.Create(&req, schemaName)
	if err != nil {
		h.handleServiceError(c, err, "patient")
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "patient created successfully", patient)
}

// Update handles PUT requests for full patient update
func (h *PatientHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid id")
		return
	}

	var req domain.PatientUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	schemaName := middleware.GetTenantSchema(c)

	patient, err := h.patientService.Update(uint(id), &req, schemaName)
	if err != nil {
		h.handleServiceError(c, err, "patient")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "patient updated successfully", patient)
}

// PartialUpdate handles PATCH requests for partial patient update
func (h *PatientHandler) PartialUpdate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid id")
		return
	}

	var req domain.PatientPartialUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	schemaName := middleware.GetTenantSchema(c)

	patient, err := h.patientService.PartialUpdate(uint(id), &req, schemaName)
	if err != nil {
		h.handleServiceError(c, err, "patient")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "patient updated successfully", patient)
}

func (h *PatientHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid id")
		return
	}

	schemaName := middleware.GetTenantSchema(c)

	if err := h.patientService.Delete(uint(id), schemaName); err != nil {
		h.handleServiceError(c, err, "patient")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "patient deleted successfully", nil)
}
