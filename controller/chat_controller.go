package controller

import (
	"calmind/helper"
	"calmind/model"
	"calmind/service"
	"calmind/usecase"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type ChatController struct {
	ChatUsecase  *usecase.ChatUsecaseImpl
	WebSocketHub *helper.WebSocketHub
}

// Tambahkan upgrader untuk WebSocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Mengizinkan semua origin (dapat disesuaikan)
		return true
	},
}

func NewChatController(chatUsecase *usecase.ChatUsecaseImpl, webSocketHub *helper.WebSocketHub) *ChatController {
	return &ChatController{
		ChatUsecase:  chatUsecase,
		WebSocketHub: webSocketHub,
	}
}

// Mengirim pesan
func (c *ChatController) SendChat(ctx echo.Context) error {
	// Ambil klaim berdasarkan middleware
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

	if claims.Role != "user" && claims.Role != "doctor" {
		return helper.JSONErrorResponse(ctx, http.StatusForbidden, "Access denied.")
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

func (c *ChatController) GetChatHistory(ctx echo.Context) error {
	// Ambil klaim berdasarkan middleware
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
	// Upgrade connection to WebSocket
	ws, err := upgrader.Upgrade(ctx.Response(), ctx.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	// Lakukan sesuatu dengan WebSocket connection
	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			log.Println("WebSocket read error:", err)
			return err
		}
		log.Printf("Received message: %s", msg)

		// Kirim kembali pesan ke klien
		if err := ws.WriteMessage(websocket.TextMessage, msg); err != nil {
			log.Println("WebSocket write error:", err)
			return err
		}
	}
}
