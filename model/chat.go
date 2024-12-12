package model

import "time"

type Chat struct {
	ID             int       `json:"id" gorm:"primaryKey"`
	ConsultationID int       `json:"consultation_id"`
	SenderID       int       `json:"sender_id"`
	Message        string    `json:"message"`
	SenderType     string    `json:"sender_type"` // "user" atau "doctor"
	CreatedAt      time.Time `json:"created_at"`

	User   *User   `gorm:"foreignKey:SenderID;references:ID"`
	Doctor *Doctor `gorm:"foreignKey:SenderID;references:ID"`
}

type ChatDTO struct {
	ID             int    `json:"id"`
	ConsultationID int    `json:"consultation_id"`
	SenderID       int    `json:"sender_id"`
	SenderName     string `json:"sender_name"`
	Message        string `json:"message"`
	SenderType     string `json:"sender_type"`
	CreatedAt      string `json:"created_at"`
}
