package handler

import (
	"net/http"
	"strconv"

	"github.com/wichai2002/his_v1/internal/delivery/http/middleware"
	"github.com/wichai2002/his_v1/internal/domain"
	"github.com/wichai2002/his_v1/pkg/utils"

	"github.com/gin-gonic/gin"
)

type StaffHandler struct {
	staffService domain.StaffService
}

func NewStaffHandler(staffService domain.StaffService) *StaffHandler {
	return &StaffHandler{
		staffService: staffService,
	}
}

// Login godoc
// @Summary Staff login
// @Description Authenticate staff and return JWT token
// @Tags staff
// @Accept json
// @Produce json
// @Param request body domain.StaffLoginRequest true "Login credentials"
// @Success 200 {object} domain.StaffLoginResponse
// @Router /staff/login [post]
func (h *StaffHandler) Login(c *gin.Context) {
	var req domain.StaffLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Get tenant schema from context
	schemaName := middleware.GetTenantSchema(c)

	response, err := h.staffService.Login(&req, schemaName)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "login successful", response)
}

// Logout godoc
// @Summary Staff logout
// @Description Logout staff (client should discard the token)
// @Tags staff
// @Security BearerAuth
// @Success 200 {object} utils.Response
// @Router /staff/logout [post]
func (h *StaffHandler) Logout(c *gin.Context) {
	// JWT is stateless, so logout is handled client-side by discarding the token
	utils.SuccessResponse(c, http.StatusOK, "logout successful", nil)
}

// GetAll godoc
// @Summary Get all staff
// @Description Get all staff members
// @Tags staff
// @Security BearerAuth
// @Produce json
// @Success 200 {object} utils.Response
// @Router /staff [get]
func (h *StaffHandler) GetAll(c *gin.Context) {
	schemaName := middleware.GetTenantSchema(c)

	staffs, err := h.staffService.GetAll(schemaName)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "success", staffs)
}

// GetByID godoc
// @Summary Get staff by ID
// @Description Get staff member by ID
// @Tags staff
// @Security BearerAuth
// @Produce json
// @Param id path int true "Staff ID"
// @Success 200 {object} utils.Response
// @Router /staff/{id} [get]
func (h *StaffHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid id")
		return
	}

	schemaName := middleware.GetTenantSchema(c)

	staff, err := h.staffService.GetByID(uint(id), schemaName)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "staff not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "success", staff)
}

// Create godoc
// @Summary Create new staff
// @Description Create new staff member (admin only)
// @Tags staff
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body domain.StaffCreateRequest true "Staff data"
// @Success 201 {object} utils.Response
// @Router /staff/create [post]
func (h *StaffHandler) Create(c *gin.Context) {
	var req domain.StaffCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	schemaName := middleware.GetTenantSchema(c)

	staff, err := h.staffService.Create(&req, schemaName)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "staff created successfully", staff)
}

// Update godoc
// @Summary Update staff
// @Description Update staff member
// @Tags staff
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Staff ID"
// @Param request body domain.StaffUpdateRequest true "Staff data"
// @Success 200 {object} utils.Response
// @Router /staff/update/{id} [put]
func (h *StaffHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid id")
		return
	}

	var req domain.StaffUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	schemaName := middleware.GetTenantSchema(c)

	staff, err := h.staffService.Update(uint(id), &req, schemaName)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "staff updated successfully", staff)
}

// Delete godoc
// @Summary Delete staff
// @Description Delete staff member (admin only)
// @Tags staff
// @Security BearerAuth
// @Param id path int true "Staff ID"
// @Success 200 {object} utils.Response
// @Router /staff/delete/{id} [delete]
func (h *StaffHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid id")
		return
	}

	schemaName := middleware.GetTenantSchema(c)

	if err := h.staffService.Delete(uint(id), schemaName); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "staff deleted successfully", nil)
}
