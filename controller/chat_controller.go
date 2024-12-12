package controller

import (
	"calmind/helper"
	"calmind/model"
	"calmind/service"
	"calmind/usecase"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ChatController struct {
	ChatUsecase usecase.ChatUsecase
}

func NewChatController(chatUsecase usecase.ChatUsecase) *ChatController {
	return &ChatController{
		ChatUsecase: chatUsecase,
	}
}

func (c *ChatController) SendChat(ctx echo.Context) error {
	var claims *service.JwtCustomClaims
	var ok bool

	// Ambil klaim berdasarkan middleware
	if ctx.Get("user") != nil {
		claims, ok = ctx.Get("user").(*service.JwtCustomClaims)
	} else if ctx.Get("doctor") != nil {
		claims, ok = ctx.Get("doctor").(*service.JwtCustomClaims)
	}

	if !ok || claims == nil {
		return helper.JSONErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized access.")
	}

	// Bind input dari body request
	var request struct {
		Message string `json:"message"`
	}

	if err := ctx.Bind(&request); err != nil || request.Message == "" {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Invalid input.")
	}

	// Tentukan tipe pengirim
	senderType := "user"
	if claims.Role == "doctor" {
		senderType = "doctor"
	}

	// Buat objek chat
	chat := model.Chat{
		UserID:     claims.UserID,
		SenderID:   claims.UserID,
		Message:    request.Message,
		SenderType: senderType,
	}

	// Simpan chat
	chatDTO, err := c.ChatUsecase.SendChat(chat)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	return helper.JSONSuccessResponse(ctx, chatDTO)
}

func (c *ChatController) GetChatHistory(ctx echo.Context) error {
	var claims *service.JwtCustomClaims
	var ok bool

	// Ambil klaim berdasarkan middleware
	if ctx.Get("user") != nil {
		claims, ok = ctx.Get("user").(*service.JwtCustomClaims)
	} else if ctx.Get("doctor") != nil {
		claims, ok = ctx.Get("doctor").(*service.JwtCustomClaims)
	}

	if !ok || claims == nil {
		return helper.JSONErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized access.")
	}

	// Ambil riwayat chat
	chats, err := c.ChatUsecase.GetChatHistory(claims.UserID)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve chat history.")
	}

	return helper.JSONSuccessResponse(ctx, chats)
}
