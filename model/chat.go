package model

import "time"

type Chat struct {
	ID         int       `json:"id" gorm:"primaryKey"`
	RoomID     string    `json:"room_id" gorm:"not null"`   // Room chat unik
	UserID     int       `json:"user_id" gorm:"not null"`   // ID pengguna
	DoctorID   int       `json:"doctor_id" gorm:"not null"` // ID dokter
	SenderID   int       `json:"sender_id" gorm:"not null"` // ID pengirim
	Message    string    `json:"message" gorm:"type:text"`  // Isi pesan
	SenderType string    `json:"sender_type"`               // Jenis pengirim (user/doctor)
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
}

type ChatDTO struct {
	ID         int    `json:"id"`
	RoomID     string `json:"room_id"`
	UserID     int    `json:"user_id"`
	DoctorID   int    `json:"doctor_id"`
	SenderID   int    `json:"sender_id"`
	SenderName string `json:"sender_name"`
	Message    string `json:"message"`
	SenderType string `json:"sender_type"`
	CreatedAt  string `json:"created_at"`
}
