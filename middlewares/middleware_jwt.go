package middlewares

import (
	"calmind/config"
	"calmind/helper"
	"calmind/service"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type JWTMiddleware struct {
	config *config.JWTConfig
}

func NewJWTMiddleware(cfg *config.JWTConfig) *JWTMiddleware {
	return &JWTMiddleware{config: cfg}
}

// Fungsi untuk memvalidasi token dari header Authorization
func (m *JWTMiddleware) validateToken(c echo.Context) (*service.JwtCustomClaims, error) {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "Token tidak ditemukan")
	}

	// Ambil token dari Authorization header
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	claims := &service.JwtCustomClaims{}

	// Parse dan validasi token JWT
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(m.config.SecretKey), nil
	})

	if err != nil {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "Token tidak valid: "+err.Error())
	}

	if !token.Valid {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "Token tidak valid")
	}

	return claims, nil
}

// Middleware untuk Admin
func (m *JWTMiddleware) HandlerAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		claims, err := m.validateToken(c)
		if err != nil {
			return helper.JSONErrorResponse(c, http.StatusUnauthorized, err.Error())
		}

		if claims.Role != "admin" {
			return helper.JSONErrorResponse(c, http.StatusForbidden, "Akses hanya untuk admin")
		}

		c.Set("admin", claims)
		return next(c)
	}
}

// Middleware untuk User
func (m *JWTMiddleware) HandlerUser(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		claims, err := m.validateToken(c)
		if err != nil {
			return helper.JSONErrorResponse(c, http.StatusUnauthorized, err.Error())
		}

		if !claims.IsVerified {
			return helper.JSONErrorResponse(c, http.StatusForbidden, "Akun belum diverifikasi. Silakan verifikasi OTP Anda.")
		}

		if claims.Role != "user" {
			return helper.JSONErrorResponse(c, http.StatusForbidden, "Akses hanya untuk user")
		}

		c.Set("user", claims)
		return next(c)
	}
}

// Middleware untuk Dokter
func (m *JWTMiddleware) HandlerDoctor(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		claims, err := m.validateToken(c)
		if err != nil {
			return helper.JSONErrorResponse(c, http.StatusUnauthorized, err.Error())
		}

		if !claims.IsVerified {
			return helper.JSONErrorResponse(c, http.StatusForbidden, "Akun belum diverifikasi. Silakan verifikasi OTP Anda.")
		}

		if claims.Role != "doctor" {
			return helper.JSONErrorResponse(c, http.StatusForbidden, "Akses hanya untuk dokter")
		}

		c.Set("doctor", claims)
		return next(c)
	}
}
