package controller

import (
	"calmind/helper"
	"calmind/service"
	"calmind/usecase"
	"net/http"

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

	response, err := c.ChatbotUsecase.GenerateResponse(claims.UserID, request.Message)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal memproses permintaan: "+err.Error())
	}

	return helper.JSONSuccessResponse(ctx, map[string]interface{}{
		"message":  request.Message,
		"response": response,
	})
}
