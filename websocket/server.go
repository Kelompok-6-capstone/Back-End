package websocket

import (
	"fmt"
	"log"
	"net/http"
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
	userID := c.QueryParam("user_id")
	doctorID := c.QueryParam("doctor_id")

	if userID == "" || doctorID == "" {
		log.Println("Missing user_id or doctor_id")
		return echo.NewHTTPError(http.StatusBadRequest, "user_id and doctor_id are required")
	}

	ws, err := s.Upgrader.Upgrade(c.Response().Writer, c.Request(), nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return err
	}
	defer ws.Close()

	roomID := fmt.Sprintf("%s-%s", userID, doctorID)
	log.Printf("WebSocket connection established for room %s", roomID)

	s.Mutex.Lock()
	s.Clients[roomID] = ws
	s.Mutex.Unlock()

	// Heartbeat interval
	go s.heartbeat(ws, roomID)

	for {
		var msg model.ChatMessage
		ws.SetReadDeadline(time.Now().Add(30 * time.Second)) // Timeout 30 detik
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("Disconnected from room %s: %v", roomID, err)
			s.Mutex.Lock()
			delete(s.Clients, roomID)
			s.Mutex.Unlock()
			break
		}

		// Handle heartbeat messages
		if msg.Message == "ping" {
			ws.WriteJSON(map[string]string{
				"message": "pong",
			})
			continue
		}

		// Log and send message
		log.Printf("Received message: %+v", msg)
		err = s.ChatUsecase.SendMessage(msg.UserID, msg.DoctorID, msg.SenderID, msg.Message)
		if err != nil {
			log.Printf("Error sending message: %v", err)
			ws.WriteJSON(map[string]string{
				"error": "Failed to send message: " + err.Error(),
			})
			continue
		}

		// Broadcast message
		s.BroadcastMessage(roomID, msg)
	}
	return nil
}

func (s *WebSocketServer) BroadcastMessage(roomID string, msg model.ChatMessage) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	if client, ok := s.Clients[roomID]; ok {
		client.WriteJSON(msg)
	}
}

func (s *WebSocketServer) heartbeat(ws *websocket.Conn, roomID string) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		<-ticker.C
		s.Mutex.Lock()
		conn, ok := s.Clients[roomID]
		s.Mutex.Unlock()

		if !ok {
			return
		}

		// Send ping
		if err := conn.WriteJSON(map[string]string{"message": "ping"}); err != nil {
			log.Printf("Heartbeat failed for room %s: %v", roomID, err)
			s.Mutex.Lock()
			delete(s.Clients, roomID)
			s.Mutex.Unlock()
			conn.Close()
			return
		}
	}
}
