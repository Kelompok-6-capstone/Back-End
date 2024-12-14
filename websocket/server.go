package websocket

import (
	"fmt"
	"net/http"
	"sync"

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
		return echo.NewHTTPError(http.StatusBadRequest, "user_id and doctor_id are required")
	}

	ws, err := s.Upgrader.Upgrade(c.Response().Writer, c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	roomID := fmt.Sprintf("%s-%s", userID, doctorID)
	s.Mutex.Lock()
	s.Clients[roomID] = ws
	s.Mutex.Unlock()

	for {
		var msg model.ChatMessage
		err := ws.ReadJSON(&msg)
		if err != nil {
			fmt.Printf("Disconnected from room %s\n", roomID)
			s.Mutex.Lock()
			delete(s.Clients, roomID)
			s.Mutex.Unlock()
			break
		}

		// Validasi dan kirim pesan
		err = s.ChatUsecase.SendMessage(msg.UserID, msg.DoctorID, msg.SenderID, msg.Message)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		// Broadcast pesan
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
