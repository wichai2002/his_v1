package mocks

import (
	"github.com/stretchr/testify/mock"
	"github.com/wichai2002/his_v1/internal/domain"
)

// MockPatientService is a mock implementation of domain.PatientService
type MockPatientService struct {
	mock.Mock
}

func NewMockPatientService() *MockPatientService {
	return &MockPatientService{}
}

func (m *MockPatientService) Search(query string, schemaName string) ([]domain.Patient, error) {
	args := m.Called(query, schemaName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Patient), args.Error(1)
}

func (m *MockPatientService) SearchByID(id uint, schemaName string) (*domain.Patient, error) {
	args := m.Called(id, schemaName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Patient), args.Error(1)
}

func (m *MockPatientService) Create(req *domain.PatientCreateRequest, schemaName string) (*domain.Patient, error) {
	args := m.Called(req, schemaName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Patient), args.Error(1)
}

func (m *MockPatientService) Update(id uint, req *domain.PatientUpdateRequest, schemaName string) (*domain.Patient, error) {
	args := m.Called(id, req, schemaName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Patient), args.Error(1)
}

func (m *MockPatientService) PartialUpdate(id uint, req *domain.PatientPartialUpdateRequest, schemaName string) (*domain.Patient, error) {
	args := m.Called(id, req, schemaName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Patient), args.Error(1)
}

func (m *MockPatientService) Delete(id uint, schemaName string) error {
	args := m.Called(id, schemaName)
	return args.Error(0)
}
