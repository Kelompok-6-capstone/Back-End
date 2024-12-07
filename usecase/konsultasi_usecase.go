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
	UpdateConsultationStatus(consultationID int, status string) error
	UpdatePaymentStatus(consultationID int, paymentStatus string) error
	UpdateMessage(consultationID int, message string) error
}

type ConsultationUsecaseImpl struct {
	ConsultationRepo repository.ConsultationRepository
}

func NewConsultationUsecase(repo repository.ConsultationRepository) ConsultationUsecase {
	return &ConsultationUsecaseImpl{ConsultationRepo: repo}
}

func (u *ConsultationUsecaseImpl) CreateConsultation(consultation *model.Consultation) error {
	return u.ConsultationRepo.CreateConsultation(consultation)
}

func (u *ConsultationUsecaseImpl) GetConsultationsAllDoctor(doctorID int) ([]model.Consultation, error) {
	var consultations []model.Consultation
	err := u.ConsultationRepo.FindByDoctorID(doctorID, &consultations)
	return consultations, err
}

func (u *ConsultationUsecaseImpl) GetConsultationByID(consultationID string) (*model.Consultation, error) {
	var consultation model.Consultation
	err := u.ConsultationRepo.FindByConsultationID(consultationID, &consultation)
	return &consultation, err
}

func (u *ConsultationUsecaseImpl) UpdateRecommendation(consultationID int, recommendation string) error {
	return u.ConsultationRepo.UpdateRecommendation(consultationID, recommendation)
}

func (u *ConsultationUsecaseImpl) UpdateConsultationStatus(consultationID int, status string) error {
	return u.ConsultationRepo.UpdateStatus(consultationID, status)
}

func (u *ConsultationUsecaseImpl) UpdatePaymentStatus(consultationID int, paymentStatus string) error {
	return u.ConsultationRepo.UpdatePaymentStatus(consultationID, paymentStatus)
}

func (u *ConsultationUsecaseImpl) UpdateMessage(consultationID int, message string) error {
	return u.ConsultationRepo.UpdateMessage(consultationID, message)
}
