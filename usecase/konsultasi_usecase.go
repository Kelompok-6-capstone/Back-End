package usecase

import (
	"calmind/model"
	"calmind/repository"
)

type ConsultationUsecase interface {
	CreateConsultation(consultation *model.Consultation) error
}

type ConsultationUsecaseImpl struct {
	ConsultationRepo repository.ConsultationRepository
}

func NewConsultationUsecase(cRepo repository.ConsultationRepository) ConsultationUsecase {
	return &ConsultationUsecaseImpl{ConsultationRepo: cRepo}
}

// Konsultasi
func (u *ConsultationUsecaseImpl) CreateConsultation(consultation *model.Consultation) error {
	return u.ConsultationRepo.CreateConsultation(consultation)
}
