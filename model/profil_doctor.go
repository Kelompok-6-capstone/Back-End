package model

import "time"

type DoctorProfile struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Name         string    `json:"name"`          // Nama lengkap dokter
	Title        string    `json:"title"`         // Gelar dan spesialisasi
	Price        float64   `json:"price"`         // Biaya konsultasi
	Experience   int       `json:"experience"`    // Lama pengalaman dalam tahun
	SuccessRate  float32   `json:"success_rate"`  // Tingkat keberhasilan (dalam %)
	Specialty    string    `json:"specialty"`     // Spesialisasi utama
	SubSpecialty string    `json:"sub_specialty"` // Sub-spesialisasi
	Schedule     string    `json:"schedule"`      // Jadwal kerja
	STRNumber    string    `json:"str_number"`    // Nomor STR
	About        string    `json:"about"`         // Tentang dokter
	Photo        string    `json:"photo"`         // URL foto dokter
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
