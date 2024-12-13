package model

import "time"

type Otp struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	Email     string    `json:"email"`
	Code      string    `json:"code"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"` // otomatis saat data dibuat
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"` // otomatis saat data diupdate
}
