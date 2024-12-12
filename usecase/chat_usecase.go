package usecase

import (
	"calmind/model"
	"calmind/repository"
	"errors"
	"time"
)

type ChatUsecase interface {
	SendChat(chat model.Chat) (*model.ChatDTO, error)
	GetChatHistory(consultationID int) ([]model.ChatDTO, error)
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
	consultation, err := uc.ConsultationRepo.GetConsultationByID(chat.ConsultationID)
	if err != nil {
		return nil, errors.New("consultation not found")
	}

	if consultation.Status != "paid" && consultation.Status != "approved" {
		return nil, errors.New("payment not completed")
	}

	endTime := consultation.StartTime.Add(time.Duration(consultation.Duration) * time.Minute)
	if time.Now().After(endTime) {
		return nil, errors.New("consultation time has ended")
	}

	chat.CreatedAt = time.Now()
	err = uc.ChatRepo.SaveChat(&chat)
	if err != nil {
		return nil, errors.New("failed to save chat")
	}

	var senderName string
	if chat.SenderType == "user" {
		user, _ := uc.ChatRepo.GetUserByID(chat.SenderID)
		senderName = user.Username
	} else if chat.SenderType == "doctor" {
		doctor, _ := uc.ChatRepo.GetDoctorByID(chat.SenderID)
		senderName = doctor.Username
	}

	return &model.ChatDTO{
		ID:             chat.ID,
		ConsultationID: chat.ConsultationID,
		SenderID:       chat.SenderID,
		SenderName:     senderName,
		Message:        chat.Message,
		SenderType:     chat.SenderType,
		CreatedAt:      chat.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (uc *ChatUsecaseImpl) GetChatHistory(consultationID int) ([]model.ChatDTO, error) {
	chats, err := uc.ChatRepo.GetChatHistory(consultationID)
	if err != nil {
		return nil, err
	}

	var chatDTOs []model.ChatDTO
	for _, chat := range chats {
		var senderName string
		if chat.SenderType == "user" && chat.User != nil {
			senderName = chat.User.Username
		} else if chat.SenderType == "doctor" && chat.Doctor != nil {
			senderName = chat.Doctor.Username
		}

		chatDTOs = append(chatDTOs, model.ChatDTO{
			ID:             chat.ID,
			ConsultationID: chat.ConsultationID,
			SenderID:       chat.SenderID,
			SenderName:     senderName,
			Message:        chat.Message,
			SenderType:     chat.SenderType,
			CreatedAt:      chat.CreatedAt.Format(time.RFC3339),
		})
	}
	return chatDTOs, nil
}
