package repository

import (
	"calmind/model"
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
)

type ConsultationRepository interface {
	CreateConsultation(*model.Consultation) (int, error)
	GetConsultationsForDoctor(doctorID int) ([]model.Consultation, error)
	FindConsultationsByDoctorAndName(doctorID int, searchName string, consultations *[]model.Consultation) error
	GetConsultationDetails(consultationID, doctorID int) (*model.Consultation, error)
	AddRecommendation(recommendation *model.Rekomendasi) error
	GetConsultationByID(consultationID int) (*model.Consultation, error)
	UpdateConsultation(consultation *model.Consultation) error
	GetConsultationsWithDoctors(userID int) ([]model.Consultation, error)
	GetPendingConsultations() ([]model.Consultation, error)
	GetDoctorByID(doctorID int) (*model.Doctor, error)
	GetActiveConsultations() ([]model.Consultation, error)
	GetAllStatusConsultations() ([]model.Consultation, error)
	GetApprovedConsultations() ([]model.Consultation, error)
	ValidateUserAndDoctor(userID, doctorID int) error
	GetConsultationByOrderID(orderID string) (*model.Consultation, error)
	GetValidConsultations(userID, doctorID int) ([]model.Consultation, error)
	GetAllConsultationsForDoctor(doctorID int) ([]model.Consultation, error)
}

type ConsultationRepositoryImpl struct {
	DB *gorm.DB
}

func NewConsultationRepositoryImpl(db *gorm.DB) *ConsultationRepositoryImpl {
	return &ConsultationRepositoryImpl{DB: db}
}

func (r *ConsultationRepositoryImpl) GetValidConsultations(userID, doctorID int) ([]model.Consultation, error) {
	var consultations []model.Consultation
	err := r.DB.Where("user_id = ? AND doctor_id = ? AND status = ? AND payment_status = ?", userID, doctorID, "approved", "paid").
		Order("start_time").
		Find(&consultations).Error

	if err != nil {
		return nil, err
	}
	return consultations, nil
}

func (r *ConsultationRepositoryImpl) CreateConsultation(consultation *model.Consultation) (int, error) {
	// Simpan konsultasi
	if err := r.DB.Create(consultation).Error; err != nil {
		return 0, fmt.Errorf("failed to create consultation: %w", err)
	}

	// Preload data User dan Doctor
	if err := r.DB.Preload("User").Preload("Doctor").Preload("Rekomendasi").First(&consultation, consultation.ID).Error; err != nil {
		fmt.Println("Error saat preload konsultasi:", err)
		return 0, err
	}

	fmt.Printf("Consultation setelah Preload: %+v\n", consultation)
	return consultation.ID, nil
}

func (r *ConsultationRepositoryImpl) ValidateUserAndDoctor(userID, doctorID int) error {
	var user model.User
	if err := r.DB.First(&user, userID).Error; err != nil {
		return fmt.Errorf("user not found with ID %d", userID)
	}

	var doctor model.Doctor
	if err := r.DB.First(&doctor, doctorID).Error; err != nil {
		return fmt.Errorf("doctor not found with ID %d", doctorID)
	}

	return nil
}

// Mendapatkan daftar konsultasi untuk dokter
func (r *ConsultationRepositoryImpl) GetAllConsultationsForDoctor(doctorID int) ([]model.Consultation, error) {
	var consultations []model.Consultation
	if err := r.DB.Preload("User").Preload("Doctor").Preload("Rekomendasi").
		Where("doctor_id = ? AND status IN ?", doctorID, []string{"paid", "approved", "expired"}).
		Find(&consultations).Error; err != nil {
		return nil, err
	}
	fmt.Printf("Consultations found: %+v\n", consultations) // Logging untuk debug
	return consultations, nil
}
func (r *ConsultationRepositoryImpl) GetConsultationsForDoctor(doctorID int) ([]model.Consultation, error) {
	var consultations []model.Consultation
	if err := r.DB.Preload("User").Preload("Doctor").Preload("Rekomendasi").
		Where("doctor_id = ? AND status IN ?", doctorID, []string{"paid", "approved"}).
		Find(&consultations).Error; err != nil {
		return nil, err
	}
	fmt.Printf("Consultations found: %+v\n", consultations) // Logging untuk debug
	return consultations, nil
}

// Mendapatkan konsultasi berdasarkan doctorID dan nama user
func (r *ConsultationRepositoryImpl) FindConsultationsByDoctorAndName(doctorID int, searchName string, consultations *[]model.Consultation) error {
	query := r.DB.Preload("User").Preload("Doctor").
		Where("doctor_id = ?", doctorID)

	// Gunakan sub-query atau join untuk mencocokkan nama user
	if searchName != "" {
		query = query.Joins("JOIN users ON users.id = consultations.user_id").
			Where("users.username LIKE ?", "%"+searchName+"%")
	}

	err := query.Find(consultations).Error
	if err != nil {
		log.Printf("Query error in FindConsultationsByDoctorAndName: %v", err) // Logging untuk debug
	}
	return err
}

// Mendapatkan detail konsultasi tertentu untuk dokter
func (r *ConsultationRepositoryImpl) GetConsultationDetails(consultationID, doctorID int) (*model.Consultation, error) {
	var consultation model.Consultation
	if err := r.DB.Preload("User").Preload("Doctor.Title").Preload("Rekomendasi").
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
func (r *ConsultationRepositoryImpl) GetActiveConsultations() ([]model.Consultation, error) {
	var consultations []model.Consultation
	err := r.DB.Where("status = ? AND start_time <= ?", "approved", time.Now()).
		Find(&consultations).Error
	if err != nil {
		return nil, err
	}
	return consultations, nil
}

func (r *ConsultationRepositoryImpl) UpdateConsultation(consultation *model.Consultation) error {
	return r.DB.Save(consultation).Error
}

func (r *ConsultationRepositoryImpl) GetConsultationByID(consultationID int) (*model.Consultation, error) {
	var consultation model.Consultation
	err := r.DB.Preload("User").Preload("Doctor.Title").Preload("Rekomendasi").First(&consultation, consultationID).Error
	if err != nil {
		return nil, err
	}
	return &consultation, nil
}

// Mendapatkan daftar konsultasi untuk user tertentu
func (r *ConsultationRepositoryImpl) GetConsultationsWithDoctors(userID int) ([]model.Consultation, error) {
	var consultations []model.Consultation
	if err := r.DB.
		Preload("Doctor.Title"). // Preload relasi Doctor
		Preload("Rekomendasi").  // Preload rekomendasi
		Preload("User").         // Preload relasi User
		Where("user_id = ?", userID).
		Find(&consultations).Error; err != nil {
		return nil, err
	}
	return consultations, nil
}

// Mendapatkan konsultasi yang menunggu persetujuan admin
func (r *ConsultationRepositoryImpl) GetPendingConsultations() ([]model.Consultation, error) {
	var consultations []model.Consultation
	if err := r.DB.Preload("User").Preload("Doctor.Title").
		Where("status = ?", "pending").
		Find(&consultations).Error; err != nil {
		return nil, err
	}
	return consultations, nil
}
func (r *ConsultationRepositoryImpl) GetApprovedConsultations() ([]model.Consultation, error) {
	var consultations []model.Consultation
	if err := r.DB.Preload("User").Preload("Doctor").
		Where("status = ?", "approved").
		Find(&consultations).Error; err != nil {
		return nil, err
	}
	return consultations, nil
}
func (r *ConsultationRepositoryImpl) GetAllStatusConsultations() ([]model.Consultation, error) {
	var consultations []model.Consultation
	if err := r.DB.Preload("User").Preload("Doctor").
		Where("status IN ?", []string{"approved", "pending"}). // Menampilkan status approved dan pending
		Find(&consultations).Error; err != nil {
		return nil, err
	}
	return consultations, nil
}

func (r *ConsultationRepositoryImpl) GetConsultationByOrderID(orderID string) (*model.Consultation, error) {
	var consultation model.Consultation
	err := r.DB.Where("order_id = ?", orderID).First(&consultation).Error
	if err != nil {
		log.Printf("Consultation not found for order_id=%s: %v", orderID, err)
		return nil, err
	}
	return &consultation, nil
}

// Mendapatkan dokter berdasarkan ID
func (r *ConsultationRepositoryImpl) GetDoctorByID(doctorID int) (*model.Doctor, error) {
	var doctor model.Doctor
	if err := r.DB.First(&doctor, doctorID).Error; err != nil {
		return nil, err
	}
	return &doctor, nil
}
