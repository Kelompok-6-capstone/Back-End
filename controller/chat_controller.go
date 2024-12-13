package controller

import (
	"calmind/helper"
	"calmind/model"
	"calmind/service"
	"calmind/usecase"
	"net/http"

	"github.com/labstack/echo/v4"
)

// ChatController handles chat-related HTTP requests.
type ChatController struct {
	ChatUsecase usecase.ChatUsecase
}

// NewChatController creates a new instance of ChatController.
func NewChatController(chatUsecase usecase.ChatUsecase) *ChatController {
	return &ChatController{
		ChatUsecase: chatUsecase,
	}
}

// SendChat handles the sending of chat messages between a user and a doctor.
func (c *ChatController) SendChat(ctx echo.Context) error {
	// Retrieve claims from middleware (user or doctor).
	var claims *service.JwtCustomClaims
	var ok bool

	if ctx.Get("user") != nil {
		claims, ok = ctx.Get("user").(*service.JwtCustomClaims)
	} else if ctx.Get("doctor") != nil {
		claims, ok = ctx.Get("doctor").(*service.JwtCustomClaims)
	}

	if !ok || claims == nil {
		return helper.JSONErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized access.")
	}

	// Parse the request body.
	var request struct {
		DoctorID int    `json:"doctor_id"`
		Message  string `json:"message"`
	}

	if err := ctx.Bind(&request); err != nil || request.Message == "" || request.DoctorID <= 0 {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Invalid input.")
	}

	// Validate chat access for the user or doctor.
	err := c.ChatUsecase.ValidateChatAccess(claims.UserID, request.DoctorID)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusForbidden, err.Error())
	}

	// Determine the sender type (user or doctor).
	senderType := "user"
	if claims.Role == "doctor" {
		senderType = "doctor"
	}

	// Construct the chat message object.
	chat := model.Chat{
		UserID:     claims.UserID,
		DoctorID:   request.DoctorID,
		SenderID:   claims.UserID,
		Message:    request.Message,
		SenderType: senderType,
	}

	// Send the chat message via the usecase.
	chatDTO, err := c.ChatUsecase.SendChat(chat)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	// Return the chat message as a response.
	return helper.JSONSuccessResponse(ctx, chatDTO)
}

// GetChatHistory retrieves the chat history for a given RoomID.
func (c *ChatController) GetChatHistory(ctx echo.Context) error {
	// Get RoomID from the URL parameter.
	roomID := ctx.Param("room_id")
	if roomID == "" {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Room ID is required.")
	}

	// Retrieve the chat history via the usecase.
	chats, err := c.ChatUsecase.GetChatHistory(roomID)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve chat history.")
	}

	// Return the chat history as a response.
	return helper.JSONSuccessResponse(ctx, chats)
}
