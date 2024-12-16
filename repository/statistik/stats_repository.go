package repository

import (
	"calmind/model"

	"gorm.io/gorm"
)

type StatsRepo interface {
	GetTotalUsers() (int64, error)
	GetTotalDoctors() (int64, error)
	GetTotalKonsultasi() (int64, error)
}

type StatsRepoImpl struct {
	DB *gorm.DB
}

func NewStatsRepo(db *gorm.DB) StatsRepo {
	return &StatsRepoImpl{DB: db}
}

// Get total number of users
func (sr *StatsRepoImpl) GetTotalUsers() (int64, error) {
	var count int64
	err := sr.DB.Model(&model.User{}).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// Get total number of doctors
func (sr *StatsRepoImpl) GetTotalDoctors() (int64, error) {
	var count int64
	err := sr.DB.Model(&model.Doctor{}).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (sr *StatsRepoImpl) GetTotalKonsultasi() (int64, error) {
	var count int64
	err := sr.DB.Model(&model.Consultation{}).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
