package model

// DTO untuk User
type UserDTO struct {
	Avatar    string `gorm:"" json:"avatar"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	TglLahir  string `gorm:"" json:"tgl_lahir"`
	Pekerjaan string `gorm:"" json:"pekerjaan"`
}

type DoctorDTO struct {
	Username string  `json:"username"`
	Email    string  `json:"email"`
	Avatar   string  `json:"avatar,omitempty"`
	Price    float64 `json:"price,omitempty"`
	About    string  `json:"about"`
}

type RecommendationDTO struct {
	ID             int    `json:"id"`
	ConsultationID int    `json:"consultation_id"`
	DoctorID       int    `json:"doctor_id"`
	Recommendation string `json:"recommendation"`
}

type ConsultationDTO struct {
	ID          int                 `json:"id"`
	Title       string              `json:"title"`
	Description string              `json:"description"`
	Duration    int                 `json:"duration"`
	Status      string              `json:"status"`
	StartTime   string              `json:"start_time,omitempty"`
	OrderID     string              `json:"order_id,omitempty"`
	User        *UserDTO            `json:"user,omitempty"`
	Doctor      *DoctorDTO          `json:"doctor,omitempty"`
	PaymentURL  string              `json:"payment_url,omitempty"`
	Rekomendasi []RecommendationDTO `json:"rekomendasi,omitempty"`
}

type SimpleConsultationDTO struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Duration    int    `json:"duration"`
	Status      string `json:"status"`
	StartTime   string `json:"start_time"`
	OrderID     string `json:"order_id"`
	PaymentURL  string `json:"payment_url"`
}
