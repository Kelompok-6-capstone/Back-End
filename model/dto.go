package model

// DTO untuk User
type UserDTO struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar,omitempty"`
}

// DTO untuk Dokter
type DoctorDTO struct {
	Username string  `json:"username"`
	Email    string  `json:"email"`
	Avatar   string  `json:"avatar,omitempty"`
	Price    float64 `json:"price,omitempty"`
}

// DTO untuk Rekomendasi
type RecommendationDTO struct {
	ID             int    `json:"id"`
	ConsultationID int    `json:"consultation_id"`
	DoctorID       int    `json:"doctor_id"`
	Recommendation string `json:"recommendation"`
}

// DTO untuk Konsultasi
type ConsultationDTO struct {
	ID          int                 `json:"id"`
	Title       string              `json:"title"`
	Description string              `json:"description"`
	Duration    int                 `json:"duration"`
	Status      string              `json:"status"`
	StartTime   string              `json:"start_time,omitempty"`
	User        *UserDTO            `json:"user,omitempty"`
	Doctor      *DoctorDTO          `json:"doctor,omitempty"`
	PaymentURL  string              `json:"payment_url,omitempty"`
	Rekomendasi []RecommendationDTO `json:"rekomendasi,omitempty"`
}
