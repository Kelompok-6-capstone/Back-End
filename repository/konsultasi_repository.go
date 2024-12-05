package repository

import (
	"calmind/model"
	"gorm.io/gorm"
)

// ConsultationRepository defines the methods for consultation-related database operations
type ConsultationRepository interface {
	CreateConsultation(consultation *model.Consultation) error
	FindByDoctorID(doctorID int, consultations *[]model.Consultation) error
	FindByConsultationID(consultationID string, consultation *model.Consultation) error
	UpdateRecommendation(consultationID int, recommendation string) error
}

// ConsultationRepositoryImpl is the implementation of ConsultationRepository
type ConsultationRepositoryImpl struct {
	DB *gorm.DB
}

// NewConsultationRepository creates a new instance of ConsultationRepository
func NewConsultationRepository(db *gorm.DB) ConsultationRepository {
	return &ConsultationRepositoryImpl{DB: db}
}

// CreateConsultation creates a new consultation record
func (r *ConsultationRepositoryImpl) CreateConsultation(consultation *model.Consultation) error {
	if consultation == nil {
		return gorm.ErrInvalidData
	}
	return r.DB.Create(consultation).Error
}

// FindByDoctorID retrieves consultations by doctor ID
func (r *ConsultationRepositoryImpl) FindByDoctorID(doctorID int, consultations *[]model.Consultation) error {
	if doctorID <= 0 {
		return gorm.ErrRecordNotFound
	}
	return r.DB.Preload("User").Where("doctor_id = ?", doctorID).Find(consultations).Error
}

// FindByConsultationID retrieves a consultation by its ID
func (r *ConsultationRepositoryImpl) FindByConsultationID(consultationID string, consultation *model.Consultation) error {
	return r.DB.Preload("User").Where("id = ?", consultationID).First(consultation).Error
}

// UpdateRecommendation updates the recommendation for a consultation
func (r *ConsultationRepositoryImpl) UpdateRecommendation(consultationID int, recommendation string) error {
	return r.DB.Model(&model.Consultation{}).Where("id = ?", consultationID).Update("rekomendasi", recommendation).Error
}
