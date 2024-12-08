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
	UpdateApprovalStatus(consultationID int, isApproved bool) error
	GetUnpaidConsultations() ([]model.Consultation, error)
	GetPendingApprovals() ([]model.Consultation, error)
	MarkExpiredConsultations() error
	GetDoctorByID(doctorID int) (*model.Doctor, error)
}

type ConsultationUsecaseImpl struct {
	ConsultationRepo repository.ConsultationRepository
}

func NewConsultationUsecase(cRepo repository.ConsultationRepository) ConsultationUsecase {
	return &ConsultationUsecaseImpl{ConsultationRepo: cRepo}
}

// Membuat konsultasi baru
func (u *ConsultationUsecaseImpl) CreateConsultation(consultation *model.Consultation) error {
	// Validasi dokter dan harga
	doctor, err := u.GetDoctorByID(consultation.DoctorID)
	if err != nil {
		return errors.New("dokter tidak ditemukan")
	}
	// Set durasi dan waktu mulai
	consultation.StartTime = time.Now()
	consultation.Duration = 120 // Default: 2 jam
	consultation.Status = "pending"
	consultation.IsPaid = false
	consultation.IsApproved = false
	consultation.Doctor = *doctor
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
	consultation, err := u.GetConsultationByID(consultationID)
	if err != nil {
		return errors.New("konsultasi tidak ditemukan")
	}
	// Validasi jika konsultasi belum disetujui atau sudah kedaluwarsa
	now := time.Now()
	if !consultation.IsApproved {
		return errors.New("konsultasi belum disetujui oleh admin")
	}
	if now.After(consultation.StartTime.Add(time.Duration(consultation.Duration) * time.Minute)) {
		return errors.New("konsultasi sudah kedaluwarsa")
	}
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

// Mengubah status persetujuan admin
func (u *ConsultationUsecaseImpl) UpdateApprovalStatus(consultationID int, isApproved bool) error {
	consultation, err := u.GetConsultationByID(consultationID)
	if err != nil {
		return errors.New("konsultasi tidak ditemukan")
	}
	if !consultation.IsPaid {
		return errors.New("pembayaran belum selesai")
	}
	return u.ConsultationRepo.UpdateApprovalStatus(consultationID, isApproved)
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

// Mendapatkan daftar konsultasi yang menunggu persetujuan admin
func (u *ConsultationUsecaseImpl) GetPendingApprovals() ([]model.Consultation, error) {
	var consultations []model.Consultation
	err := u.ConsultationRepo.FindPendingApproval(&consultations)
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

// Mendapatkan data dokter berdasarkan ID
func (u *ConsultationUsecaseImpl) GetDoctorByID(doctorID int) (*model.Doctor, error) {
	var doctor model.Doctor
	err := u.ConsultationRepo.FindDoctorByID(doctorID, &doctor)
	if err != nil {
		return nil, err
	}
	return &doctor, nil
}
