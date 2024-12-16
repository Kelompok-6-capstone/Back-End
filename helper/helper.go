package helper

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// WebSocket upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow any origin for simplicity; refine for production.
	},
}

// Client struct to hold connection and user details
type Client struct {
	ID   string
	Name string
	Conn *websocket.Conn
}

// ChatSession to hold exactly two clients
type ChatSession struct {
	User1 *Client
	User2 *Client
}

var sessions = make(map[string]*ChatSession) // Mapping session ID to ChatSession
var sessionsLock sync.Mutex

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	// Parse query params for user details
	userID := r.URL.Query().Get("id")
	userName := r.URL.Query().Get("name")
	sessionID := r.URL.Query().Get("session")

	if userID == "" || userName == "" || sessionID == "" {
		http.Error(w, "Missing required query parameters", http.StatusBadRequest)
		return
	}

	// Upgrade to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}
	defer conn.Close()

	client := &Client{ID: userID, Name: userName, Conn: conn}

	// Add client to the session
	session := addClientToSession(sessionID, client)

	log.Printf("User %s joined session %s", userName, sessionID)

	// Listen for incoming messages
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("User %s disconnected: %v", userName, err)
			break
		}

		log.Printf("Message from %s: %s", userName, message)

		// Broadcast message to the other user in the session
		broadcastMessage(session, client, message)
	}

	// Remove client on disconnect
	removeClientFromSession(sessionID, client)
}

func addClientToSession(sessionID string, client *Client) *ChatSession {
	sessionsLock.Lock()
	defer sessionsLock.Unlock()

	session, exists := sessions[sessionID]
	if !exists {
		session = &ChatSession{}
		sessions[sessionID] = session
	}

	if session.User1 == nil {
		session.User1 = client
	} else if session.User2 == nil {
		session.User2 = client
	} else {
		log.Printf("Session %s is full", sessionID)
		client.Conn.Close()
	}

	return session
}

func removeClientFromSession(sessionID string, client *Client) {
	sessionsLock.Lock()
	defer sessionsLock.Unlock()

	session := sessions[sessionID]
	if session != nil {
		if session.User1 == client {
			session.User1 = nil
		} else if session.User2 == client {
			session.User2 = nil
		}

		// Delete session if empty
		if session.User1 == nil && session.User2 == nil {
			delete(sessions, sessionID)
		}
	}
}

func broadcastMessage(session *ChatSession, sender *Client, message []byte) {
	if session.User1 != nil && session.User1 != sender {
		session.User1.Conn.WriteMessage(websocket.TextMessage, message)
	}
	if session.User2 != nil && session.User2 != sender {
		session.User2.Conn.WriteMessage(websocket.TextMessage, message)
	}
}
