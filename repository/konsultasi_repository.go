package repository

import (
	"calmind/model"
	"errors"
	"time"

	"gorm.io/gorm"
)

type ConsultationRepository interface {
	CreateConsultation(*model.Consultation) error
	ApprovePayment(consultationID int) error
	// EndConsultation(consultationID int) error
	// PayConsultation(consultationID int, amount float64) error
	// GetConsultationsForDoctor(doctorID int) ([]model.Consultation, error)
	// GetConsultationDetails(consultationID, doctorID int) (*model.Consultation, error)
	// AddRecommendation(recommendation *model.Rekomendasi) error
	// GetAdminViewConsultation(consultationID int) (*model.Consultation, error)
	// GetConsultationByID(consultationID int) (*model.Consultation, error)
	// UpdateConsultation(consultation *model.Consultation) error
	// GetActiveConsultations() ([]model.Consultation, error)
	// GetConsultationsWithDoctors(userID int) ([]model.Consultation, error)
	// GetPendingConsultations() ([]model.Consultation, error)
	// AddIncomeForDoctor(doctorID int, amount float64) error
	// AddIncomeForAdmin(amount float64) error
	// GetDoctorByID(doctorID int) (*model.Doctor, error)
}

type ConsultationRepositoryImpl struct {
	DB *gorm.DB
}

func NewConsultationRepositoryImpl(db *gorm.DB) *ConsultationRepositoryImpl {
	return &ConsultationRepositoryImpl{DB: db}
}

// Membuat konsultasi baru
func (r *ConsultationRepositoryImpl) CreateConsultation(consultation *model.Consultation) (int, error) {
	err := r.DB.Create(consultation).Error
	if err != nil {
		return 0, err
	}
	return consultation.ID, nil
}

// Menyetujui pembayaran
func (r *ConsultationRepositoryImpl) ApprovePayment(consultationID int) error {
	var consultation model.Consultation
	if err := r.DB.First(&consultation, consultationID).Error; err != nil {
		return err
	}

	if consultation.IsPaid && !consultation.IsApproved {
		consultation.IsApproved = true
		consultation.StartTime = time.Now()
		consultation.Status = "active"
		return r.DB.Save(&consultation).Error
	}

	return errors.New("payment not completed or already approved")
}

// // Mengakhiri konsultasi (hanya untuk internal logika kedaluwarsa)
// func (r *ConsultationRepositoryImpl) EndConsultation(consultationID int) error {
// 	var consultation model.Consultation
// 	if err := r.DB.First(&consultation, consultationID).Error; err != nil {
// 		return err
// 	}

// 	if time.Since(consultation.StartTime).Minutes() > float64(consultation.Duration) {
// 		consultation.Status = "ended"
// 		return r.DB.Save(&consultation).Error
// 	}

// 	return errors.New("consultation is still active")
// }

// // Membayar konsultasi
// func (r *ConsultationRepositoryImpl) PayConsultation(consultationID int, amount float64) error {
// 	var consultation model.Consultation
// 	if err := r.DB.First(&consultation, consultationID).Error; err != nil {
// 		return err
// 	}

// 	if consultation.IsPaid {
// 		return errors.New("consultation already paid")
// 	}

// 	if amount < 100000 { // Harga default dokter
// 		return errors.New("insufficient payment amount")
// 	}

// 	consultation.IsPaid = true
// 	return r.DB.Save(&consultation).Error
// }

// // Mendapatkan daftar konsultasi untuk dokter
// func (r *ConsultationRepositoryImpl) GetConsultationsForDoctor(doctorID int) ([]model.Consultation, error) {
// 	var consultations []model.Consultation
// 	if err := r.DB.Preload("User").Where("doctor_id = ? AND status = ?", doctorID, "active").Find(&consultations).Error; err != nil {
// 		return nil, err
// 	}
// 	return consultations, nil
// }

// // Mendapatkan detail konsultasi untuk dokter
// func (r *ConsultationRepositoryImpl) GetConsultationDetails(consultationID, doctorID int) (*model.Consultation, error) {
// 	var consultation model.Consultation
// 	if err := r.DB.Preload("User").Preload("Doctor").Where("id = ? AND doctor_id = ?", consultationID, doctorID).First(&consultation).Error; err != nil {
// 		return nil, err
// 	}
// 	return &consultation, nil
// }

// // Menambahkan rekomendasi
// func (r *ConsultationRepositoryImpl) AddRecommendation(recommendation *model.Rekomendasi) error {
// 	return r.DB.Create(recommendation).Error
// }

// // Mendapatkan detail konsultasi untuk admin
// func (r *ConsultationRepositoryImpl) GetAdminViewConsultation(consultationID int) (*model.Consultation, error) {
// 	var consultation model.Consultation
// 	if err := r.DB.Preload("User").Preload("Doctor").First(&consultation, consultationID).Error; err != nil {
// 		return nil, err
// 	}
// 	return &consultation, nil
// }

// // Mendapatkan konsultasi berdasarkan ID
// func (r *ConsultationRepositoryImpl) GetConsultationByID(consultationID int) (*model.Consultation, error) {
// 	var consultation model.Consultation
// 	if err := r.DB.First(&consultation, consultationID).Error; err != nil {
// 		return nil, err
// 	}
// 	return &consultation, nil
// }

// // Memperbarui konsultasi
// func (r *ConsultationRepositoryImpl) UpdateConsultation(consultation *model.Consultation) error {
// 	return r.DB.Save(consultation).Error
// }

// // Mendapatkan daftar konsultasi aktif
// func (r *ConsultationRepositoryImpl) GetActiveConsultations() ([]model.Consultation, error) {
// 	var consultations []model.Consultation
// 	if err := r.DB.Where("status = ?", "active").Find(&consultations).Error; err != nil {
// 		return nil, err
// 	}
// 	return consultations, nil
// }

// // Mendapatkan daftar konsultasi dengan dokter untuk user tertentu
// func (r *ConsultationRepositoryImpl) GetConsultationsWithDoctors(userID int) ([]model.Consultation, error) {
// 	var consultations []model.Consultation
// 	err := r.DB.Preload("Doctor").Where("user_id = ?", userID).Find(&consultations).Error
// 	return consultations, err
// }

// // Mendapatkan daftar konsultasi yang belum disetujui pembayarannya
// func (r *ConsultationRepositoryImpl) GetPendingConsultations() ([]model.Consultation, error) {
// 	var consultations []model.Consultation
// 	err := r.DB.Preload("User").Where("is_paid = ? AND is_approved = ?", true, false).Find(&consultations).Error
// 	return consultations, err
// }

// func (r *ConsultationRepositoryImpl) AddIncomeForDoctor(doctorID int, amount float64) error {
// 	return r.DB.Model(&model.Doctor{}).Where("id = ?", doctorID).Update("income", gorm.Expr("income + ?", amount)).Error
// }

// func (r *ConsultationRepositoryImpl) AddIncomeForAdmin(amount float64) error {
// 	adminAccount := &model.Admin{ID: 1} // Asumsikan admin ID selalu 1
// 	return r.DB.Model(adminAccount).Update("income", gorm.Expr("income + ?", amount)).Error
// }

// func (r *ConsultationRepositoryImpl) GetDoctorByID(doctorID int) (*model.Doctor, error) {
// 	var doctor model.Doctor
// 	if err := r.DB.First(&doctor, doctorID).Error; err != nil {
// 		return nil, err
// 	}
// 	return &doctor, nil
// }
