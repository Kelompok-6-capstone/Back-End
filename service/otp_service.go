package service

import (
	"crypto/rand"
	"encoding/base32"
	"time"
)

type OtpService interface {
	GenerateOtp() string
	IsOtpExpired(expiry time.Time) bool
}

type OtpServiceImpl struct{}

func NewOtpService() OtpService {
	return &OtpServiceImpl{}
}

func (s *OtpServiceImpl) GenerateOtp() string {
	b := make([]byte, 3)
	rand.Read(b)
	encoded := base32.StdEncoding.EncodeToString(b)
	return encoded[:4]
}

func (s *OtpServiceImpl) IsOtpExpired(expiry time.Time) bool {
	return time.Now().After(expiry)
}
