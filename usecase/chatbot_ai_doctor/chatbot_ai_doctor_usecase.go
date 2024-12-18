package usecase

import (
	"calmind/helper"
	"calmind/model"
	repository "calmind/repository/chatbot_ai_doctor"
	"context"
	"fmt"
	"strings"
)

type DoctorChatbotUsecase interface {
	GenerateDoctorRecommendation(doctorID int, message string) (string, error)
}

type DoctorChatbotUsecaseImpl struct {
	ChatbotRepo repository.DoctorChatbotRepository
}

func NewDoctorChatbotUsecase(chatbotRepo repository.DoctorChatbotRepository) DoctorChatbotUsecase {
	return &DoctorChatbotUsecaseImpl{ChatbotRepo: chatbotRepo}
}

func (u *DoctorChatbotUsecaseImpl) GenerateDoctorRecommendation(doctorID int, message string) (string, error) {
	// Validasi bahwa pesan harus berisi permintaan rekomendasi
	if !strings.Contains(message, "rekomendasi") {
		return "", fmt.Errorf("pesan tidak terkait dengan rekomendasi perawatan, hanya pertanyaan terkait rekomendasi perawatan yang bisa diproses")
	}

	// Prompt kustom untuk dokter
	customPrompt := "Anda adalah asisten AI yang membantu dokter dalam memberikan rekomendasi perawatan pasien. " +
		"Hanya berikan jawaban terkait rekomendasi perawatan berdasarkan permintaan dokter berikut: " + message

	// Ambil log percakapan dokter sebelumnya
	logs, err := u.ChatbotRepo.GetLogsByDoctorID(doctorID)
	if err != nil {
		return "", err
	}

	// Gabungkan percakapan sebelumnya
	var contextBuilder strings.Builder
	for _, log := range logs {
		contextBuilder.WriteString(log.Message + "\n")
		contextBuilder.WriteString(log.Response + "\n")
	}
	contextBuilder.WriteString(customPrompt + "\n")

	// Panggil AI untuk respons
	ctx := context.Background()
	response, err := helper.ResponseAI(ctx, contextBuilder.String())
	if err != nil {
		return "", err
	}

	// Simpan log percakapan baru
	newLog := model.Chatbot{
		UserID:   doctorID,
		Response: response,
	}
	if err := u.ChatbotRepo.SaveLog(&newLog); err != nil {
		return "", err
	}

	return response, nil
}

