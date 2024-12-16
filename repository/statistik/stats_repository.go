package repository

import (
	"gorm.io/gorm"
)

type StatsRepository interface {
	GetTotalUsers() (int64, error)
	GetTotalDoctors() (int64, error)
	GetTotalConsultations() (int64, error)
	GetTotalConsultationsByPaymentStatus(paymentStatus string) (int64, error)
}

type StatsRepositoryImpl struct {
	DB *gorm.DB
}

func NewStatsRepository(db *gorm.DB) StatsRepository {
	return &StatsRepositoryImpl{DB: db}
}

func (r *StatsRepositoryImpl) GetTotalUsers() (int64, error) {
	var totalUsers int64
	err := r.DB.Table("users").Count(&totalUsers).Error
	return totalUsers, err
}

func (r *StatsRepositoryImpl) GetTotalDoctors() (int64, error) {
	var totalDoctors int64
	err := r.DB.Table("doctors").Count(&totalDoctors).Error
	return totalDoctors, err
}

func (r *StatsRepositoryImpl) GetTotalConsultations() (int64, error) {
	var totalConsultations int64
	err := r.DB.Table("consultations").Count(&totalConsultations).Error
	return totalConsultations, err
}

func (r *StatsRepositoryImpl) GetTotalConsultationsByPaymentStatus(paymentStatus string) (int64, error) {
	var count int64
	err := r.DB.Table("consultations").
		Where("payment_status = ?", paymentStatus).
		Count(&count).Error
	return count, err
}
