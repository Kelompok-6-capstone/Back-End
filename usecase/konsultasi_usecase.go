package usecase

import (
	"calmind/model"
	"calmind/repository"
	"errors"
	"time"
)

type ConsultationUsecase interface {
	CreateConsultation(userID, doctorID int, title, description string) error
	PayConsultation(userID, consultationID int, amount float64) error
	GetConsultationByID(consultationID int) (*model.Consultation, error)
	MarkExpiredConsultations() error
	ApprovePayment(adminID, consultationID int) error
	ViewConsultationDetailsForAdmin(consultationID int) (*model.Consultation, error)
	GetConsultationsForDoctor(doctorID int) ([]model.Consultation, error)
	ViewConsultationDetails(doctorID, consultationID int) (*model.Consultation, error)
	AddRecommendation(doctorID, consultationID int, recommendation string) error
}

type ConsultationUsecaseImpl struct {
	Repo *repository.ConsultationRepositoryImpl
}

func NewConsultationUsecaseImpl(repo *repository.ConsultationRepositoryImpl) *ConsultationUsecaseImpl {
	return &ConsultationUsecaseImpl{Repo: repo}
}

// Membuat konsultasi
func (uc *ConsultationUsecaseImpl) CreateConsultation(userID, doctorID int, title, description string) error {
	consultation := &model.Consultation{
		UserID:      userID,
		DoctorID:    doctorID,
		Title:       title,
		Description: description,
		Status:      "pending",
		IsPaid:      false,
		IsApproved:  false,
	}
	return uc.Repo.CreateConsultation(consultation)
}

// Membayar konsultasi
func (uc *ConsultationUsecaseImpl) PayConsultation(userID, consultationID int, amount float64) error {
	consultation, err := uc.Repo.GetConsultationByID(consultationID)
	if err != nil {
		return errors.New("consultation not found")
	}

	if consultation.UserID != userID {
		return errors.New("unauthorized access to consultation")
	}

	if consultation.IsPaid {
		return errors.New("consultation already paid")
	}

	// Logika pembayaran di sini, misalnya cek jumlah pembayaran
	if amount < 100000 { // Harga default 100rb
		return errors.New("insufficient amount for consultation")
	}

	// Tandai sebagai telah dibayar
	consultation.IsPaid = true
	return uc.Repo.UpdateConsultation(consultation)
}

// Mendapatkan konsultasi berdasarkan ID
func (uc *ConsultationUsecaseImpl) GetConsultationByID(consultationID int) (*model.Consultation, error) {
	consultation, err := uc.Repo.GetConsultationByID(consultationID)
	if err != nil {
		return nil, err
	}

	// Periksa apakah konsultasi telah kedaluwarsa
	now := time.Now()
	endTime := consultation.StartTime.Add(time.Duration(consultation.Duration) * time.Minute)
	if now.After(endTime) && consultation.Status != "ended" {
		consultation.Status = "ended"
		err = uc.Repo.UpdateConsultation(consultation)
		if err != nil {
			return nil, errors.New("failed to update consultation status")
		}
	}

	return consultation, nil
}

// Menandai konsultasi yang sudah kedaluwarsa
func (uc *ConsultationUsecaseImpl) MarkExpiredConsultations() error {
	consultations, err := uc.Repo.GetActiveConsultations()
	if err != nil {
		return err
	}

	now := time.Now()
	for _, consultation := range consultations {
		endTime := consultation.StartTime.Add(time.Duration(consultation.Duration) * time.Minute)
		if now.After(endTime) {
			consultation.Status = "ended"
			if err := uc.Repo.UpdateConsultation(&consultation); err != nil {
				return errors.New("failed to update consultation status")
			}
		}
	}
	return nil
}

// Menyetujui pembayaran
func (uc *ConsultationUsecaseImpl) ApprovePayment(adminID, consultationID int) error {
	consultation, err := uc.Repo.GetConsultationByID(consultationID)
	if err != nil {
		return errors.New("consultation not found")
	}

	if consultation.IsApproved {
		return errors.New("payment already approved")
	}

	consultation.IsApproved = true
	consultation.StartTime = time.Now() // Durasi dimulai
	consultation.Status = "active"

	return uc.Repo.UpdateConsultation(consultation)
}

// Mendapatkan detail konsultasi untuk admin
func (uc *ConsultationUsecaseImpl) ViewConsultationDetailsForAdmin(consultationID int) (*model.Consultation, error) {
	return uc.Repo.GetConsultationByID(consultationID)
}

// Mendapatkan daftar konsultasi untuk dokter
func (uc *ConsultationUsecaseImpl) GetConsultationsForDoctor(doctorID int) ([]model.Consultation, error) {
	return uc.Repo.GetConsultationsForDoctor(doctorID)
}

// Mendapatkan detail konsultasi tertentu untuk dokter
func (uc *ConsultationUsecaseImpl) ViewConsultationDetails(doctorID, consultationID int) (*model.Consultation, error) {
	consultation, err := uc.Repo.GetConsultationDetails(consultationID, doctorID)
	if err != nil {
		return nil, err
	}

	if consultation.DoctorID != doctorID {
		return nil, errors.New("unauthorized access to consultation")
	}

	return consultation, nil
}

// Menambahkan rekomendasi untuk konsultasi
func (uc *ConsultationUsecaseImpl) AddRecommendation(doctorID, consultationID int, recommendation string) error {
	consultation, err := uc.Repo.GetConsultationDetails(consultationID, doctorID)
	if err != nil {
		return err
	}

	if consultation.DoctorID != doctorID {
		return errors.New("unauthorized access to consultation")
	}

	recommendationObj := &model.Rekomendasi{
		ConsultationID: consultation.ID,
		DoctorID:       doctorID,
		Rekomendasi:    recommendation,
	}
	return uc.Repo.AddRecommendation(recommendationObj)
}
