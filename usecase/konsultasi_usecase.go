package usecase

import (
	"calmind/model"
	"calmind/repository"
)

type ConsultationUsecase interface {
	CreateConsultation(consultation *model.Consultation) error
	GetConsultationsAllDoctor(doctorID int) ([]model.Consultation, error)
	GetConsultationByID(consultationID string) (*model.Consultation, error)
	UpdateRecommendation(consultationID int, recommendation string) error
}

type ConsultationUsecaseImpl struct {
	ConsultationRepo repository.ConsultationRepository
}

func NewConsultationUsecase(cRepo repository.ConsultationRepository) ConsultationUsecase {
	return &ConsultationUsecaseImpl{ConsultationRepo: cRepo}
}

// Membuat konsultasi baru
func (u *ConsultationUsecaseImpl) CreateConsultation(consultation *model.Consultation) error {
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
func (u *ConsultationUsecaseImpl) GetConsultationByID(consultationID string) (*model.Consultation, error) {
	var consultation model.Consultation
	err := u.ConsultationRepo.FindByConsultationID(consultationID, &consultation)
	if err != nil {
		return nil, err
	}
	return &consultation, nil
}

// Memperbarui rekomendasi konsultasi
func (u *ConsultationUsecaseImpl) UpdateRecommendation(consultationID int, recommendation string) error {
	err := u.ConsultationRepo.UpdateRecommendation(consultationID, recommendation)
	return err
}
