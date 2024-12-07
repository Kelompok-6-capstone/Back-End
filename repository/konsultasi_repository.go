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
	UpdateStatus(consultationID int, status string) error
	UpdatePaymentStatus(consultationID int, paymentStatus string) error
	UpdateMessage(consultationID int, message string) error
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

func (r *ConsultationRepositoryImpl) FindByDoctorID(doctorID int, consultations *[]model.Consultation) error {
	return r.DB.Preload("User").Where("doctor_id = ?", doctorID).Find(consultations).Error
}

func (r *ConsultationRepositoryImpl) FindByConsultationID(consultationID string, consultation *model.Consultation) error {
	return r.DB.Preload("User").Preload("Doctor").Where("id = ?", consultationID).First(consultation).Error
}

func (r *ConsultationRepositoryImpl) UpdateRecommendation(consultationID int, recommendation string) error {
	return r.DB.Model(&model.Consultation{}).Where("id = ?", consultationID).Update("recommendation", recommendation).Error
}

func (r *ConsultationRepositoryImpl) UpdateStatus(consultationID int, status string) error {
	return r.DB.Model(&model.Consultation{}).Where("id = ?", consultationID).Update("status", status).Error
}

func (r *ConsultationRepositoryImpl) UpdatePaymentStatus(consultationID int, paymentStatus string) error {
	return r.DB.Model(&model.Consultation{}).Where("id = ?", consultationID).Update("payment_status", paymentStatus).Error
}

func (r *ConsultationRepositoryImpl) UpdateMessage(consultationID int, message string) error {
	return r.DB.Model(&model.Consultation{}).Where("id = ?", consultationID).Update("message", message).Error
}
