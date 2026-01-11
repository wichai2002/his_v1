package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/wichai2002/his_v1/internal/domain"
	"github.com/wichai2002/his_v1/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

type staffService struct {
	staffRepo  domain.StaffRepository
	jwtService jwt.JWTService
}

func NewStaffService(
	staffRepo domain.StaffRepository,
	jwtService jwt.JWTService,
) domain.StaffService {
	return &staffService{
		staffRepo:  staffRepo,
		jwtService: jwtService,
	}
}

// Login authenticates a staff member and returns a JWT token
func (s *staffService) Login(req *domain.StaffLoginRequest, schemaName string) (*domain.StaffLoginResponse, error) {
	// Tenant schema provides isolation, just search by username
	staff, err := s.staffRepo.GetByUsername(req.Username, schemaName)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(staff.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Generate JWT token with schema information
	token, err := s.jwtService.GenerateToken(staff.ID, staff.Username, staff.IsAdmin, schemaName)
	if err != nil {
		return nil, err
	}

	return &domain.StaffLoginResponse{
		Token: token,
		Staff: domain.StaffResponse{
			ID:          staff.ID,
			Username:    staff.Username,
			StaffCode:   staff.StaffCode,
			PhoneNumber: staff.PhoneNumber,
			Email:       staff.Email,
			FirstName:   staff.FirstName,
			LastName:    staff.LastName,
			IsAdmin:     staff.IsAdmin,
		},
	}, nil
}

func (s *staffService) GetAll(schemaName string) ([]domain.Staff, error) {
	return s.staffRepo.GetAll(schemaName)
}

func (s *staffService) GetByID(id uint, schemaName string) (*domain.Staff, error) {
	return s.staffRepo.GetByID(id, schemaName)
}

// Create staff
func (s *staffService) Create(req *domain.StaffCreateRequest, schemaName string) (*domain.Staff, error) {
	// Generate staff code using timestamp-based generation
	staffCode := fmt.Sprintf("STF%d", time.Now().UnixNano()%100000000)

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	staff := &domain.Staff{
		Username:    req.Username,
		Password:    string(hashedPassword),
		StaffCode:   staffCode,
		PhoneNumber: req.PhoneNumber,
		Email:       req.Email,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		IsAdmin:     req.IsAdmin,
	}

	if err := s.staffRepo.Create(staff, schemaName); err != nil {
		return nil, err
	}

	return staff, nil
}

func (s *staffService) Update(id uint, req *domain.StaffUpdateRequest, schemaName string) (*domain.Staff, error) {
	staff, err := s.staffRepo.GetByID(id, schemaName)
	if err != nil {
		return nil, err
	}

	// Update fields
	staff.StaffCode = req.StaffCode
	staff.PhoneNumber = req.PhoneNumber
	staff.Email = req.Email
	staff.FirstName = req.FirstName
	staff.LastName = req.LastName
	staff.IsAdmin = req.IsAdmin

	if err := s.staffRepo.Update(staff, schemaName); err != nil {
		return nil, err
	}

	return staff, nil
}

func (s *staffService) Delete(id uint, schemaName string) error {
	return s.staffRepo.Delete(id, schemaName)
}
