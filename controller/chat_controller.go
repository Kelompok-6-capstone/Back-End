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

// SendChat handles sending chat messages between a user and a doctor.
func (c *ChatController) SendChat(ctx echo.Context) error {
	// Retrieve claims from middleware (user or doctor).
	var claims *service.JwtCustomClaims
	var ok bool

	if ctx.Get("user") != nil {
		claims, ok = ctx.Get("user").(*service.JwtCustomClaims)
	} else if ctx.Get("doctor") != nil {
		claims, ok = ctx.Get("doctor").(*service.JwtCustomClaims)
	} else {
		return helper.JSONErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized access.")
	}

	if !ok || claims == nil {
		return helper.JSONErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized access.")
	}

	// Parse the request body.
	var request struct {
		UserID   int    `json:"user_id,omitempty"`   // Required for doctors
		DoctorID int    `json:"doctor_id,omitempty"` // Required for users
		Message  string `json:"message"`
	}

	if err := ctx.Bind(&request); err != nil || request.Message == "" {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Invalid input.")
	}

	// Validate input based on the role.
	if claims.Role == "user" && request.DoctorID <= 0 {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Doctor ID is required for users.")
	}

	if claims.Role == "doctor" && request.UserID <= 0 {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "User ID is required for doctors.")
	}

	// Determine sender and receiver.
	senderID := claims.UserID
	receiverID := request.DoctorID // Default for user
	if claims.Role == "doctor" {
		senderID = claims.UserID
		receiverID = request.UserID
	}

	// Validate chat access.
	err := c.ChatUsecase.ValidateChatAccess(claims.UserID, receiverID)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusForbidden, err.Error())
	}

	// Build the chat message.
	chat := model.Chat{
		UserID:     receiverID,
		DoctorID:   senderID,
		SenderID:   senderID,
		Message:    request.Message,
		SenderType: claims.Role,
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
