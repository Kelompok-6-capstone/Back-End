package repository

import (
	"calmind/model"
	"errors"

	"gorm.io/gorm"
)

type DoctorRepository interface {
	GetByEmail(email string) (*model.Doctor, error)
	CreateDoctor(*model.Doctor) error
}

type doctorRepositoryState struct {
	DB *gorm.DB
}

func NewDoctorAuthRepository(db *gorm.DB) DoctorRepository {
	return &doctorRepositoryState{DB: db}
}

func (r *doctorRepositoryState) CreateDoctor(doctor *model.Doctor) error {
	var existingDoctor model.Doctor
	err := r.DB.Where("email = ?", doctor.Email).First(&existingDoctor).Error
	if err == nil {
		return errors.New("email already registered")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
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

