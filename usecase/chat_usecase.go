package usecase

import (
	"calmind/model"
	"calmind/repository"
	"errors"
	"time"
)

type ChatUsecase interface {
	SendChat(chat model.Chat) (*model.ChatDTO, error)
	GetChatHistory(userID int) ([]model.ChatDTO, error)
}

type ChatUsecaseImpl struct {
	ChatRepo         repository.ChatRepository
	ConsultationRepo repository.ConsultationRepository
}

func NewChatUsecaseImpl(chatRepo repository.ChatRepository, consultationRepo repository.ConsultationRepository) *ChatUsecaseImpl {
	return &ChatUsecaseImpl{
		ChatRepo:         chatRepo,
		ConsultationRepo: consultationRepo,
	}
}

func (uc *ChatUsecaseImpl) SendChat(chat model.Chat) (*model.ChatDTO, error) {
	// Validasi konsultasi dengan pembayaran selesai
	consultations, err := uc.ConsultationRepo.GetPaidConsultationsByUserID(chat.UserID)
	if err != nil || len(consultations) == 0 {
		return nil, errors.New("no valid consultations found for this user")
	}

	// Validasi jika ada konsultasi aktif
	isValid := false
	for _, consultation := range consultations {
		endTime := consultation.StartTime.Add(time.Duration(consultation.Duration) * time.Minute)
		if time.Now().Before(endTime) {
			isValid = true
			break
		}
	}

	if !isValid {
		return nil, errors.New("all consultations have ended or are invalid")
	}

	// Simpan pesan chat
	chat.CreatedAt = time.Now()
	err = uc.ChatRepo.SaveChat(&chat)
	if err != nil {
		return nil, errors.New("failed to save chat")
	}

	// Ambil nama pengirim
	var senderName string
	if chat.SenderType == "user" {
		user, _ := uc.ChatRepo.GetUserByID(chat.SenderID)
		senderName = user.Username
	} else if chat.SenderType == "doctor" {
		doctor, _ := uc.ChatRepo.GetDoctorByID(chat.SenderID)
		senderName = doctor.Username
	}

	return &model.ChatDTO{
		ID:         chat.ID,
		UserID:     chat.UserID,
		SenderID:   chat.SenderID,
		SenderName: senderName,
		Message:    chat.Message,
		SenderType: chat.SenderType,
		CreatedAt:  chat.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (uc *ChatUsecaseImpl) GetChatHistory(userID int) ([]model.ChatDTO, error) {
	// Ambil riwayat chat
	chats, err := uc.ChatRepo.GetChatHistoryByUser(userID)
	if err != nil {
		return nil, err
	}

	// Konversi ke DTO
	var chatDTOs []model.ChatDTO
	for _, chat := range chats {
		var senderName string
		if chat.SenderType == "user" && chat.User != nil {
			senderName = chat.User.Username
		} else if chat.SenderType == "doctor" && chat.Doctor != nil {
			senderName = chat.Doctor.Username
		}

		chatDTOs = append(chatDTOs, model.ChatDTO{
			ID:         chat.ID,
			UserID:     chat.UserID,
			SenderID:   chat.SenderID,
			SenderName: senderName,
			Message:    chat.Message,
			SenderType: chat.SenderType,
			CreatedAt:  chat.CreatedAt.Format(time.RFC3339),
		})
	}
	return chatDTOs, nil
}
