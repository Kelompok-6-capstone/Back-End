package model

import "time"

type Artikel struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`                          // Primary key
	AdminID   int       `gorm:"not null;index" json:"admin_id"`                              // Foreign key dengan index
	Admin     Admin     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"admin"` // Relasi ke Admin
	Judul     string    `gorm:"type:varchar(255);not null" json:"judul"`                     // Panjang maksimum judul ditentukan
	Gambar    string    `json:"gambar"`
	Isi       string    `gorm:"not null" json:"isi"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
