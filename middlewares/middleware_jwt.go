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

// Validate Admin Token
func (m *JWTMiddleware) HandlerAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return helper.JSONErrorResponse(c, http.StatusUnauthorized, "Token tidak ditemukan")
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
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

		if claims.Role != "admin" {
			return helper.JSONErrorResponse(c, http.StatusForbidden, "access forbidden")
		}

		c.Set("admin", claims)
		return next(c)
	}
}

// Validate User Token
func (m *JWTMiddleware) HandlerUser(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return helper.JSONErrorResponse(c, http.StatusUnauthorized, "Token tidak ditemukan")
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
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

// Validate Doctor Token
func (m *JWTMiddleware) HandlerDoctor(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return helper.JSONErrorResponse(c, http.StatusUnauthorized, "Token tidak ditemukan")
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
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
