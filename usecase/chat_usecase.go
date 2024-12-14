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
	// Ambil semua konsultasi valid
	consultations, err := uc.ConsultationRepo.GetValidConsultations(userID, doctorID)
	if err != nil {
		log.Printf("Error fetching consultations: %v", err)
		return err
	}

	// Jika tidak ada konsultasi valid
	if len(consultations) == 0 {
		log.Printf("No valid consultations found for user_id=%d, doctor_id=%d", userID, doctorID)
		return errors.New("no valid consultation found between user and doctor")
	}

	// Hitung total waktu durasi dan waktu berakhir
	var totalDuration time.Duration
	var startTime time.Time

	for i, consultation := range consultations {
		totalDuration += time.Duration(consultation.Duration) * time.Minute
		if i == 0 {
			startTime = consultation.StartTime
		}
	}

	// Waktu akhir dari total durasi
	endTime := startTime.Add(totalDuration)

	log.Printf("Start Time: %v", startTime)
	log.Printf("Total Duration: %v", totalDuration)
	log.Printf("End Time: %v", endTime)
	log.Printf("Current Time: %v", time.Now())

	// Validasi apakah total waktu sudah habis
	if time.Now().After(endTime) {
		log.Printf("Consultation time expired for user_id=%d, doctor_id=%d", userID, doctorID)
		return errors.New("consultation time has expired")
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
	// Ambil semua konsultasi valid
	consultations, err := uc.ConsultationRepo.GetValidConsultations(userID, doctorID)
	if err != nil {
		log.Printf("Error fetching consultations: %v", err)
		return nil, err
	}

	if len(consultations) == 0 {
		return nil, errors.New("no valid consultation found between user and doctor")
	}

	// Hitung total waktu durasi dan waktu berakhir
	var totalDuration time.Duration
	var startTime time.Time

	for i, consultation := range consultations {
		totalDuration += time.Duration(consultation.Duration) * time.Minute
		if i == 0 {
			startTime = consultation.StartTime
		}
	}

	// Waktu akhir dari total durasi
	endTime := startTime.Add(totalDuration)

	if time.Now().After(endTime) {
		return nil, errors.New("consultation time has expired")
	}

	// Ambil pesan
	return uc.ChatRepo.GetMessages(userID, doctorID)
}
