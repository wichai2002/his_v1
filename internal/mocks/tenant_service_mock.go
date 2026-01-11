package mocks

import (
	"github.com/stretchr/testify/mock"
	"github.com/wichai2002/his_v1/internal/domain"
)

// MockTenantService is a mock implementation of domain.TenantService
type MockTenantService struct {
	mock.Mock
}

func NewMockTenantService() *MockTenantService {
	return &MockTenantService{}
}

func (m *MockTenantService) GetBySubdomain(subdomain string) (*domain.Tenant, error) {
	args := m.Called(subdomain)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Tenant), args.Error(1)
}

func (m *MockTenantService) GetBySchemaName(schemaName string) (*domain.Tenant, error) {
	args := m.Called(schemaName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Tenant), args.Error(1)
}

func (m *MockTenantService) CreateTenant(req *domain.TenantCreateRequest) (*domain.Tenant, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Tenant), args.Error(1)
}

func (m *MockTenantService) CreateTenantSchema(schemaName string) error {
	args := m.Called(schemaName)
	return args.Error(0)
}

func (m *MockTenantService) MigrateTenantSchema(schemaName string) error {
	args := m.Called(schemaName)
	return args.Error(0)
}

func (m *MockTenantService) SetupTenantWithAdmin(
	tenantCode, name, subdomain, hospitalName, hospitalCode string,
	address *string,
	adminUsername, adminPassword, adminEmail string,
) (*domain.Tenant, error) {
	args := m.Called(tenantCode, name, subdomain, hospitalName, hospitalCode, address, adminUsername, adminPassword, adminEmail)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Tenant), args.Error(1)
}

func (m *MockTenantService) GenerateHN(schemaName string) (string, error) {
	args := m.Called(schemaName)
	return args.String(0), args.Error(1)
}
