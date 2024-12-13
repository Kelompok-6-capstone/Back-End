package repository

import (
	"calmind/model"
	"fmt"
	"time"

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
	GetActiveConsultations(userID int, doctorID int) ([]model.Consultation, error)
	GetAllStatusConsultations() ([]model.Consultation, error)
	GetApprovedConsultations() ([]model.Consultation, error)
	GetAllActiveConsultations() ([]model.Consultation, error)
}

type ConsultationRepositoryImpl struct {
	DB *gorm.DB
}

func NewConsultationRepositoryImpl(db *gorm.DB) *ConsultationRepositoryImpl {
	return &ConsultationRepositoryImpl{DB: db}
}

// Fungsi untuk mendapatkan konsultasi aktif berdasarkan userID dan doctorID
func (r *ConsultationRepositoryImpl) GetActiveConsultations(userID int, doctorID int) ([]model.Consultation, error) {
	var consultations []model.Consultation
	err := r.DB.Where("user_id = ? AND doctor_id = ? AND status IN ('approved', 'paid') AND start_time + interval duration minute > ?", userID, doctorID, time.Now()).
		Find(&consultations).Error
	if err != nil {
		return nil, err
	}
	return consultations, nil
}

func (r *ConsultationRepositoryImpl) GetAllActiveConsultations() ([]model.Consultation, error) {
	var consultations []model.Consultation
	err := r.DB.Where("status IN ('approved', 'paid') AND start_time + interval duration minute > ?", time.Now()).
		Find(&consultations).Error
	if err != nil {
		return nil, err
	}
	return consultations, nil
}

// Fungsi untuk membuat konsultasi baru
func (r *ConsultationRepositoryImpl) CreateConsultation(consultation *model.Consultation) (int, error) {
	if err := r.DB.Create(consultation).Error; err != nil {
		fmt.Println("Error saat menyimpan konsultasi:", err)
		return 0, err
	}

	if err := r.DB.Preload("User").Preload("Doctor").First(&consultation, consultation.ID).Error; err != nil {
		fmt.Println("Error saat preload konsultasi:", err)
		return 0, err
	}

	return consultation.ID, nil
}

// Mendapatkan daftar konsultasi untuk dokter tertentu
func (r *ConsultationRepositoryImpl) GetConsultationsForDoctor(doctorID int) ([]model.Consultation, error) {
	var consultations []model.Consultation
	err := r.DB.Preload("User").Preload("Doctor").Preload("Rekomendasi").
		Where("doctor_id = ? AND status IN ?", doctorID, []string{"paid", "approved"}).
		Find(&consultations).Error
	return consultations, err
}

// Mendapatkan detail konsultasi berdasarkan consultationID dan doctorID
func (r *ConsultationRepositoryImpl) GetConsultationDetails(consultationID, doctorID int) (*model.Consultation, error) {
	var consultation model.Consultation
	err := r.DB.Preload("User").Preload("Doctor").Preload("Rekomendasi").
		Where("id = ? AND doctor_id = ?", consultationID, doctorID).
		First(&consultation).Error
	return &consultation, err
}

// Menambahkan rekomendasi untuk konsultasi tertentu
func (r *ConsultationRepositoryImpl) AddRecommendation(recommendation *model.Rekomendasi) error {
	return r.DB.Create(recommendation).Error
}

// Mendapatkan konsultasi berdasarkan ID
func (r *ConsultationRepositoryImpl) GetConsultationByID(consultationID int) (*model.Consultation, error) {
	var consultation model.Consultation
	err := r.DB.Preload("User").Preload("Doctor").First(&consultation, consultationID).Error
	return &consultation, err
}

// Update data konsultasi
func (r *ConsultationRepositoryImpl) UpdateConsultation(consultation *model.Consultation) error {
	return r.DB.Save(consultation).Error
}

// Mendapatkan daftar konsultasi untuk user tertentu
func (r *ConsultationRepositoryImpl) GetConsultationsWithDoctors(userID int) ([]model.Consultation, error) {
	var consultations []model.Consultation
	err := r.DB.Preload("Doctor").Preload("Rekomendasi").Preload("User").
		Where("user_id = ?", userID).Find(&consultations).Error
	return consultations, err
}

// Mendapatkan konsultasi yang menunggu persetujuan admin
func (r *ConsultationRepositoryImpl) GetPendingConsultations() ([]model.Consultation, error) {
	var consultations []model.Consultation
	err := r.DB.Preload("User").Preload("Doctor").
		Where("status = ?", "pending").
		Find(&consultations).Error
	return consultations, err
}

// Mendapatkan konsultasi dengan status disetujui (approved)
func (r *ConsultationRepositoryImpl) GetApprovedConsultations() ([]model.Consultation, error) {
	var consultations []model.Consultation
	err := r.DB.Preload("User").Preload("Doctor").
		Where("status = ?", "approved").
		Find(&consultations).Error
	return consultations, err
}

// Mendapatkan semua konsultasi dengan status approved atau pending
func (r *ConsultationRepositoryImpl) GetAllStatusConsultations() ([]model.Consultation, error) {
	var consultations []model.Consultation
	err := r.DB.Preload("User").Preload("Doctor").
		Where("status IN ?", []string{"approved", "pending"}).
		Find(&consultations).Error
	return consultations, err
}

// Mendapatkan dokter berdasarkan ID
func (r *ConsultationRepositoryImpl) GetDoctorByID(doctorID int) (*model.Doctor, error) {
	var doctor model.Doctor
	err := r.DB.First(&doctor, doctorID).Error
	return &doctor, err
}
