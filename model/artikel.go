package model

import "time"

type Artikel struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	AdminID   int       `gorm:"not null" json:"admin_id"`                       // Foreign key
	Admin     Admin     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"` // Relasi ke Admin
	Judul     string    `gorm:"not null" json:"judul"`
	Gambar    string    `json:"gambar"`
	Isi       string    `gorm:"not null" json:"isi"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
