package controller

import (
	"calmind/helper"
	"calmind/service"
	"calmind/usecase"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type ChatbotController struct {
	ChatbotUsecase usecase.ChatbotUsecase
}

func NewChatbotController(chatbotUsecase usecase.ChatbotUsecase) *ChatbotController {
	return &ChatbotController{ChatbotUsecase: chatbotUsecase}
}

func (c *ChatbotController) GetChatResponse(ctx echo.Context) error {
	claims, _ := ctx.Get("user").(*service.JwtCustomClaims)

	var request struct {
		Message string `json:"message"`
	}

	if err := ctx.Bind(&request); err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Input tidak valid")
	}

	// Respon statis untuk pesan tertentu
	standardResponses := map[string]string{
		"hallo":        "Selamat datang di Calmind! Apakah ada yang bisa kami bantu?",
		"terima kasih": "Sama-sama! Jika ada pertanyaan lain, jangan ragu untuk bertanya.",
	}

	// Konversi pesan user ke huruf kecil
	lowerMessage := strings.ToLower(request.Message)

	// Jika pesan cocok dengan respon statis
	if response, exists := standardResponses[lowerMessage]; exists {
		return helper.JSONSuccessResponse(ctx, map[string]interface{}{
			"message":  request.Message,
			"response": response,
		})
	}

	// Jika tidak cocok, gunakan AI untuk menjawab
	response, err := c.ChatbotUsecase.GenerateResponse(claims.UserID, request.Message)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal memproses permintaan: "+err.Error())
	}

	// Tambahkan ucapan terima kasih jika pesan mengandung kata "terima kasih"
	if strings.Contains(lowerMessage, "terima kasih") {
		response += " Sama-sama! Jangan ragu untuk bertanya lagi."
	}

	return helper.JSONSuccessResponse(ctx, map[string]interface{}{
		"message":  request.Message,
		"response": response,
	})
}
