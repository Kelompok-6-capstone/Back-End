package repository

import (
	"calmind/model"
	"gorm.io/gorm"
)

type ConsultationRepository interface {
	CreateConsultation(consultation *model.Consultation) error
	FindByDoctorID(doctorID int, consultations *[]model.Consultation) error
	FindByConsultationID(consultationID string, consultation *model.Consultation) error
	UpdateRecommendation(consultationID int, recommendation string) error
}

type ConsultationRepositoryImpl struct {
	DB *gorm.DB
}

func NewConsultationRepository(db *gorm.DB) ConsultationRepository {
	return &ConsultationRepositoryImpl{DB: db}
}

// Membuat konsultasi
func (r *ConsultationRepositoryImpl) CreateConsultation(consultation *model.Consultation) error {
	if consultation == nil {
		return gorm.ErrInvalidData
	}
	return r.DB.Create(consultation).Error
}

// Mendapatkan daftar konsultasi berdasarkan doctorID
func (r *ConsultationRepositoryImpl) FindByDoctorID(doctorID int, consultations *[]model.Consultation) error {
	if doctorID <= 0 {
		return gorm.ErrRecordNotFound
	}
	return r.DB.Preload("User").Where("doctor_id = ?", doctorID).Find(consultations).Error
}

// Mendapatkan detail konsultasi berdasarkan consultationID
func (r *ConsultationRepositoryImpl) FindByConsultationID(consultationID string, consultation *model.Consultation) error {
    return r.DB.Preload("User").Where("id = ?", consultationID).First(consultation).Error
}


// Memperbarui rekomendasi konsultasi
func (r *ConsultationRepositoryImpl) UpdateRecommendation(consultationID int, recommendation string) error {
    return r.DB.Model(&model.Consultation{}).Where("id = ?", consultationID).Update("rekomendasi", recommendation).Error
}

