package model

import "time"

type Consultation struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    int       `json:"user_id" gorm:"not null"`
	DoctorID  int       `json:"doctor_id" gorm:"not null"`
	Message   string    `json:"message" gorm:"type:text;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"` // otomatis saat data dibuat
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"` // otomatis saat data diupdate
}
