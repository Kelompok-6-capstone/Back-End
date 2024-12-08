package repository

import (
	"calmind/model"
	"time"

	"gorm.io/gorm"
)

// ConsultationRepository defines the methods for consultation-related database operations
type ConsultationRepository interface {
	CreateConsultation(consultation *model.Consultation) error
	FindByDoctorID(doctorID int, consultations *[]model.Consultation) error
	FindByConsultationID(consultationID int, consultation *model.Consultation) error
	UpdateRecommendation(consultationID int, recommendation string) error
	UpdatePaymentStatus(consultationID int, isPaid bool) error
	UpdateApprovalStatus(consultationID int, isApproved bool) error
	FindUnpaidConsultations(consultations *[]model.Consultation) error
	FindPendingApproval(consultations *[]model.Consultation) error
	ExpireConsultations() error
	FindDoctorByID(doctorID int, doctor *model.Doctor) error
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
	return r.DB.Preload("User").Where("doctor_id = ? AND is_paid = ? AND is_approved = ?", doctorID, true, true).Find(consultations).Error
}

// FindByConsultationID retrieves a consultation by its ID
func (r *ConsultationRepositoryImpl) FindByConsultationID(consultationID int, consultation *model.Consultation) error {
	return r.DB.Preload("User").Preload("Doctor").Where("id = ?", consultationID).First(consultation).Error
}

// UpdateRecommendation updates the recommendation for a consultation
func (r *ConsultationRepositoryImpl) UpdateRecommendation(consultationID int, recommendation string) error {
	return r.DB.Model(&model.Consultation{}).Where("id = ?", consultationID).Update("rekomendasi", recommendation).Error
}

// UpdatePaymentStatus updates the payment status of a consultation
func (r *ConsultationRepositoryImpl) UpdatePaymentStatus(consultationID int, isPaid bool) error {
	return r.DB.Model(&model.Consultation{}).Where("id = ?", consultationID).Update("is_paid", isPaid).Error
}

// UpdateApprovalStatus updates the approval status of a consultation by admin
func (r *ConsultationRepositoryImpl) UpdateApprovalStatus(consultationID int, isApproved bool) error {
	return r.DB.Model(&model.Consultation{}).Where("id = ?", consultationID).Update("is_approved", isApproved).Error
}

// FindUnpaidConsultations retrieves consultations with unpaid status
func (r *ConsultationRepositoryImpl) FindUnpaidConsultations(consultations *[]model.Consultation) error {
	return r.DB.Preload("User").Preload("Doctor").Where("is_paid = ?", false).Find(consultations).Error
}

// FindPendingApproval retrieves consultations with paid status but not approved by admin
func (r *ConsultationRepositoryImpl) FindPendingApproval(consultations *[]model.Consultation) error {
	return r.DB.Preload("User").Preload("Doctor").Where("is_paid = ? AND is_approved = ?", true, false).Find(consultations).Error
}

// ExpireConsultations marks consultations as expired if they exceed their duration
func (r *ConsultationRepositoryImpl) ExpireConsultations() error {
	now := time.Now()
	return r.DB.Model(&model.Consultation{}).
		Where("TIMESTAMPADD(MINUTE, duration, start_time) < ? AND is_paid = ? AND status != ?", now, true, "expired").
		Update("status", "expired").Error
}

// FindDoctorByID retrieves doctor details by ID
func (r *ConsultationRepositoryImpl) FindDoctorByID(doctorID int, doctor *model.Doctor) error {
	if doctorID <= 0 {
		return gorm.ErrRecordNotFound
	}
	return r.DB.First(doctor, doctorID).Error
}
