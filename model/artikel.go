package model

import "time"

type Artikel struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	AdminID   int       `gorm:"not null" json:"admin_id"`
	Admin     Admin     `gorm:"foreignKey:AdminID" json:"admin"`
	Judul     string    `gorm:"not null" json:"judul"`
	DeleteURL string    `json:"delete_url"`
	Gambar    string    `json:"gambar"`
	Isi       string    `gorm:"not null" json:"isi"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
