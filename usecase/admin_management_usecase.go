package usecase

import (
	"calmind/model"
	"calmind/repository"
	"time"
)

type AdminManagementUsecase interface {
	GetAllUsers() ([]UserDTO, error)
	GetAllDoctors() ([]DoctorDTO, error)
	DeleteUser(id int) (*model.User, error)
	DeleteDoctor(id int) (*model.Doctor, error)
}

type AdminManagementUsecaseImpl struct {
	Repo repository.AdminManagementRepo
}

func NewAdminManagementUsecase(repo repository.AdminManagementRepo) AdminManagementUsecase {
	return &AdminManagementUsecaseImpl{Repo: repo}
}

// DTO Struct untuk User
type UserDTO struct {
	ID               int                     `json:"id"`
	Username         string                  `json:"username"`
	Email            string                  `json:"email"`
	NoHp             string                  `json:"no_hp"`
	Alamat           string                  `json:"alamat"`
	TglLahir         string                  `json:"tgl_lahir"`
	JenisKelamin     string                  `json:"jenis_kelamin"`
	Pekerjaan        string                  `json:"pekerjaan"`
	IsVerified       bool                    `json:"is_verified"`
	CreatedAt        time.Time               `json:"created_at"`
	UpdatedAt        time.Time               `json:"updated_at"`
	LastConsultation *ConsultationSummaryDTO `json:"last_consultation"`
}

// DTO Struct untuk Doctor
type DoctorDTO struct {
	ID               int                     `json:"id"`
	Username         string                  `json:"username"`
	Email            string                  `json:"email"`
	Price            float64                 `json:"price"`
	Experience       int                     `json:"experience"`
	JenisKelamin     string                  `json:"jenis_kelamin"`
	IsVerified       bool                    `json:"is_verified"`
	IsActive         bool                    `json:"is_active"`
	Title            model.Title             `json:"title"`
	Tags             []model.Tags            `json:"tags"`
	CreatedAt        time.Time               `json:"created_at"`
	UpdatedAt        time.Time               `json:"updated_at"`
	LastConsultation *ConsultationSummaryDTO `json:"last_consultation"`
}

// DTO Struct untuk Konsultasi Terakhir
type ConsultationSummaryDTO struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	Description string    `json:"description"`
}

// Mendapatkan semua pengguna dengan konsultasi terakhir mereka
func (au *AdminManagementUsecaseImpl) GetAllUsers() ([]UserDTO, error) {
	users, err := au.Repo.FindAllUsersWithLastConsultation()
	if err != nil {
		return nil, err
	}

	var result []UserDTO
	for _, user := range users {
		var lastConsultation *ConsultationSummaryDTO
		if len(user.Consultations) > 0 {
			consultation := user.Consultations[0]
			lastConsultation = &ConsultationSummaryDTO{
				ID:          consultation.ID,
				Title:       consultation.Title,
				Status:      consultation.Status,
				CreatedAt:   consultation.CreatedAt,
				Description: consultation.Description,
			}
		}

		result = append(result, UserDTO{
			ID:               user.ID,
			Username:         user.Username,
			Email:            user.Email,
			NoHp:             user.NoHp,
			Alamat:           user.Alamat,
			TglLahir:         user.TglLahir,
			JenisKelamin:     user.JenisKelamin,
			Pekerjaan:        user.Pekerjaan,
			IsVerified:       user.IsVerified,
			CreatedAt:        user.CreatedAt,
			UpdatedAt:        user.UpdatedAt,
			LastConsultation: lastConsultation,
		})
	}

	return result, nil
}

// Mendapatkan semua dokter dengan konsultasi terakhir mereka
func (au *AdminManagementUsecaseImpl) GetAllDoctors() ([]DoctorDTO, error) {
	doctors, err := au.Repo.FindAllDoctorsWithLastConsultation()
	if err != nil {
		return nil, err
	}

	var result []DoctorDTO
	for _, doctor := range doctors {
		var lastConsultation *ConsultationSummaryDTO
		if len(doctor.Consultations) > 0 {
			consultation := doctor.Consultations[0]
			lastConsultation = &ConsultationSummaryDTO{
				ID:          consultation.ID,
				Title:       consultation.Title,
				Status:      consultation.Status,
				CreatedAt:   consultation.CreatedAt,
				Description: consultation.Description,
			}
		}

		result = append(result, DoctorDTO{
			ID:               doctor.ID,
			Username:         doctor.Username,
			Email:            doctor.Email,
			Price:            doctor.Price,
			Experience:       doctor.Experience,
			JenisKelamin:     doctor.JenisKelamin,
			IsVerified:       doctor.IsVerified,
			IsActive:         doctor.IsActive,
			Title:            doctor.Title,
			Tags:             doctor.Tags,
			CreatedAt:        doctor.CreatedAt,
			UpdatedAt:        doctor.UpdatedAt,
			LastConsultation: lastConsultation,
		})
	}

	return result, nil
}

// Menghapus pengguna berdasarkan ID
func (au *AdminManagementUsecaseImpl) DeleteUser(id int) (*model.User, error) {
	return au.Repo.DeleteUser(id)
}

// Menghapus dokter berdasarkan ID
func (au *AdminManagementUsecaseImpl) DeleteDoctor(id int) (*model.Doctor, error) {
	return au.Repo.DeleteDoctor(id)
}