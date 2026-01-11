package mocks

import (
	"github.com/stretchr/testify/mock"
	"github.com/wichai2002/his_v1/internal/domain"
)

// MockStaffService is a mock implementation of domain.StaffService
type MockStaffService struct {
	mock.Mock
}

func NewMockStaffService() *MockStaffService {
	return &MockStaffService{}
}

func (m *MockStaffService) Login(req *domain.StaffLoginRequest, schemaName string) (*domain.StaffLoginResponse, error) {
	args := m.Called(req, schemaName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.StaffLoginResponse), args.Error(1)
}

func (m *MockStaffService) GetAll(schemaName string) ([]domain.Staff, error) {
	args := m.Called(schemaName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Staff), args.Error(1)
}

func (m *MockStaffService) GetByID(id uint, schemaName string) (*domain.Staff, error) {
	args := m.Called(id, schemaName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Staff), args.Error(1)
}

func (m *MockStaffService) Create(req *domain.StaffCreateRequest, schemaName string) (*domain.Staff, error) {
	args := m.Called(req, schemaName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Staff), args.Error(1)
}

func (m *MockStaffService) Update(id uint, req *domain.StaffUpdateRequest, schemaName string) (*domain.Staff, error) {
	args := m.Called(id, req, schemaName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Staff), args.Error(1)
}

func (m *MockStaffService) Delete(id uint, schemaName string) error {
	args := m.Called(id, schemaName)
	return args.Error(0)
}
