package controller

import (
	"calmind/helper"
	"calmind/service"
	usecase "calmind/usecase/chatbot_ai_doctor"
	"strings"

	"net/http"

	"github.com/labstack/echo/v4"
)

type DoctorChatbotController struct {
	ChatbotUsecase usecase.DoctorChatbotUsecase
}

func NewDoctorChatbotController(chatbotUsecase usecase.DoctorChatbotUsecase) *DoctorChatbotController {
	return &DoctorChatbotController{ChatbotUsecase: chatbotUsecase}
}

func (c *DoctorChatbotController) GetDoctorRecommendation(ctx echo.Context) error {
	claims, _ := ctx.Get("doctor").(*service.JwtCustomClaims)

	var request struct {
		Message string `json:"message"`
	}

	if err := ctx.Bind(&request); err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "input tidak valid")
	}

	// Panggil usecase untuk mendapatkan rekomendasi
	response, err := c.ChatbotUsecase.GenerateDoctorRecommendation(claims.UserID, request.Message)
	if err != nil {
		if strings.Contains(err.Error(), "pesan tidak terkait dengan rekomendasi perawatan") {
			return helper.JSONErrorResponse(ctx, http.StatusBadRequest, err.Error()) // Menggunakan 400 Bad Request
		}
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "gagal memproses permintaan: "+err.Error())
	}

	return helper.JSONSuccessResponse(ctx, map[string]interface{}{
		"response": response,
	})
}

