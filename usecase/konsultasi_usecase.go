package usecase

import (
	"calmind/model"
	"calmind/repository"
	"errors"
	"time"
)

type ConsultationUsecase interface {
	CreateConsultation(consultation *model.Consultation) error
	GetConsultationsAllDoctor(doctorID int) ([]model.Consultation, error)
	GetConsultationByID(consultationID int) (*model.Consultation, error)
	UpdateRecommendation(consultationID int, recommendation string) error
	UpdatePaymentStatus(consultationID int, isPaid bool) error
	GetUnpaidConsultations() ([]model.Consultation, error)
	MarkExpiredConsultations() error
}

type ConsultationUsecaseImpl struct {
	ConsultationRepo repository.ConsultationRepository
}

func NewConsultationUsecase(cRepo repository.ConsultationRepository) ConsultationUsecase {
	return &ConsultationUsecaseImpl{ConsultationRepo: cRepo}
}

// Membuat konsultasi baru
func (u *ConsultationUsecaseImpl) CreateConsultation(consultation *model.Consultation) error {
	// Set waktu mulai konsultasi ke waktu saat ini
	consultation.StartTime = time.Now()
	return u.ConsultationRepo.CreateConsultation(consultation)
}

// Mendapatkan daftar konsultasi berdasarkan doctorID
func (u *ConsultationUsecaseImpl) GetConsultationsAllDoctor(doctorID int) ([]model.Consultation, error) {
	var consultations []model.Consultation
	if err := u.ConsultationRepo.FindByDoctorID(doctorID, &consultations); err != nil {
		return nil, err
	}
	return consultations, nil
}

// Mendapatkan detail konsultasi berdasarkan consultationID
func (u *ConsultationUsecaseImpl) GetConsultationByID(consultationID int) (*model.Consultation, error) {
	var consultation model.Consultation
	err := u.ConsultationRepo.FindByConsultationID(consultationID, &consultation)
	if err != nil {
		return nil, err
	}
	return &consultation, nil
}

// Memperbarui rekomendasi konsultasi
func (u *ConsultationUsecaseImpl) UpdateRecommendation(consultationID int, recommendation string) error {
	return u.ConsultationRepo.UpdateRecommendation(consultationID, recommendation)
}

// Mengubah status pembayaran
func (u *ConsultationUsecaseImpl) UpdatePaymentStatus(consultationID int, isPaid bool) error {
	consultation, err := u.GetConsultationByID(consultationID)
	if err != nil {
		return errors.New("konsultasi tidak ditemukan")
	}
	if consultation.IsPaid && isPaid {
		return errors.New("pembayaran sudah selesai sebelumnya")
	}
	return u.ConsultationRepo.UpdatePaymentStatus(consultationID, isPaid)
}

// Mendapatkan daftar konsultasi yang belum dibayar
func (u *ConsultationUsecaseImpl) GetUnpaidConsultations() ([]model.Consultation, error) {
	var consultations []model.Consultation
	err := u.ConsultationRepo.FindUnpaidConsultations(&consultations)
	if err != nil {
		return nil, err
	}
	return consultations, nil
}

// Menandai konsultasi yang kedaluwarsa
func (u *ConsultationUsecaseImpl) MarkExpiredConsultations() error {
	err := u.ConsultationRepo.ExpireConsultations()
	if err != nil {
		return err
	}
	return nil
}
