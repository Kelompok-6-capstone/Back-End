package helper

import (
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID             int
	Conn           *websocket.Conn
	Send           chan []byte
	ConsultationID int
}

type WebSocketHub struct {
	Clients    map[int][]*Client
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan []byte
	mu         sync.Mutex
}

func NewWebSocketHub() *WebSocketHub {
	return &WebSocketHub{
		Clients:    make(map[int][]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan []byte),
	}
}

func (hub *WebSocketHub) Run() {
	for {
		select {
		case client := <-hub.Register:
			hub.mu.Lock()
			hub.Clients[client.ConsultationID] = append(hub.Clients[client.ConsultationID], client)
			hub.mu.Unlock()
		case client := <-hub.Unregister:
			hub.mu.Lock()
			clients := hub.Clients[client.ConsultationID]
			for i, c := range clients {
				if c == client {
					hub.Clients[client.ConsultationID] = append(clients[:i], clients[i+1:]...)
					break
				}
			}
			hub.mu.Unlock()
		}
	}
}

// HandleConnections menangani koneksi WebSocket
func (hub *WebSocketHub) HandleConnections(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}

	consultationIDStr := r.URL.Query().Get("consultation_id")
	consultationID, err := strconv.Atoi(consultationIDStr)
	if err != nil {
		log.Println("Invalid consultation_id:", consultationIDStr)
		conn.Close()
		return
	}

	client := &Client{
		Conn:           conn,
		Send:           make(chan []byte),
		ConsultationID: consultationID,
	}

	hub.Register <- client

	go client.ReadMessages(hub)
	go client.WriteMessages()
}

// BroadcastToConsultation mengirim pesan ke klien berdasarkan ConsultationID
func (hub *WebSocketHub) BroadcastToConsultation(consultationID int, message []byte) {
	hub.mu.Lock()
	defer hub.mu.Unlock()

	clients := hub.Clients[consultationID]
	for _, client := range clients {
		select {
		case client.Send <- message:
		default:
			log.Println("Failed to send message, removing client")
			close(client.Send)
			hub.Unregister <- client
		}
	}
}

func (c *Client) ReadMessages(hub *WebSocketHub) {
	defer func() {
		hub.Unregister <- c
		c.Conn.Close()
	}()
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}
		hub.Broadcast <- message
	}
}

func (c *Client) WriteMessages() {
	for message := range c.Send {
		if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Println("Error writing message:", err)
			break
		}
	}
}
