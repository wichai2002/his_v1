package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService interface {
	GenerateToken(userID uint, username string, isAdmin bool, schemaName string) (string, error)
	ValidateToken(token string) (*Claims, error)
}

type Claims struct {
	UserID     uint   `json:"user_id"`
	Username   string `json:"username"`
	IsAdmin    bool   `json:"is_admin"`
	TenantID   uint   `json:"tenant_id"`
	SchemaName string `json:"schema_name"` // Tenant schema for multi-tenancy
	jwt.RegisteredClaims
}

type jwtService struct {
	secretKey string
	expiresIn time.Duration
}

func NewJWTService(secretKey string, expiresIn time.Duration) JWTService {
	return &jwtService{
		secretKey: secretKey,
		expiresIn: expiresIn,
	}
}

func (s *jwtService) GenerateToken(userID uint, username string, isAdmin bool, schemaName string) (string, error) {
	claims := &Claims{
		UserID:     userID,
		Username:   username,
		IsAdmin:    isAdmin,
		SchemaName: schemaName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.expiresIn)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secretKey))
}

func (s *jwtService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
