package middlewares

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

func (m *JWTMiddleware) HandlerAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		h, err := c.Cookie("token_admin")
		if err != nil {
			return helper.JSONErrorResponse(c, http.StatusUnauthorized, "Token admin tidak ditemukan")
		}

		tokenString := h.Value
		claims := &service.JwtCustomClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(m.config.SecretKey), nil
		})

		if err != nil || !token.Valid {
			return helper.JSONErrorResponse(c, http.StatusUnauthorized, "Token admin tidak valid")
		}

		if claims.Role != "admin" {
			return helper.JSONErrorResponse(c, http.StatusForbidden, "Akses hanya untuk admin")
		}

		// Validasi bahwa admin_id tidak kosong
		if claims.UserID == 0 {
			return helper.JSONErrorResponse(c, http.StatusUnauthorized, "Klaim JWT tidak valid atau admin_id kosong")
		}

		// Set klaim ke context
		c.Set("admin", claims)
		return next(c)
	}
}

func (m *JWTMiddleware) HandlerUser(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		h, err := c.Cookie("token_user")
		if err != nil {
			return helper.JSONErrorResponse(c, http.StatusUnauthorized, "gagal login token tidak ditemukan")
		}

		tokenString := h.Value
		claims := &service.JwtCustomClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(m.config.SecretKey), nil
		})

		if err != nil {
			return helper.JSONErrorResponse(c, http.StatusUnauthorized, "invalid token")
		}

		if !token.Valid {
			return helper.JSONErrorResponse(c, http.StatusUnauthorized, "invalid token")
		}

		if !claims.IsVerified {
			return helper.JSONErrorResponse(c, http.StatusForbidden, "access forbidden: account not verified. Please verify your OTP.")
		}

		if claims.Role != "user" {
			return helper.JSONErrorResponse(c, http.StatusForbidden, "access forbidden")
		}

		c.Set("user", claims)

		return next(c)
	}
}

func (m *JWTMiddleware) HandlerDoctor(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		h, err := c.Cookie("token_doctor")
		if err != nil {
			return helper.JSONErrorResponse(c, http.StatusUnauthorized, "gagal login token tidak ditemukan")
		}

		tokenString := h.Value
		claims := &service.JwtCustomClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(m.config.SecretKey), nil
		})

		if err != nil {
			return helper.JSONErrorResponse(c, http.StatusUnauthorized, "invalid token")
		}

		if !token.Valid {
			return helper.JSONErrorResponse(c, http.StatusUnauthorized, "invalid token")
		}

		if !claims.IsVerified {
			return helper.JSONErrorResponse(c, http.StatusForbidden, "access forbidden: account not verified. Please verify your OTP.")
		}

		if claims.Role != "doctor" {
			return helper.JSONErrorResponse(c, http.StatusForbidden, "access forbidden")
		}

		c.Set("doctor", claims)

		return next(c)
	}
}
