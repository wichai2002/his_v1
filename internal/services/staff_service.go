package services

import (
	"errors"

	"github.com/wichai2002/his_v1/internal/domain"
	"github.com/wichai2002/his_v1/pkg/jwt"

	"golang.org/x/crypto/bcrypt"
)

type staffService struct {
	staffRepo       domain.StaffRepository
	jwtService      jwt.JWTService
	hospitalService domain.HospitalService
}

func NewStaffService(staffRepo domain.StaffRepository, jwtService jwt.JWTService, hospitalService domain.HospitalService) domain.StaffService {
	return &staffService{
		staffRepo:       staffRepo,
		jwtService:      jwtService,
		hospitalService: hospitalService,
	}
}

// login staff and return token and staff information
func (s *staffService) Login(req *domain.StaffLoginRequest) (*domain.StaffLoginResponse, error) {
	staff, err := s.staffRepo.GetByUsernameAndHospitalID(req.Username, req.HospitalID)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(staff.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	token, err := s.jwtService.GenerateToken(staff.ID, staff.Username, staff.IsAdmin, staff.HospitalID)
	if err != nil {
		return nil, errors.New("failed to generate token")
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
			HospitalID:  staff.HospitalID,
			IsAdmin:     staff.IsAdmin,
		},
	}, nil
}

func (s *staffService) GetAll() ([]domain.Staff, error) {
	return s.staffRepo.GetAll()
}

func (s *staffService) GetByID(id uint) (*domain.Staff, error) {
	return s.staffRepo.GetByID(id)
}

// Create creates a new staff member
func (s *staffService) Create(req *domain.StaffCreateRequest) (*domain.Staff, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Generate staff code from hospital
	// format: STAFF-HospitalCode-RunningNumber
	staffCode, err := s.hospitalService.GenerateStaffCode(req.HospitalID)
	if err != nil {
		return nil, errors.New("failed to generate staff code")
	}

	staff := &domain.Staff{
		Username:    req.Username,
		Password:    string(hashedPassword),
		StaffCode:   staffCode,
		PhoneNumber: req.PhoneNumber,
		Email:       req.Email,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		HospitalID:  req.HospitalID,
		IsAdmin:     req.IsAdmin,
	}

	if err := s.staffRepo.Create(staff); err != nil {
		return nil, err
	}

	return staff, nil
}

func (s *staffService) Update(id uint, req *domain.StaffUpdateRequest) (*domain.Staff, error) {
	staff, err := s.staffRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if req.StaffCode != "" {
		staff.StaffCode = req.StaffCode
	}
	if req.PhoneNumber != "" {
		staff.PhoneNumber = req.PhoneNumber
	}
	if req.Email != "" {
		staff.Email = req.Email
	}
	if req.FirstName != "" {
		staff.FirstName = req.FirstName
	}
	if req.LastName != "" {
		staff.LastName = req.LastName
	}
	if req.HospitalID != 0 {
		staff.HospitalID = req.HospitalID
	}
	staff.IsAdmin = req.IsAdmin

	if err := s.staffRepo.Update(staff); err != nil {
		return nil, err
	}

	return staff, nil
}

func (s *staffService) Delete(id uint) error {
	return s.staffRepo.Delete(id)
}
