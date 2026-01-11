package handler

import (
	"net/http"
	"strconv"

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

func (h *PatientHandler) Search(c *gin.Context) {
	hospitalID, _ := c.Get("hospital_id")

	patient, err := h.patientService.Search(c.Query("query"), hospitalID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "patient not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "success", patient)
}

func (h *PatientHandler) Create(c *gin.Context) {
	var req domain.PatientCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	hospitalID, _ := c.Get("hospital_id")

	patient, err := h.patientService.Create(&req, hospitalID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
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

	patient, err := h.patientService.Update(uint(id), &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
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

	patient, err := h.patientService.PartialUpdate(uint(id), &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
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

	if err := h.patientService.Delete(uint(id)); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "patient deleted successfully", nil)
}
