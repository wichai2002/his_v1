package services

import (
	"time"

	"github.com/wichai2002/his_v1/internal/domain"
)

type patientService struct {
	patientRepo     domain.PatientRepository
	hospitalService domain.HospitalService
}

func NewPatientService(patientRepo domain.PatientRepository, hospitalService domain.HospitalService) domain.PatientService {
	return &patientService{
		patientRepo:     patientRepo,
		hospitalService: hospitalService,
	}
}

func (s *patientService) Search(query string, hospitalID uint) ([]domain.Patient, error) {
	return s.patientRepo.Search(query, hospitalID)
}

func (s *patientService) SearchByID(id uint, hospitalID uint) (*domain.Patient, error) {
	return s.patientRepo.SearchByID(id, hospitalID)
}

// Create a new patient
func (s *patientService) Create(req *domain.PatientCreateRequest, hospitalID uint) (*domain.Patient, error) {
	// Generate HN from hospital
	hn, err := s.hospitalService.GenerateHN(hospitalID)
	if err != nil {
		return nil, err
	}

	var dob time.Time
	if req.DateOfBirth != "" {
		dob, err = time.Parse("2006-01-02", req.DateOfBirth)
		if err != nil {
			return nil, err
		}
	}

	patient := &domain.Patient{
		FirstNameTH:  req.FirstNameTH,
		LastNameTH:   req.LastNameTH,
		MiddleNameTH: req.MiddleNameTH,
		FirstNameEN:  req.FirstNameEN,
		LastNameEN:   req.LastNameEN,
		MiddleNameEN: req.MiddleNameEN,
		DateOfBirth:  dob,
		NickNameTH:   req.NickNameTH,
		NickNameEN:   req.NickNameEN,
		PatientHN:    hn,
		NationalID:   req.NationalID,
		PassportID:   req.PassportID,
		PhoneNumber:  req.PhoneNumber,
		Email:        req.Email,
		Gender:       req.Gender,
		Nationality:  req.Nationality,
		BloodGrp:     domain.BloodGrp(req.BloodGrp),
		HospitalID:   hospitalID,
	}

	if err := s.patientRepo.Create(patient); err != nil {
		return nil, err
	}

	return patient, nil
}

// Update performs a full update (PUT) - replaces all fields
func (s *patientService) Update(id uint, req *domain.PatientUpdateRequest) (*domain.Patient, error) {
	patient, err := s.patientRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Parse date of birth
	dob, err := time.Parse("2006-01-02", req.DateOfBirth)
	if err != nil {
		return nil, err
	}

	// Replace all fields (full update)
	patient.FirstNameTH = req.FirstNameTH
	patient.LastNameTH = req.LastNameTH
	patient.MiddleNameTH = req.MiddleNameTH
	patient.FirstNameEN = req.FirstNameEN
	patient.LastNameEN = req.LastNameEN
	patient.MiddleNameEN = req.MiddleNameEN
	patient.DateOfBirth = dob
	patient.NickNameTH = req.NickNameTH
	patient.NickNameEN = req.NickNameEN
	patient.NationalID = req.NationalID
	patient.PassportID = req.PassportID
	patient.PhoneNumber = req.PhoneNumber
	patient.Email = req.Email
	patient.Gender = domain.Gender(req.Gender)
	patient.Nationality = req.Nationality
	patient.BloodGrp = domain.BloodGrp(req.BloodGrp)

	if err := s.patientRepo.Update(patient); err != nil {
		return nil, err
	}

	return patient, nil
}

// PartialUpdate performs a partial update (PATCH) - only updates provided fields
func (s *patientService) PartialUpdate(id uint, req *domain.PatientPartialUpdateRequest) (*domain.Patient, error) {
	if _, err := s.patientRepo.GetByID(id); err != nil {
		return nil, err
	}

	updates, err := req.ToMap()
	if err != nil {
		return nil, err
	}

	if err := s.patientRepo.PartialUpdate(id, updates); err != nil {
		return nil, err
	}

	return s.patientRepo.GetByID(id)
}

func (s *patientService) Delete(id uint) error {
	return s.patientRepo.Delete(id)
}
