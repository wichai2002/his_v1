package mocks

import (
	"github.com/stretchr/testify/mock"
	"github.com/wichai2002/his_v1/internal/domain"
)

// MockStaffRepository is a mock implementation of domain.StaffRepository
type MockStaffRepository struct {
	mock.Mock
}

func NewMockStaffRepository() *MockStaffRepository {
	return &MockStaffRepository{}
}

func (m *MockStaffRepository) GetAll(schemaName string) ([]domain.Staff, error) {
	args := m.Called(schemaName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Staff), args.Error(1)
}

func (m *MockStaffRepository) GetByID(id uint, schemaName string) (*domain.Staff, error) {
	args := m.Called(id, schemaName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Staff), args.Error(1)
}

func (m *MockStaffRepository) GetByUsername(username string, schemaName string) (*domain.Staff, error) {
	args := m.Called(username, schemaName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Staff), args.Error(1)
}

func (m *MockStaffRepository) Create(staff *domain.Staff, schemaName string) error {
	args := m.Called(staff, schemaName)
	return args.Error(0)
}

func (m *MockStaffRepository) Update(staff *domain.Staff, schemaName string) error {
	args := m.Called(staff, schemaName)
	return args.Error(0)
}

func (m *MockStaffRepository) Delete(id uint, schemaName string) error {
	args := m.Called(id, schemaName)
	return args.Error(0)
}
