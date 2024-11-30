package repository

import (
	"calmind/model"

	"gorm.io/gorm"
)

type DoctorProfilRepository interface {
	GetByID(id int) (*model.Doctor, error)
	UpdateByID(id int, doctor *model.Doctor) (*model.Doctor, error)
}

type DoctorProfilRepositoryImpl struct {
	DB *gorm.DB
}

func NewDoctorProfilRepository(db *gorm.DB) DoctorProfilRepository {
	return &DoctorProfilRepositoryImpl{DB: db}
}

func (r *DoctorProfilRepositoryImpl) GetByID(id int) (*model.Doctor, error) {
	var doctor model.Doctor
	err := r.DB.Where("id = ?", id).First(&doctor).Error
	if err != nil {
		return nil, err
	}
	return &doctor, nil
}

func (r *DoctorProfilRepositoryImpl) UpdateByID(id int, doctor *model.Doctor) (*model.Doctor, error) {
	var existingDoctor model.Doctor
	err := r.DB.Where("id = ?", id).First(&existingDoctor).Error
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if doctor.Username != "" {
		existingDoctor.Username = doctor.Username
	}
	if doctor.NoHp != "" {
		existingDoctor.NoHp = doctor.NoHp
	}

	// Save updates to the database
	err = r.DB.Save(&existingDoctor).Error
	if err != nil {
		return nil, err
	}

	return &existingDoctor, nil
}
