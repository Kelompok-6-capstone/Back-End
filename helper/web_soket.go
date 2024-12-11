package helper

import (
	"calmind/model"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Channel untuk broadcast artikel
var Broadcast = make(chan model.Artikel)

// WebSocket upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Thread-safe map untuk menyimpan klien WebSocket
var clients sync.Map

// Handle koneksi WebSocket baru
func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade WebSocket: %v", err)
		http.Error(w, "Gagal meng-upgrade koneksi WebSocket", http.StatusInternalServerError)
		return
	}
	defer ws.Close()

	clients.Store(ws, true)

	// Heartbeat untuk memastikan koneksi aktif
	ws.SetReadDeadline(time.Now().Add(60 * time.Second))
	ws.SetPongHandler(func(string) error {
		ws.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			if err := ws.WriteMessage(websocket.PingMessage, nil); err != nil {
				clients.Delete(ws)
				ws.Close()
				break
			}
		}
	}()

	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			clients.Delete(ws)
			log.Println("WebSocket read error:", err)
			break
		}
	}
}

// Broadcast artikel ke semua klien WebSocket
func BroadcastMessages() {
	for {
		artikel := <-Broadcast
		clients.Range(func(key, value interface{}) bool {
			ws := key.(*websocket.Conn)
			err := ws.WriteJSON(artikel)
			if err != nil {
				ws.Close()
				clients.Delete(ws)
				fmt.Println("WebSocket write error:", err)
			}
			return true
		})
	}
}
