package usecase

import (
	"calmind/model"
	"calmind/repository"
	"errors"
	"log"
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
	if err != nil {
		log.Printf("Error fetching consultation: %v", err)
		return err
	}
	if consultation == nil {
		log.Printf("No valid consultation found for user_id=%d, doctor_id=%d", userID, doctorID)
		return errors.New("no valid consultation found between user and doctor")
	}

	// Tambahkan logging waktu konsultasi
	log.Printf("Start Time: %v", consultation.StartTime)
	log.Printf("Duration: %d minutes", consultation.Duration)
	log.Printf("End Time: %v", consultation.StartTime.Add(time.Duration(consultation.Duration)*time.Minute))
	log.Printf("Current Time: %v", time.Now())
	log.Printf("Is consultation expired: %v", time.Now().After(consultation.StartTime.Add(time.Duration(consultation.Duration)*time.Minute)))

	log.Printf("Valid consultation found: %+v", consultation)

	// Validasi status konsultasi
	if consultation.Status != "approved" || consultation.PaymentStatus != "paid" ||
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
	log.Printf("Saving message: %+v", chat)
	return uc.ChatRepo.SaveMessage(chat)
}

func (uc *ChatUsecaseImpl) GetMessages(userID, doctorID int) ([]model.ChatMessage, error) {
	// Validasi konsultasi valid
	consultation, err := uc.ConsultationRepo.GetConsultationBetweenUserAndDoctor(userID, doctorID)
	if err != nil {
		return nil, err // Error dari database
	}
	if consultation == nil {
		return nil, errors.New("no valid consultation found between user and doctor")
	}

	log.Printf("Start Time: %v", consultation.StartTime)
	log.Printf("Duration: %d minutes", consultation.Duration)
	log.Printf("End Time: %v", consultation.StartTime.Add(time.Duration(consultation.Duration)*time.Minute))
	log.Printf("Current Time: %v", time.Now())
	log.Printf("Is consultation expired: %v", time.Now().After(consultation.StartTime.Add(time.Duration(consultation.Duration)*time.Minute)))

	// Validasi status konsultasi
	if consultation.Status != "approved" || consultation.PaymentStatus != "paid" ||
		time.Now().After(consultation.StartTime.Add(time.Duration(consultation.Duration)*time.Minute)) {
		return nil, errors.New("consultation is not valid for chat")
	}

	// Ambil pesan
	return uc.ChatRepo.GetMessages(userID, doctorID)
}
