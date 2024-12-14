package model

import "time"

type Consultation struct {
	ID            int           `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID        int           `json:"user_id" gorm:"not null"`   // Foreign key ke tabel Users
	DoctorID      int           `json:"doctor_id" gorm:"not null"` // Foreign key ke tabel Doctors
	User          User          `json:"user" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
	Doctor        Doctor        `json:"doctor" gorm:"foreignKey:DoctorID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
	Title         string        `json:"title" gorm:"not null"`
	Description   string        `json:"description" gorm:"type:text;not null"`
	Duration      int           `json:"duration" gorm:"default:120"`                            // Durasi (default: 2 jam)
	IsApproved    bool          `json:"is_approved" gorm:"default:false"`                       // Status persetujuan admin
	PaymentStatus string        `json:"payment_status" gorm:"type:varchar(20)"`                 // pending, completed, failed
	OrderID       string        `json:"order_id"`                                               // ID pembayaran (Midtrans)
	Status        string        `json:"status" gorm:"type:varchar(20);default:'pending';index"` // Indeks Status
	StartTime     time.Time     `json:"start_time"`
	CreatedAt     time.Time     `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time     `json:"updated_at" gorm:"autoUpdateTime"`
	Rekomendasi   []Rekomendasi `json:"rekomendasi" gorm:"foreignKey:ConsultationID"`
}

type Rekomendasi struct {
	ID             int       `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID         int       `gorm:"not null"` // Foreign key ke User
	ConsultationID int       `json:"consultation_id" gorm:"not null;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
	DoctorID       int       `json:"doctor_id" gorm:"not null;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
	Rekomendasi    string    `json:"rekomendasi" gorm:"type:text;not null"`
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
}
