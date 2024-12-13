package usecase

import (
	"calmind/model"
	"calmind/repository"
	"errors"
	"fmt"
	"time"
)

type ChatUsecase interface {
	ValidateChatAccess(userID int, receiverID int) error
	SendChat(chat model.Chat) (*model.ChatDTO, error)
	GetChatHistory(roomID string) ([]model.ChatDTO, error)
}

type ChatUsecaseImpl struct {
	ChatRepo         repository.ChatRepository
	ConsultationRepo repository.ConsultationRepository
}

func NewChatUsecase(chatRepo repository.ChatRepository, ConsultationRepo repository.ConsultationRepository) *ChatUsecaseImpl {
	return &ChatUsecaseImpl{
		ChatRepo:         chatRepo,
		ConsultationRepo: ConsultationRepo,
	}
}

func (uc *ChatUsecaseImpl) ValidateChatAccess(userID int, receiverID int) error {
	// Validate if there is an active consultation between the user and the doctor.
	activeConsultations, err := uc.ConsultationRepo.GetActiveConsultations(userID, receiverID)
	if err != nil || len(activeConsultations) == 0 {
		return errors.New("no active consultation found, please complete a payment to start chatting")
	}
	return nil
}

func (uc *ChatUsecaseImpl) SendChat(chat model.Chat) (*model.ChatDTO, error) {
	// Generate RoomID if it doesn't exist.
	roomID := fmt.Sprintf("room-%d-%d", chat.UserID, chat.DoctorID)
	chat.RoomID = roomID
	chat.CreatedAt = time.Now()

	// Save the chat to the database.
	err := uc.ChatRepo.SaveChat(&chat)
	if err != nil {
		return nil, fmt.Errorf("failed to save chat: %v", err)
	}

	// Return the chat data as DTO.
	return &model.ChatDTO{
		ID:         chat.ID,
		RoomID:     chat.RoomID,
		UserID:     chat.UserID,
		DoctorID:   chat.DoctorID,
		SenderID:   chat.SenderID,
		Message:    chat.Message,
		SenderType: chat.SenderType,
		CreatedAt:  chat.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (uc *ChatUsecaseImpl) GetChatHistory(roomID string) ([]model.ChatDTO, error) {
	chats, err := uc.ChatRepo.GetChatsByRoomID(roomID)
	if err != nil {
		return nil, err
	}

	var chatDTOs []model.ChatDTO
	for _, chat := range chats {
		chatDTOs = append(chatDTOs, model.ChatDTO{
			ID:         chat.ID,
			RoomID:     chat.RoomID,
			UserID:     chat.UserID,
			DoctorID:   chat.DoctorID,
			SenderID:   chat.SenderID,
			Message:    chat.Message,
			SenderType: chat.SenderType,
			CreatedAt:  chat.CreatedAt.Format(time.RFC3339),
		})
	}
	return chatDTOs, nil
}
