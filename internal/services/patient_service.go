package services

import (
	"fmt"
	"strings"
	"time"

	"github.com/wichai2002/his_v1/internal/domain"
	"gorm.io/gorm"
)

type patientService struct {
	patientRepo   domain.PatientRepository
	tenantService domain.TenantService
}

func NewPatientService(patientRepo domain.PatientRepository, tenantService domain.TenantService) domain.PatientService {
	return &patientService{
		patientRepo:   patientRepo,
		tenantService: tenantService,
	}
}

// parseDateOfBirth parses and validates date of birth
func parseDateOfBirth(dateStr string) (time.Time, error) {
	if dateStr == "" {
		return time.Time{}, fmt.Errorf("%w: date of birth is required", domain.ErrInvalidInput)
	}

	dob, err := time.Parse(domain.DateFormat, dateStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("%w: expected format YYYY-MM-DD", domain.ErrInvalidDateFormat)
	}

	now := time.Now()

	// Validate date is not in the future
	if dob.After(now) {
		return time.Time{}, domain.ErrFutureDateOfBirth
	}

	// Validate date is not too old (e.g., more than 150 years ago)
	maxAge := now.AddDate(-domain.MaxPatientAge, 0, 0)
	if dob.Before(maxAge) {
		return time.Time{}, domain.ErrDateOfBirthTooOld
	}

	return dob, nil
}

// wrapError wraps repository errors with domain-specific errors
func wrapError(err error) error {
	if err == nil {
		return nil
	}

	// Check for GORM's record not found error
	if err == gorm.ErrRecordNotFound {
		return domain.ErrNotFound
	}

	// Check for duplicate key errors (PostgreSQL)
	errStr := err.Error()
	if strings.Contains(errStr, "duplicate key") || strings.Contains(errStr, "unique constraint") {
		return fmt.Errorf("%w: %s", domain.ErrDuplicateEntry, errStr)
	}

	return err
}

func (s *patientService) Search(query string, schemaName string) ([]domain.Patient, error) {
	patients, err := s.patientRepo.Search(query, schemaName)
	if err != nil {
		return nil, wrapError(err)
	}
	return patients, nil
}

func (s *patientService) SearchByID(id uint, schemaName string) (*domain.Patient, error) {
	patient, err := s.patientRepo.SearchByID(id, schemaName)
	if err != nil {
		return nil, wrapError(err)
	}
	return patient, nil
}

// Create a new patient
func (s *patientService) Create(req *domain.PatientCreateRequest, schemaName string) (*domain.Patient, error) {
	// Validate and parse date of birth
	dob, err := parseDateOfBirth(req.DateOfBirth)
	if err != nil {
		return nil, err
	}

	// Generate HN using tenant's hospital code and running number
	// Format: hospitalCode-HNRunning (e.g., "HOSP0001-00000001")
	hn, err := s.tenantService.GenerateHN(schemaName)
	if err != nil {
		return nil, fmt.Errorf("failed to generate HN: %w", err)
	}

	patient := &domain.Patient{
		FirstNameTH:  strings.TrimSpace(req.FirstNameTH),
		LastNameTH:   strings.TrimSpace(req.LastNameTH),
		MiddleNameTH: strings.TrimSpace(req.MiddleNameTH),
		FirstNameEN:  strings.TrimSpace(req.FirstNameEN),
		LastNameEN:   strings.TrimSpace(req.LastNameEN),
		MiddleNameEN: strings.TrimSpace(req.MiddleNameEN),
		DateOfBirth:  dob,
		NickNameTH:   strings.TrimSpace(req.NickNameTH),
		NickNameEN:   strings.TrimSpace(req.NickNameEN),
		PatientHN:    hn,
		NationalID:   strings.TrimSpace(req.NationalID),
		PassportID:   strings.TrimSpace(req.PassportID),
		PhoneNumber:  strings.TrimSpace(req.PhoneNumber),
		Email:        strings.TrimSpace(req.Email),
		Gender:       req.Gender,
		Nationality:  strings.TrimSpace(req.Nationality),
		BloodGrp:     domain.BloodGrp(req.BloodGrp),
	}

	if err := s.patientRepo.Create(patient, schemaName); err != nil {
		return nil, wrapError(err)
	}

	return patient, nil
}

// Update performs a full update (PUT) - replaces all fields
func (s *patientService) Update(id uint, req *domain.PatientUpdateRequest, schemaName string) (*domain.Patient, error) {
	patient, err := s.patientRepo.GetByID(id, schemaName)
	if err != nil {
		return nil, wrapError(err)
	}

	// Validate and parse date of birth
	dob, err := parseDateOfBirth(req.DateOfBirth)
	if err != nil {
		return nil, err
	}

	// Replace all fields (full update)
	patient.FirstNameTH = strings.TrimSpace(req.FirstNameTH)
	patient.LastNameTH = strings.TrimSpace(req.LastNameTH)
	patient.MiddleNameTH = strings.TrimSpace(req.MiddleNameTH)
	patient.FirstNameEN = strings.TrimSpace(req.FirstNameEN)
	patient.LastNameEN = strings.TrimSpace(req.LastNameEN)
	patient.MiddleNameEN = strings.TrimSpace(req.MiddleNameEN)
	patient.DateOfBirth = dob
	patient.NickNameTH = strings.TrimSpace(req.NickNameTH)
	patient.NickNameEN = strings.TrimSpace(req.NickNameEN)
	patient.NationalID = strings.TrimSpace(req.NationalID)
	patient.PassportID = strings.TrimSpace(req.PassportID)
	patient.PhoneNumber = strings.TrimSpace(req.PhoneNumber)
	patient.Email = strings.TrimSpace(req.Email)
	patient.Gender = domain.Gender(req.Gender)
	patient.Nationality = strings.TrimSpace(req.Nationality)
	patient.BloodGrp = domain.BloodGrp(req.BloodGrp)

	if err := s.patientRepo.Update(patient, schemaName); err != nil {
		return nil, wrapError(err)
	}

	return patient, nil
}

// PartialUpdate performs a partial update (PATCH) - only updates provided fields
func (s *patientService) PartialUpdate(id uint, req *domain.PatientPartialUpdateRequest, schemaName string) (*domain.Patient, error) {
	// Convert request to update map (validates date if present)
	updates, err := req.ToMap()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrInvalidDateFormat, err)
	}

	// Skip database call if there's nothing to update
	if len(updates) == 0 {
		patient, err := s.patientRepo.GetByID(id, schemaName)
		if err != nil {
			return nil, wrapError(err)
		}
		return patient, nil
	}

	// Perform the partial update - repository will handle record existence check
	if err := s.patientRepo.PartialUpdate(id, updates, schemaName); err != nil {
		return nil, wrapError(err)
	}

	// Fetch the updated patient
	patient, err := s.patientRepo.GetByID(id, schemaName)
	if err != nil {
		return nil, wrapError(err)
	}

	return patient, nil
}

func (s *patientService) Delete(id uint, schemaName string) error {
	// First check if patient exists
	_, err := s.patientRepo.GetByID(id, schemaName)
	if err != nil {
		return wrapError(err)
	}

	if err := s.patientRepo.Delete(id, schemaName); err != nil {
		return wrapError(err)
	}
	return nil
}
