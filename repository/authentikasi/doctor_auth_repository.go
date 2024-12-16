package repository

import (
	"calmind/model"
	"errors"

	"gorm.io/gorm"
)

type DoctorRepository interface {
	GetByEmail(email string) (*model.Doctor, error)
	CreateDoctor(*model.Doctor) error
	UpdateDokterVerificationStatus(email string, isVerified bool) error
}

type doctorRepositoryState struct {
	DB *gorm.DB
}

func NewDoctorAuthRepository(db *gorm.DB) DoctorRepository {
	return &doctorRepositoryState{DB: db}
}

func (r *doctorRepositoryState) CreateDoctor(doctor *model.Doctor) error {
	// Tetapkan default title_id jika tidak diberikan
	if doctor.TitleID == 0 {
		doctor.TitleID = 1 // ID default dari title
	}

	var existingDoctor model.Doctor
	err := r.DB.Where("email = ?", doctor.Email).First(&existingDoctor).Error
	if err == nil {
		return errors.New("Email sudah terdaftar, gunakan email lain.")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("Terjadi kesalahan saat memeriksa email.")
	}

	return r.DB.Create(doctor).Error
}

func (r *doctorRepositoryState) GetByEmail(email string) (*model.Doctor, error) {
	var doctor model.Doctor
	err := r.DB.Where("email = ?", email).First(&doctor).Error
	if err != nil {
		return nil, err
	}
	return &doctor, nil
}

func (r *doctorRepositoryState) UpdateDokterVerificationStatus(email string, isVerified bool) error {
	return r.DB.Model(&model.Doctor{}).Where("email = ?", email).Update("is_verified", isVerified).Error
}
