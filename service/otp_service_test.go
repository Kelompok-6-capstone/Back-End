package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateOtp(t *testing.T) {
	otpService := NewOtpService()

	// Generate OTP
	otp := otpService.GenerateOtp()

	// Validasi panjang OTP
	assert.Len(t, otp, 4, "OTP harus terdiri dari 4 karakter")

	// Validasi bahwa OTP terdiri dari karakter yang valid (base32)
	for _, char := range otp {
		assert.True(t, (char >= 'A' && char <= 'Z') || (char >= '2' && char <= '7'), "OTP hanya boleh berisi karakter base32")
	}
}
