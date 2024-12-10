package repository

import (
	"calmind/model"
	"errors"
	"time"

	"gorm.io/gorm"
)

type ConsultationRepository interface {
	CreateConsultation(*model.Consultation) (int, error)
	ApprovePayment(consultationID int) error
	GetConsultationsForDoctor(doctorID int) ([]model.Consultation, error)
	GetConsultationDetails(consultationID, doctorID int) (*model.Consultation, error)
	AddRecommendation(recommendation *model.Rekomendasi) error
	GetAdminViewConsultation(consultationID int) (*model.Consultation, error)
	GetConsultationByID(consultationID int) (*model.Consultation, error)
	UpdateConsultation(consultation *model.Consultation) error
	GetActiveConsultations() ([]model.Consultation, error)
	GetConsultationsWithDoctors(userID int) ([]model.Consultation, error)
	GetPendingConsultations() ([]model.Consultation, error)
	AddIncomeForDoctor(doctorID int, amount float64) error
	AddIncomeForAdmin(amount float64) error
	GetDoctorByID(doctorID int) (*model.Doctor, error)
	GetPendingPayments() ([]model.Consultation, error)
}

type ConsultationRepositoryImpl struct {
	DB *gorm.DB
}

func NewConsultationRepositoryImpl(db *gorm.DB) *ConsultationRepositoryImpl {
	return &ConsultationRepositoryImpl{DB: db}
}

func (r *ConsultationRepositoryImpl) CreateConsultation(consultation *model.Consultation) (int, error) {
	if err := r.DB.Preload("Doctor").Create(consultation).Error; err != nil {
		return 0, err
	}
	return consultation.ID, nil
}

// **2. Approve Payment**
func (r *ConsultationRepositoryImpl) ApprovePayment(consultationID int) error {
	var consultation model.Consultation
	if err := r.DB.First(&consultation, consultationID).Error; err != nil {
		return errors.New("consultation not found")
	}

	if !consultation.IsPaid {
		return errors.New("payment not completed")
	}

	if consultation.IsApproved {
		return errors.New("payment already approved")
	}

	consultation.IsApproved = true
	consultation.Status = "active"
	consultation.StartTime = time.Now()

	return r.DB.Save(&consultation).Error
}

func (r *ConsultationRepositoryImpl) GetConsultationsForDoctor(doctorID int) ([]model.Consultation, error) {
	var consultations []model.Consultation
	if err := r.DB.Preload("User").Preload("Rekomendasi").
		Where("doctor_id = ? AND status = ?", doctorID, "active").
		Find(&consultations).Error; err != nil {
		return nil, err
	}
	return consultations, nil
}

func (r *ConsultationRepositoryImpl) GetConsultationDetails(consultationID, doctorID int) (*model.Consultation, error) {
	var consultation model.Consultation
	if err := r.DB.Preload("User").Preload("Doctor").Preload("Rekomendasi").
		Where("id = ? AND doctor_id = ?", consultationID, doctorID).
		First(&consultation).Error; err != nil {
		return nil, err
	}
	return &consultation, nil
}

func (r *ConsultationRepositoryImpl) AddRecommendation(recommendation *model.Rekomendasi) error {
	return r.DB.Create(recommendation).Error
}

func (r *ConsultationRepositoryImpl) GetAdminViewConsultation(consultationID int) (*model.Consultation, error) {
	var consultation model.Consultation
	if err := r.DB.Preload("User").Preload("Doctor").Preload("Rekomendasi").
		First(&consultation, consultationID).Error; err != nil {
		return nil, err
	}
	return &consultation, nil
}

func (r *ConsultationRepositoryImpl) GetConsultationByID(consultationID int) (*model.Consultation, error) {
	var consultation model.Consultation
	if err := r.DB.Preload("User").Preload("Doctor").Preload("Rekomendasi").
		First(&consultation, consultationID).Error; err != nil {
		return nil, err
	}
	return &consultation, nil
}

func (r *ConsultationRepositoryImpl) UpdateConsultation(consultation *model.Consultation) error {
	return r.DB.Save(consultation).Error
}

func (r *ConsultationRepositoryImpl) GetActiveConsultations() ([]model.Consultation, error) {
	var consultations []model.Consultation
	if err := r.DB.Preload("User").Preload("Doctor").Preload("Rekomendasi").
		Where("status = ?", "active").
		Find(&consultations).Error; err != nil {
		return nil, err
	}
	return consultations, nil
}
func (r *ConsultationRepositoryImpl) GetConsultationsWithDoctors(userID int) ([]model.Consultation, error) {
	var consultations []model.Consultation
	if err := r.DB.Preload("Doctor").Preload("Rekomendasi").
		Where("user_id = ?", userID).
		Find(&consultations).Error; err != nil {
		return nil, err
	}
	return consultations, nil
}

func (r *ConsultationRepositoryImpl) GetPendingConsultations() ([]model.Consultation, error) {
	var consultations []model.Consultation
	if err := r.DB.Preload("User").Preload("Doctor").Preload("Rekomendasi").
		Where("status = ?", "pending").
		Find(&consultations).Error; err != nil {
		return nil, err
	}
	return consultations, nil
}

func (r *ConsultationRepositoryImpl) AddIncomeForDoctor(doctorID int, amount float64) error {
	return r.DB.Model(&model.Doctor{}).
		Where("id = ?", doctorID).
		Update("income", gorm.Expr("income + ?", amount)).
		Error
}

func (r *ConsultationRepositoryImpl) AddIncomeForAdmin(amount float64) error {
	return r.DB.Exec("UPDATE admins SET income = income + ? WHERE id = ?", amount, 1).Error
}

func (r *ConsultationRepositoryImpl) GetDoctorByID(doctorID int) (*model.Doctor, error) {
	var doctor model.Doctor
	if err := r.DB.First(&doctor, doctorID).Error; err != nil {
		return nil, err
	}
	return &doctor, nil
}

func (r *ConsultationRepositoryImpl) GetPendingPayments() ([]model.Consultation, error) {
	var consultations []model.Consultation
	if err := r.DB.Preload("User").Preload("Doctor").
		Where("is_paid = ? AND is_approved = ?", true, false).
		Find(&consultations).Error; err != nil {
		return nil, err
	}
	return consultations, nil
}
