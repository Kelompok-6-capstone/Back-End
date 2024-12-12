package controller

import (
	"calmind/helper"
	"calmind/model"
	"calmind/service"
	"calmind/usecase"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type ChatController struct {
	ChatUsecase  *usecase.ChatUsecaseImpl
	WebSocketHub *helper.WebSocketHub
}

func NewChatController(chatUsecase *usecase.ChatUsecaseImpl, webSocketHub *helper.WebSocketHub) *ChatController {
	return &ChatController{
		ChatUsecase:  chatUsecase,
		WebSocketHub: webSocketHub,
	}
}

// Mengirim pesan
func (c *ChatController) SendChat(ctx echo.Context) error {
	claims, ok := ctx.Get("user").(*service.JwtCustomClaims)
	if !ok || claims == nil {
		return helper.JSONErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized access.")
	}

	var request struct {
		ConsultationID int    `json:"consultation_id"`
		Message        string `json:"message"`
	}

	if err := ctx.Bind(&request); err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Invalid input.")
	}

	senderType := "user"
	if claims.Role == "doctor" {
		senderType = "doctor"
	}

	chat := model.Chat{
		ConsultationID: request.ConsultationID,
		SenderID:       claims.UserID,
		Message:        request.Message,
		SenderType:     senderType,
	}

	chatDTO, err := c.ChatUsecase.SendChat(chat)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	// Broadcast ke klien yang relevan berdasarkan ConsultationID
	chatBytes, _ := json.Marshal(chatDTO)
	c.WebSocketHub.BroadcastToConsultation(request.ConsultationID, chatBytes)

	return helper.JSONSuccessResponse(ctx, chatDTO)
}

// Mendapatkan riwayat pesan
func (c *ChatController) GetChatHistory(ctx echo.Context) error {
	claims, ok := ctx.Get("user").(*service.JwtCustomClaims)
	if !ok || claims == nil {
		return helper.JSONErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized access.")
	}

	consultationID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Invalid consultation ID.")
	}

	chats, err := c.ChatUsecase.GetChatHistory(consultationID)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve chat history.")
	}

	return helper.JSONSuccessResponse(ctx, chats)
}

// WebSocket handler
func (c *ChatController) WebSocketHandler(ctx echo.Context) error {
	c.WebSocketHub.HandleConnections(ctx.Response(), ctx.Request())
	return nil
}
