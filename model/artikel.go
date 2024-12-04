package model

import "time"

type Artikel struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement"`
	AdminID   int       `json:"admin_id" gorm:"not null"` // Foreign key
	Admin     Admin     `json:"admin" gorm:"foreignKey:AdminID"`
	Judul     string    `json:"judul" gorm:"not null"`
	Gambar    string    `json:"gambar"`
	Isi       string    `json:"isi" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
