package model

import "time"

type Consultation struct {
	ID             int       `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID         int       `json:"user_id" gorm:"not null"`
	DoctorID       int       `json:"doctor_id" gorm:"not null"`
	StartTime      time.Time `json:"start_time"`
	EndTime        time.Time `json:"end_time"`
	Message        string    `json:"message" gorm:"type:text"`
	Recommendation string    `json:"recommendation" gorm:"type:text"`
	Status         string    `json:"status" gorm:"default:'pending'"`
	PaymentStatus  string    `json:"payment_status" gorm:"default:'unpaid'"`
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	User           User      `json:"user" gorm:"foreignKey:UserID"`
	Doctor         Doctor    `json:"doctor" gorm:"foreignKey:DoctorID"`
}
