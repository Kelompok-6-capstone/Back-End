package model

import "time"

type Consultation struct {
	ID          int       `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID      int       `json:"user_id" gorm:"not null"`
	DoctorID    int       `json:"doctor_id" gorm:"not null"`
	Message     string    `json:"message" gorm:"type:text;not null"`
	Rekomendasi string    `json:"rekomendasi" gorm:"type:text"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	User        User      `json:"user" gorm:"foreignKey:UserID"`
}
