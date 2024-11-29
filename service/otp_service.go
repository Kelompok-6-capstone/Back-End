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
	b := make([]byte, 5)
	rand.Read(b)
	return base32.StdEncoding.EncodeToString(b)
}

func (s *OtpServiceImpl) IsOtpExpired(expiry time.Time) bool {
	return time.Now().After(expiry)
}
