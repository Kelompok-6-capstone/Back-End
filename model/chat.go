package model

import "time"

type ChatMessage struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    int       `json:"user_id" gorm:"not null"`
	DoctorID  int       `json:"doctor_id" gorm:"not null"`
	SenderID  int       `json:"sender_id" gorm:"not null"`
	Message   string    `json:"message" gorm:"type:text"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}
