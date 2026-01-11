package services

import (
	"fmt"

	"github.com/wichai2002/his_v1/internal/domain"
)

type hospitalService struct {
	hospitalRepo domain.HospitalRepository
}

func NewHospitalService(hospitalRepo domain.HospitalRepository) domain.HospitalService {
	return &hospitalService{
		hospitalRepo: hospitalRepo,
	}
}

func (s *hospitalService) GetAllPublic() ([]domain.HospitalPublicResponse, error) {
	return s.hospitalRepo.GetAllPublic()
}

func (s *hospitalService) GetAll() ([]domain.Hospital, error) {
	return s.hospitalRepo.GetAll()
}

func (s *hospitalService) GetByID(id uint) (*domain.Hospital, error) {
	return s.hospitalRepo.GetByID(id)
}

func (s *hospitalService) Create(hospital *domain.Hospital) error {
	return s.hospitalRepo.Create(hospital)
}

func (s *hospitalService) Update(hospital *domain.Hospital) error {
	return s.hospitalRepo.Update(hospital)
}

func (s *hospitalService) Delete(id uint) error {
	return s.hospitalRepo.Delete(id)
}

// GenerateHN generates a new HN for a patient
func (s *hospitalService) GenerateHN(hospitalID uint) (string, error) {
	hospital, err := s.hospitalRepo.GetByID(hospitalID)
	if err != nil {
		return "", err
	}

	runningNumber, err := s.hospitalRepo.IncrementHNRunningNumber(hospitalID)
	if err != nil {
		return "", err
	}

	// Format: HospitalCode-RunningNumber (e.g., HOS001-000001)
	hn := fmt.Sprintf("%s-%06d", hospital.HospitalCode, runningNumber)
	return hn, nil
}

// GenerateStaffCode generates a new staff code
func (s *hospitalService) GenerateStaffCode(hospitalID uint) (string, error) {
	hospital, err := s.hospitalRepo.GetByID(hospitalID)
	if err != nil {
		return "", err
	}

	runningNumber, err := s.hospitalRepo.IncrementStaffRunningNumber(hospitalID)
	if err != nil {
		return "", err
	}

	// Format: STAFF-HospitalCode-RunningNumber (e.g., STAFF-HOS001-000001)
	staffCode := fmt.Sprintf("STAFF-%s-%06d", hospital.HospitalCode, runningNumber)
	return staffCode, nil
}
