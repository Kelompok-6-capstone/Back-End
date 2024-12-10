package repository

import (
	"calmind/model"

	"gorm.io/gorm"
)

type ConsultationRepository interface {
	CreateConsultation(*model.Consultation) (int, error)
	GetConsultationsForDoctor(doctorID int) ([]model.Consultation, error)
	GetConsultationDetails(consultationID, doctorID int) (*model.Consultation, error)
	AddRecommendation(recommendation *model.Rekomendasi) error
	GetConsultationByID(consultationID int) (*model.Consultation, error)
	UpdateConsultation(consultation *model.Consultation) error
	GetConsultationsWithDoctors(userID int) ([]model.Consultation, error)
	GetPendingConsultations() ([]model.Consultation, error)
	GetDoctorByID(doctorID int) (*model.Doctor, error)
	GetActiveConsultations() ([]model.Consultation, error)
}

type ConsultationRepositoryImpl struct {
	DB *gorm.DB
}

func NewConsultationRepositoryImpl(db *gorm.DB) *ConsultationRepositoryImpl {
	return &ConsultationRepositoryImpl{DB: db}
}

// Membuat konsultasi baru
func (r *ConsultationRepositoryImpl) CreateConsultation(consultation *model.Consultation) (int, error) {
	if err := r.DB.Preload("User").Preload("Doctor").Create(consultation).Error; err != nil {
		return 0, err
	}
	return consultation.ID, nil
}

// Mendapatkan daftar konsultasi untuk dokter
func (r *ConsultationRepositoryImpl) GetConsultationsForDoctor(doctorID int) ([]model.Consultation, error) {
	var consultations []model.Consultation
	if err := r.DB.Preload("User").Preload("Rekomendasi").
		Where("doctor_id = ? AND status = ?", doctorID, "active").
		Find(&consultations).Error; err != nil {
		return nil, err
	}
	return consultations, nil
}

// Mendapatkan detail konsultasi tertentu untuk dokter
func (r *ConsultationRepositoryImpl) GetConsultationDetails(consultationID, doctorID int) (*model.Consultation, error) {
	var consultation model.Consultation
	if err := r.DB.Preload("User").Preload("Doctor").Preload("Rekomendasi").
		Where("id = ? AND doctor_id = ?", consultationID, doctorID).
		First(&consultation).Error; err != nil {
		return nil, err
	}
	return &consultation, nil
}

// Menambahkan rekomendasi untuk konsultasi
func (r *ConsultationRepositoryImpl) AddRecommendation(recommendation *model.Rekomendasi) error {
	return r.DB.Create(recommendation).Error
}

// Mendapatkan konsultasi berdasarkan ID
func (r *ConsultationRepositoryImpl) GetConsultationByID(consultationID int) (*model.Consultation, error) {
	var consultation model.Consultation
	if err := r.DB.
		Preload("User").        // Preload relasi User
		Preload("Doctor").      // Preload relasi Doctor
		Preload("Rekomendasi"). // Preload rekomendasi
		First(&consultation, consultationID).Error; err != nil {
		return nil, err
	}
	return &consultation, nil
}

func (r *ConsultationRepositoryImpl) GetActiveConsultations() ([]model.Consultation, error) {
	var consultations []model.Consultation
	err := r.DB.Where("status = ?", "active").Find(&consultations).Error
	if err != nil {
		return nil, err
	}
	return consultations, nil
}

func (r *ConsultationRepositoryImpl) UpdateConsultation(consultation *model.Consultation) error {
	return r.DB.Save(consultation).Error
}

// Mendapatkan daftar konsultasi untuk user tertentu
func (r *ConsultationRepositoryImpl) GetConsultationsWithDoctors(userID int) ([]model.Consultation, error) {
	var consultations []model.Consultation
	if err := r.DB.
		Preload("Doctor").      // Preload relasi Doctor
		Preload("Rekomendasi"). // Preload rekomendasi
		Preload("User").        // Preload relasi User
		Where("user_id = ?", userID).
		Find(&consultations).Error; err != nil {
		return nil, err
	}
	return consultations, nil
}

// Mendapatkan konsultasi yang menunggu persetujuan admin
func (r *ConsultationRepositoryImpl) GetPendingConsultations() ([]model.Consultation, error) {
	var consultations []model.Consultation
	if err := r.DB.Preload("User").Preload("Doctor").
		Where("status = ?", "pending").
		Find(&consultations).Error; err != nil {
		return nil, err
	}
	return consultations, nil
}

// Mendapatkan dokter berdasarkan ID
func (r *ConsultationRepositoryImpl) GetDoctorByID(doctorID int) (*model.Doctor, error) {
	var doctor model.Doctor
	if err := r.DB.First(&doctor, doctorID).Error; err != nil {
		return nil, err
	}
	return &doctor, nil
}
