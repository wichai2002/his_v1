package handler

import (
	"net/http"
	"strconv"

	"github.com/wichai2002/his_v1/internal/domain"
	"github.com/wichai2002/his_v1/pkg/utils"

	"github.com/gin-gonic/gin"
)

type HospitalHandler struct {
	hospitalService domain.HospitalService
}

func NewHospitalHandler(hospitalService domain.HospitalService) *HospitalHandler {
	return &HospitalHandler{
		hospitalService: hospitalService,
	}
}

func (h *HospitalHandler) GetAllPublic(c *gin.Context) {
	hospitals, err := h.hospitalService.GetAllPublic()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "success", hospitals)
}

// GetAll godoc
// @Summary Get all hospitals
// @Description Get all hospitals
// @Tags hospital
// @Security BearerAuth
// @Produce json
// @Success 200 {object} utils.Response
// @Router /hospital [get]
func (h *HospitalHandler) GetAll(c *gin.Context) {
	hospitals, err := h.hospitalService.GetAll()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "success", hospitals)
}

// GetByID godoc
// @Summary Get hospital by ID
// @Description Get hospital by ID
// @Tags hospital
// @Security BearerAuth
// @Produce json
// @Param id path int true "Hospital ID"
// @Success 200 {object} utils.Response
// @Router /hospital/{id} [get]
func (h *HospitalHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid id")
		return
	}

	hospital, err := h.hospitalService.GetByID(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "hospital not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "success", hospital)
}

// Create godoc
// @Summary Create new hospital
// @Description Create new hospital (admin only)
// @Tags hospital
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body domain.Hospital true "Hospital data"
// @Success 201 {object} utils.Response
// @Router /hospital/create [post]
func (h *HospitalHandler) Create(c *gin.Context) {
	var hospital domain.Hospital
	if err := c.ShouldBindJSON(&hospital); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.hospitalService.Create(&hospital); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "hospital created successfully", hospital)
}

// Update godoc
// @Summary Update hospital
// @Description Update hospital (admin only)
// @Tags hospital
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Hospital ID"
// @Param request body domain.Hospital true "Hospital data"
// @Success 200 {object} utils.Response
// @Router /hospital/update/{id} [put]
func (h *HospitalHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid id")
		return
	}

	var hospital domain.Hospital
	if err := c.ShouldBindJSON(&hospital); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	hospital.ID = uint(id)
	if err := h.hospitalService.Update(&hospital); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "hospital updated successfully", hospital)
}

// Delete godoc
// @Summary Delete hospital
// @Description Delete hospital (admin only)
// @Tags hospital
// @Security BearerAuth
// @Param id path int true "Hospital ID"
// @Success 200 {object} utils.Response
// @Router /hospital/delete/{id} [delete]
func (h *HospitalHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.hospitalService.Delete(uint(id)); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "hospital deleted successfully", nil)
}
