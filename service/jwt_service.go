package service

import (
	"calmind/config"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtCustomClaims struct {
	UserID     int    `json:"user_id"`
	Email      string `json:"email"`
	Role       string `json:"role"`
	IsVerified bool   `json:"is_verified"`
	jwt.RegisteredClaims
}

type JWTService interface {
	GenerateJWT(email string, userID int, role string, isVerified bool) (string, error)
}

type jwtService struct {
	config *config.JWTConfig
}

func NewJWTService(cfg *config.JWTConfig) JWTService {
	return &jwtService{config: cfg}
}
func (s *jwtService) GenerateJWT(email string, id int, role string, isVerified bool) (string, error) {
	if s.config.SecretKey == "" {
		return "", errors.New("secret key is required")
	}

	claims := &JwtCustomClaims{
		UserID:     id,
		Email:      email,
		Role:       role,
		IsVerified: isVerified,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)),
		},
	}
	// Membuat token dengan claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(s.config.SecretKey))

	if err != nil {
		return "", err
	}
	return t, nil
}
