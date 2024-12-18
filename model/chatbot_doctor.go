package model

import "time"

type Chatbot struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	UserID    int       `json:"user_id"`    // ID dokter
	Message   string    `json:"message"`    // Pesan yang dikirim oleh dokter
	Response  string    `json:"response"`   // Respons dari AI
	CreatedAt time.Time `json:"created_at"` // Waktu percakapan dibuat
	UpdatedAt time.Time `json:"updated_at"` // Waktu percakapan diperbarui
}
