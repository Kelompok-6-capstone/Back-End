package usecase

import (
	"calmind/helper"
	"calmind/model"
	"calmind/repository"
	"context"
	"strings"
)

type ChatbotUsecase interface {
	GenerateResponse(userID int, message string) (string, error)
}

type ChatbotUsecaseImpl struct {
	ChatLogRepo repository.ChatLogRepository
}

func NewChatbotUsecase(chatLogRepo repository.ChatLogRepository) ChatbotUsecase {
	return &ChatbotUsecaseImpl{ChatLogRepo: chatLogRepo}
}

func (u *ChatbotUsecaseImpl) GenerateResponse(userID int, message string) (string, error) {
	// Ambil log percakapan sebelumnya
	logs, err := u.ChatLogRepo.GetLogsByUserID(userID)
	if err != nil {
		return "", err
	}

	// Gabungkan percakapan sebelumnya
	var contextBuilder strings.Builder
	for _, log := range logs {
		contextBuilder.WriteString(log.Message + "\n")
		contextBuilder.WriteString(log.Response + "\n")
	}
	contextBuilder.WriteString(message + "\n")

	// Panggil AI untuk respons
	ctx := context.Background()
	response, err := helper.ResponseAI(ctx, contextBuilder.String())
	if err != nil {
		return "", err
	}

	// Simpan log percakapan baru
	newLog := model.ChatLog{
		UserID:   userID,
		Message:  message,
		Response: response,
	}
	if err := u.ChatLogRepo.SaveLog(&newLog); err != nil {
		return "", err
	}

	return response, nil
}
