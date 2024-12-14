package websocket

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"calmind/model"
	"calmind/usecase"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type WebSocketServer struct {
	ChatUsecase usecase.ChatUsecase
	Clients     map[string]*websocket.Conn // Key: user-doctor pair
	Mutex       *sync.Mutex
	Upgrader    websocket.Upgrader
}

func NewWebSocketServer(chatUsecase usecase.ChatUsecase) *WebSocketServer {
	return &WebSocketServer{
		ChatUsecase: chatUsecase,
		Clients:     make(map[string]*websocket.Conn),
		Mutex:       &sync.Mutex{},
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins for now
			},
		},
	}
}

func (s *WebSocketServer) HandleWebSocket(c echo.Context) error {
	userIDStr := c.QueryParam("user_id")
	doctorIDStr := c.QueryParam("doctor_id")

	if userIDStr == "" || doctorIDStr == "" {
		log.Println("Missing user_id or doctor_id")
		return echo.NewHTTPError(http.StatusBadRequest, "user_id and doctor_id are required")
	}

	// Convert user_id and doctor_id to integers
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		log.Printf("Invalid user_id format: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "user_id must be a valid integer")
	}

	doctorID, err := strconv.Atoi(doctorIDStr)
	if err != nil {
		log.Printf("Invalid doctor_id format: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "doctor_id must be a valid integer")
	}

	ws, err := s.Upgrader.Upgrade(c.Response().Writer, c.Request(), nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return err
	}
	defer func() {
		log.Printf("Closing WebSocket connection for room %s-%s", userIDStr, doctorIDStr)
		ws.Close()
	}()

	roomID := fmt.Sprintf("%d-%d", userID, doctorID)
	log.Printf("WebSocket connection established for room %s", roomID)

	// Register client connection
	s.Mutex.Lock()
	s.Clients[roomID] = ws
	s.Mutex.Unlock()

	// Heartbeat handler
	go s.handleHeartbeat(ws, roomID)

	// Read messages
	for {
		var msg model.ChatMessage
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("Disconnected from room %s: %v", roomID, err)
			s.Mutex.Lock()
			delete(s.Clients, roomID)
			s.Mutex.Unlock()
			break
		}

		// Log received message
		log.Printf("Received message: %+v", msg)

		// Validate message fields
		if msg.UserID == 0 || msg.DoctorID == 0 || msg.Message == "" {
			log.Printf("Invalid message data: %+v", msg)
			ws.WriteJSON(map[string]string{"error": "Invalid message data"})
			continue
		}

		// Process and send the message using usecase
		err = s.ChatUsecase.SendMessage(msg.UserID, msg.DoctorID, msg.SenderID, msg.Message)
		if err != nil {
			log.Printf("Error sending message: %v", err)
			ws.WriteJSON(map[string]string{
				"error": "Failed to send message: " + err.Error(),
			})
			continue
		}

		// Broadcast the message
		s.BroadcastMessage(roomID, msg)
	}
	return nil
}

func (s *WebSocketServer) BroadcastMessage(roomID string, msg model.ChatMessage) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	if client, ok := s.Clients[roomID]; ok {
		err := client.WriteJSON(msg)
		if err != nil {
			log.Printf("Error broadcasting message to room %s: %v", roomID, err)
		}
	}
}

func (s *WebSocketServer) handleHeartbeat(ws *websocket.Conn, roomID string) {
	for {
		err := ws.WriteMessage(websocket.PingMessage, nil)
		if err != nil {
			log.Printf("Heartbeat failed for room %s: %v", roomID, err)
			s.Mutex.Lock()
			delete(s.Clients, roomID)
			s.Mutex.Unlock()
			break
		}
		log.Printf("Heartbeat sent to room %s", roomID)
		// Delay between heartbeats
		time.Sleep(30 * time.Second)
	}
}
