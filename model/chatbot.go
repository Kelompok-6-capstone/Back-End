package model

import "time"

type ChatLog struct {
	ID        int       `gorm:"primaryKey;autoIncrement"`
	UserID    int       `gorm:"not null"`
	Message   string    `gorm:"type:text;not null"` // Pesan dari user
	Response  string    `gorm:"type:text;not null"` // Respons chatbot
	CreatedAt time.Time `gorm:"autoCreateTime"`     // Waktu percakapan
}
