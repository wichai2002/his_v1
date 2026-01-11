package mocks

import (
	"github.com/stretchr/testify/mock"
	"github.com/wichai2002/his_v1/internal/domain"
)

// MockPatientRepository is a mock implementation of domain.PatientRepository
type MockPatientRepository struct {
	mock.Mock
}

func NewMockPatientRepository() *MockPatientRepository {
	return &MockPatientRepository{}
}

func (m *MockPatientRepository) GetAll(schemaName string) ([]domain.Patient, error) {
	args := m.Called(schemaName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Patient), args.Error(1)
}

func (m *MockPatientRepository) GetByID(id uint, schemaName string) (*domain.Patient, error) {
	args := m.Called(id, schemaName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Patient), args.Error(1)
}

func (m *MockPatientRepository) Search(query string, schemaName string) ([]domain.Patient, error) {
	args := m.Called(query, schemaName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Patient), args.Error(1)
}

func (m *MockPatientRepository) SearchByID(id uint, schemaName string) (*domain.Patient, error) {
	args := m.Called(id, schemaName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Patient), args.Error(1)
}

func (m *MockPatientRepository) Create(patient *domain.Patient, schemaName string) error {
	args := m.Called(patient, schemaName)
	return args.Error(0)
}

func (m *MockPatientRepository) Update(patient *domain.Patient, schemaName string) error {
	args := m.Called(patient, schemaName)
	return args.Error(0)
}

func (m *MockPatientRepository) PartialUpdate(id uint, updates map[string]interface{}, schemaName string) error {
	args := m.Called(id, updates, schemaName)
	return args.Error(0)
}

func (m *MockPatientRepository) Delete(id uint, schemaName string) error {
	args := m.Called(id, schemaName)
	return args.Error(0)
}
