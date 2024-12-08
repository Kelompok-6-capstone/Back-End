package model

import "time"

type Consultation struct {
	ID          int       `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID      int       `json:"user_id" gorm:"not null"`
	DoctorID    int       `json:"doctor_id" gorm:"not null"`
	Message     string    `json:"message" gorm:"type:text;not null"`
	Rekomendasi string    `json:"rekomendasi" gorm:"type:text"`                     // Rekomendasi dari dokter
	IsPaid      bool      `json:"is_paid" gorm:"default:false"`                     // Status pembayaran
	IsApproved  bool      `json:"is_approved" gorm:"default:false"`                 // Status persetujuan admin
	Status      string    `json:"status" gorm:"type:varchar(20);default:'pending'"` // Status konsultasi (e.g., pending, active, expired)
	StartTime   time.Time `json:"start_time" gorm:"not null"`                       // Waktu mulai konsultasi
	Duration    int       `json:"duration" gorm:"default:120"`                      // Durasi dalam menit (default: 2 jam)
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`                 // Waktu pembuatan konsultasi
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`                 // Waktu pembaruan konsultasi
	User        User      `json:"user" gorm:"foreignKey:UserID"`                    // Relasi dengan user
	Doctor      Doctor    `json:"doctor" gorm:"foreignKey:DoctorID"`                // Relasi dengan dokter
}
