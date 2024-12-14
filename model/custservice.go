package model

import "time"

type CustService struct {
	ID         int       `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID     int       `json:"user_id"`
	Message    string    `json:"message"`
	Answer     string    `json:"answer"`
	IsAnswered bool      `json:"is_answered"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
	AnsweredAt time.Time `json:"answered_at" gorm:"autoUpdateTime"`
}
