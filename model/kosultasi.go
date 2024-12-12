package model

import "time"

type Consultation struct {
	ID            int           `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID        int           `json:"user_id" gorm:"not null"`
	DoctorID      int           `json:"doctor_id" gorm:"not null"`
	User          User          `json:"user" gorm:"foreignKey:UserID"`     // Relasi dengan user
	Doctor        Doctor        `json:"doctor" gorm:"foreignKey:DoctorID"` // Relasi dengan dokter
	Title         string        `json:"title" gorm:"not null"`
	Description   string        `json:"description" gorm:"type:text;not null"`
	Duration      int           `json:"duration" gorm:"default:120"`                            // Durasi dalam menit (default: 2 jam)
	IsApproved    bool          `json:"is_approved" gorm:"default:false"`                       // Status persetujuan admin
	PaymentStatus string        `json:"payment_status"`                                         // pending, completed, failed
	OrderID       string        `json:"order_id"`                                               // Order ID dari Midtrans
	Status        string        `json:"status" gorm:"type:varchar(20);default:'pending';index"` // Indeks pada Status
	StartTime     time.Time     `json:"start_time"`
	CreatedAt     time.Time     `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time     `json:"updated_at" gorm:"autoUpdateTime"`
	Rekomendasi   []Rekomendasi `json:"rekomendasi" gorm:"foreignKey:ConsultationID"` // Rekomendasi dari dokter
}

type Rekomendasi struct {
	ID             int    `json:"id" gorm:"primaryKey;autoIncrement"`
	ConsultationID int    `json:"consultation_id" gorm:"not null"`
	DoctorID       int    `json:"doctor_id" gorm:"not null"`
	Rekomendasi    string `json:"rekomendasi" gorm:"type:text;not null"`
}
