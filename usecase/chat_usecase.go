package usecase

import (
	"calmind/model"
	"calmind/repository"
	"errors"
	"time"
)

type ChatUsecase interface {
	SendMessage(userID, doctorID, senderID int, message string) error
	GetMessages(userID, doctorID int) ([]model.ChatMessage, error)
}

type ChatUsecaseImpl struct {
	ChatRepo         repository.ChatRepository
	ConsultationRepo repository.ConsultationRepository
}

func (uc *ChatUsecaseImpl) SendMessage(userID, doctorID, senderID int, message string) error {
	// Validasi konsultasi valid
	consultation, err := uc.ConsultationRepo.GetConsultationBetweenUserAndDoctor(userID, doctorID)
	if err != nil || consultation.Status != "approved" || consultation.PaymentStatus != "paid" ||
		time.Now().After(consultation.StartTime.Add(time.Duration(consultation.Duration)*time.Minute)) {
		return errors.New("consultation is not valid for chat")
	}

	// Simpan pesan
	chat := &model.ChatMessage{
		UserID:   userID,
		DoctorID: doctorID,
		SenderID: senderID,
		Message:  message,
	}
	return uc.ChatRepo.SaveMessage(chat)
}

func (uc *ChatUsecaseImpl) GetMessages(userID, doctorID int) ([]model.ChatMessage, error) {
	// Validasi konsultasi valid
	consultation, err := uc.ConsultationRepo.GetConsultationBetweenUserAndDoctor(userID, doctorID)
	if err != nil || consultation.Status != "approved" || consultation.PaymentStatus != "paid" ||
		time.Now().After(consultation.StartTime.Add(time.Duration(consultation.Duration)*time.Minute)) {
		return nil, errors.New("consultation is not valid for chat")
	}

	// Ambil pesan
	return uc.ChatRepo.GetMessages(userID, doctorID)
}
