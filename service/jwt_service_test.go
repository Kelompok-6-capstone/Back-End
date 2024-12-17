package service

import (
	"calmind/config"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestGenerateJWT_Success(t *testing.T) {
	// Arrange
	cfg := &config.JWTConfig{
		SecretKey: "test-secret-key",
	}
	jwtService := NewJWTService(cfg)
	email := "user@example.com"
	userID := 1
	role := "user"
	isVerified := true

	// Act
	token, err := jwtService.GenerateJWT(email, userID, role, isVerified)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Parse token untuk memverifikasi hasil
	parsedToken, _ := jwt.ParseWithClaims(token, &JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.SecretKey), nil
	})

	assert.True(t, parsedToken.Valid)
	claims, ok := parsedToken.Claims.(*JwtCustomClaims)
	assert.True(t, ok)
	assert.Equal(t, email, claims.Email)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, role, claims.Role)
	assert.Equal(t, isVerified, claims.IsVerified)
	assert.WithinDuration(t, time.Now().Add(72*time.Hour), claims.ExpiresAt.Time, 1*time.Second)
}

func TestGenerateJWT_EmptySecretKey(t *testing.T) {
	// Arrange
	cfg := &config.JWTConfig{
		SecretKey: "", // Secret key kosong
	}
	jwtService := NewJWTService(cfg)
	email := "user@example.com"
	userID := 1
	role := "user"
	isVerified := true

	// Act
	token, err := jwtService.GenerateJWT(email, userID, role, isVerified)

	// Assert
	assert.Error(t, err)
	assert.EqualError(t, err, "secret key is required")
	assert.Empty(t, token)
}

func TestGenerateJWT_SignedStringError(t *testing.T) {
	// Arrange
	cfg := &config.JWTConfig{
		SecretKey: "test-secret-key",
	}
	jwtService := &jwtService{config: cfg}

	// Buat token dummy untuk memanipulasi `SignedString` error
	claims := &JwtCustomClaims{
		UserID: 1,
		Email:  "user@example.com",
		Role:   "user",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Act
	// Override `SignedString` untuk menghasilkan error
	jwtSigningMethod := token.Method.(*jwt.SigningMethodHMAC)
	jwtSigningMethod.Hash = 0 // Force error by using invalid hash
	tok, err := jwtService.GenerateJWT("user@example.com", 1, "user", true)

	// Assert
	assert.Error(t, err)
	assert.Empty(t, tok)
}
