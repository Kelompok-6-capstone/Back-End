package model

import "time"

type Chat struct {
	ID         int       `json:"id" gorm:"primaryKey"`             // ID chat sebagai primary key
	UserID     int       `json:"user_id" gorm:"not null"`          // ID pengguna untuk menghubungkan chat ke pengguna
	SenderID   int       `json:"sender_id" gorm:"not null"`        // ID pengirim (bisa user atau doctor)
	Message    string    `json:"message" gorm:"type:text"`         // Isi pesan chat
	SenderType string    `json:"sender_type" gorm:"not null"`      // "user" atau "doctor"
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"` // Waktu pembuatan chat

	// Relasi dengan model User dan Doctor
	User   *User   `gorm:"foreignKey:UserID;references:ID"`   // Relasi ke User berdasarkan UserID
	Doctor *Doctor `gorm:"foreignKey:SenderID;references:ID"` // Relasi ke Doctor berdasarkan SenderID
}

type ChatDTO struct {
	ID         int    `json:"id"`          // ID chat
	UserID     int    `json:"user_id"`     // ID pengguna
	SenderID   int    `json:"sender_id"`   // ID pengirim
	SenderName string `json:"sender_name"` // Nama pengirim (diambil dari User atau Doctor)
	Message    string `json:"message"`     // Isi pesan
	SenderType string `json:"sender_type"` // "user" atau "doctor"
	CreatedAt  string `json:"created_at"`  // Waktu pembuatan dalam format string
}
