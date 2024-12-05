package config

import "os"

type JWTConfig struct {
	SecretKey string
}

func NewJWTConfig() *JWTConfig {
	return &JWTConfig{
		SecretKey: os.Getenv("JWT_SECRET_KEY"), // Ganti dengan secret key yang aman
	}
}
