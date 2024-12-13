package usecase

import (
	"calmind/model"
	"calmind/repository"
	"errors"
	"fmt"
	"time"
)

// ChatUsecase defines the interface for chat-related operations.
type ChatUsecase interface {
	GetOrCreateRoom(userID int, doctorID int) (string, error)
	ValidateChatAccess(userID int, doctorID int) error
	SendChat(chat model.Chat) (*model.ChatDTO, error)
	GetChatHistory(roomID string) ([]model.ChatDTO, error)
}

// ChatUsecaseImpl is the implementation of ChatUsecase.
type ChatUsecaseImpl struct {
	ChatRepo         repository.ChatRepository
	ConsultationRepo repository.ConsultationRepository
}

// NewChatUsecaseImpl creates a new instance of ChatUsecaseImpl.
func NewChatUsecaseImpl(chatRepo repository.ChatRepository, consultationRepo repository.ConsultationRepository) *ChatUsecaseImpl {
	return &ChatUsecaseImpl{
		ChatRepo:         chatRepo,
		ConsultationRepo: consultationRepo,
	}
}

// GetOrCreateRoom retrieves or creates a chat room between a user and a doctor.
func (uc *ChatUsecaseImpl) GetOrCreateRoom(userID int, doctorID int) (string, error) {
	// Check if the chat room already exists
	room, err := uc.ChatRepo.GetRoomByUserAndDoctor(userID, doctorID)
	if err == nil && room != nil {
		return room.RoomID, nil
	}

	// Create a new chat room
	roomID := fmt.Sprintf("room-%d-%d", userID, doctorID)
	err = uc.ChatRepo.CreateRoom(roomID, userID, doctorID)
	if err != nil {
		return "", fmt.Errorf("failed to create chat room: %w", err)
	}
	return roomID, nil
}

// ValidateChatAccess checks if the user and doctor have an active consultation.
func (uc *ChatUsecaseImpl) ValidateChatAccess(userID int, doctorID int) error {
	consultations, err := uc.ConsultationRepo.GetActiveConsultations(userID, doctorID)
	if err != nil {
		return fmt.Errorf("failed to check consultation status: %w", err)
	}

	if len(consultations) == 0 {
		return errors.New("no active consultation found, please complete a payment to start chatting")
	}

	return nil
}

// SendChat handles the sending of chat messages between a user and a doctor.
func (uc *ChatUsecaseImpl) SendChat(chat model.Chat) (*model.ChatDTO, error) {
	// Validate access to chat
	err := uc.ValidateChatAccess(chat.UserID, chat.DoctorID)
	if err != nil {
		return nil, err
	}

	// Get or create the chat room
	roomID, err := uc.GetOrCreateRoom(chat.UserID, chat.DoctorID)
	if err != nil {
		return nil, err
	}

	// Assign the RoomID to the chat and set the timestamp
	chat.RoomID = roomID
	chat.CreatedAt = time.Now()

	// Save the chat message
	err = uc.ChatRepo.SaveChat(&chat)
	if err != nil {
		return nil, fmt.Errorf("failed to save chat message: %w", err)
	}

	// Convert to ChatDTO and return
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

// GetChatHistory retrieves the chat history for a given RoomID.
func (uc *ChatUsecaseImpl) GetChatHistory(roomID string) ([]model.ChatDTO, error) {
	// Fetch chat messages by RoomID
	chats, err := uc.ChatRepo.GetChatsByRoomID(roomID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch chat history: %w", err)
	}

	// Convert to DTOs for the response
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
