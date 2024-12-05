package repository

import (
	"calmind/model"
	"errors"
	"time"

	"gorm.io/gorm"
)

type OtpRepository interface {
	GenerateOtp(email string, code string, expiresAt time.Time) error
	GetOtpByEmail(email string) (*model.Otp, error)
	DeleteOtpByEmail(email string) error
}

type OtpRepositoryImpl struct {
	DB *gorm.DB
}

func NewOtpRepository(db *gorm.DB) OtpRepository {
	return &OtpRepositoryImpl{DB: db}
}

func (r *OtpRepositoryImpl) GenerateOtp(email string, code string, expiresAt time.Time) error {
	otp := &model.Otp{
		Email:     email,
		Code:      code,
		ExpiresAt: expiresAt,
	}
	return r.DB.Create(otp).Error
}

func (r *OtpRepositoryImpl) GetOtpByEmail(email string) (*model.Otp, error) {
	var otp model.Otp
	err := r.DB.Where("email = ?", email).First(&otp).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &otp, err
}

func (r *OtpRepositoryImpl) DeleteOtpByEmail(email string) error {
	return r.DB.Where("email = ?", email).Delete(&model.Otp{}).Error
}
