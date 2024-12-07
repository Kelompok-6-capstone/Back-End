package model

import "time"

type Payment struct {
	ID             int       `json:"id" gorm:"primaryKey;autoIncrement"`
	ConsultationID int       `json:"consultation_id" gorm:"not null"`
	Amount         int       `json:"amount" gorm:"not null"`
	Status         string    `json:"status" gorm:"default:'pending'"`
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
