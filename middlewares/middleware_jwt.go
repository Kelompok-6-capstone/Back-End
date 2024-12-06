package middlewares

import (
	"calmind/config"
	"calmind/service"
	"log"
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

// Fungsi umum untuk memvalidasi token dan mengambil klaim
func (m *JWTMiddleware) validateToken(c echo.Context) (*service.JwtCustomClaims, error) {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "Token tidak ditemukan")
	}

	// Ambil token dari header Authorization
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	claims := &service.JwtCustomClaims{}

	// Parse token JWT dan validasi klaim
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(m.config.SecretKey), nil
	})

	if err != nil {
		log.Printf("Token error: %v", err)
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "Token tidak valid")
	}

	if !token.Valid {
		log.Println("Token tidak valid")
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "Token tidak valid")
	}

	return claims, nil
}

// Middleware untuk Admin
func (m *JWTMiddleware) HandlerAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		claims, err := m.validateToken(c)
		if err != nil {
			return err
		}

		if claims.Role != "admin" {
			return echo.NewHTTPError(http.StatusForbidden, "Akses hanya untuk admin")
		}

		c.Set("admin", claims) // Simpan klaim di context
		return next(c)
	}
}

// Middleware untuk User
func (m *JWTMiddleware) HandlerUser(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		claims, err := m.validateToken(c)
		if err != nil {
			return err
		}

		if claims.Role != "user" {
			return echo.NewHTTPError(http.StatusForbidden, "Akses hanya untuk user")
		}

		if !claims.IsVerified {
			return echo.NewHTTPError(http.StatusForbidden, "Akun belum diverifikasi. Harap verifikasi OTP Anda.")
		}

		c.Set("user", claims) // Simpan klaim di context
		return next(c)
	}
}

// Middleware untuk Dokter
func (m *JWTMiddleware) HandlerDoctor(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		claims, err := m.validateToken(c)
		if err != nil {
			return err
		}

		if claims.Role != "doctor" {
			return echo.NewHTTPError(http.StatusForbidden, "Akses hanya untuk dokter")
		}

		if !claims.IsVerified {
			return echo.NewHTTPError(http.StatusForbidden, "Akun belum diverifikasi. Harap verifikasi OTP Anda.")
		}

		c.Set("doctor", claims) // Simpan klaim di context
		return next(c)
	}
}
