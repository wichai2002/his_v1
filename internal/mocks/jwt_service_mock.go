package mocks

import (
	"github.com/stretchr/testify/mock"
	"github.com/wichai2002/his_v1/pkg/jwt"
)

// MockJWTService is a mock implementation of jwt.JWTService
type MockJWTService struct {
	mock.Mock
}

func NewMockJWTService() *MockJWTService {
	return &MockJWTService{}
}

func (m *MockJWTService) GenerateToken(userID uint, username string, isAdmin bool, schemaName string) (string, error) {
	args := m.Called(userID, username, isAdmin, schemaName)
	return args.String(0), args.Error(1)
}

func (m *MockJWTService) ValidateToken(tokenString string) (*jwt.Claims, error) {
	args := m.Called(tokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*jwt.Claims), args.Error(1)
}
