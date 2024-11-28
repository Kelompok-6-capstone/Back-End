package middleware

import (
	"calmind/config"
	"calmind/helper"
	"calmind/service"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type JWTMiddleware struct {
	config *config.JWTConfig
}

func NewJWTMiddleware(cfg *config.JWTConfig) *JWTMiddleware {
	return &JWTMiddleware{config: cfg}
}

func (m *JWTMiddleware) HandlerUser(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Ambil token dari cookie
		h, err := c.Cookie("token_user")
		if err != nil {
			return helper.JSONErrorResponse(c, http.StatusUnauthorized, "gagal login token tidak ditemukan")
		}

		tokenString := h.Value
		claims := &service.JwtCustomClaims{}

		// Parse dan validasi token
		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(m.config.SecretKey), nil
		})

		if err != nil {
			return helper.JSONErrorResponse(c, http.StatusUnauthorized, "invalid token")
		}

		// Pastikan token valid
		if !token.Valid {
			return helper.JSONErrorResponse(c, http.StatusUnauthorized, "invalid token")
		}

		// Validasi role user
		if claims.Role != "user" {
			return helper.JSONErrorResponse(c, http.StatusForbidden, "access forbidden")
		}

		// Simpan klaim di context untuk digunakan di handler berikutnya
		c.Set("user", claims)

		return next(c)
	}
}

func (m *JWTMiddleware) HandlerAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Ambil token dari cookie
		h, err := c.Cookie("token_admin")
		if err != nil {
			return helper.JSONErrorResponse(c, http.StatusUnauthorized, "gagal login token tidak ditemukan")
		}

		tokenString := h.Value
		claims := &service.JwtCustomClaims{}

		// Parse dan validasi token
		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(m.config.SecretKey), nil
		})

		if err != nil {
			return helper.JSONErrorResponse(c, http.StatusUnauthorized, "invalid token")
		}

		// Pastikan token valid
		if !token.Valid {
			return helper.JSONErrorResponse(c, http.StatusUnauthorized, "invalid token")
		}

		// Validasi role admin
		if claims.Role != "admin" {
			return helper.JSONErrorResponse(c, http.StatusForbidden, "access forbidden")
		}

		// Simpan klaim di context untuk digunakan di handler berikutnya
		c.Set("admin", claims)

		return next(c)
	}
}
