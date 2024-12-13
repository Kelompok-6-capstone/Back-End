package model

// DTO untuk User
type UserDTO struct {
	ID        int    `json:"id"` // Tambahkan ID user
	Avatar    string `json:"avatar"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	TglLahir  string `json:"tgl_lahir"`
	Pekerjaan string `json:"pekerjaan"`
}

type DoctorDTO struct {
	ID       int    `json:"id"` // Tambahkan ID dokter
	Avatar   string `json:"avatar"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Price    int    `json:"price"`
	About    string `json:"about"`
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
	StartTime   string              `json:"start_time"`
	OrderID     string              `json:"order_id"`
	UserID      int                 `json:"user_id"`   // Tambahkan UserID
	DoctorID    int                 `json:"doctor_id"` // Tambahkan DoctorID
	User        *UserDTO            `json:"user"`
	Doctor      *DoctorDTO          `json:"doctor"`
	Rekomendasi []RecommendationDTO `json:"rekomendasi"`
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
