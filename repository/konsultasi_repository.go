package repository

import (
	"calmind/model"

	"gorm.io/gorm"
)

// ConsultationRepository
type ConsultationRepository interface {
	CreateConsultation(consultation *model.Consultation) error
}

type ConsultationRepositoryImpl struct {
	DB *gorm.DB
}

func NewConsultationRepository(db *gorm.DB) ConsultationRepository {
	return &ConsultationRepositoryImpl{DB: db}
}

func (r *ConsultationRepositoryImpl) CreateConsultation(consultation *model.Consultation) error {
	return r.DB.Create(consultation).Error
}
