package repository

import (
	"calmind/model"

	"gorm.io/gorm"
)

type DoctorProfilRepository interface {
	GetByID(id int) (*model.Doctor, error)
	UpdateByID(id int, doctor *model.Doctor) (*model.Doctor, error)
	UpdateDoctorActiveStatus(id int, isActive bool) error
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

	// Update fields
	if doctor.Username != "" {
		existingDoctor.Username = doctor.Username
	}
	if doctor.NoHp != "" {
		existingDoctor.NoHp = doctor.NoHp
	}
	if doctor.Avatar != "" {
		existingDoctor.Avatar = doctor.Avatar
	}
	if doctor.DateOfBirth != "" {
		existingDoctor.DateOfBirth = doctor.DateOfBirth
	}
	if doctor.Address != "" {
		existingDoctor.Address = doctor.Address
	}
	if doctor.Schedule != "" {
		existingDoctor.Schedule = doctor.Schedule
	}
	if doctor.Title != "" {
		existingDoctor.Title = doctor.Title
	}
	if doctor.Price > 0 {
		existingDoctor.Price = doctor.Price
	}
	if doctor.Experience > 0 {
		existingDoctor.Experience = doctor.Experience
	}
	if doctor.STRNumber != "" {
		existingDoctor.STRNumber = doctor.STRNumber
	}
	if doctor.About != "" {
		existingDoctor.About = doctor.About
	}

	// Save updated doctor
	err = r.DB.Save(&existingDoctor).Error
	if err != nil {
		return nil, err
	}

	return &existingDoctor, nil
}

func (r *DoctorProfilRepositoryImpl) UpdateDoctorActiveStatus(id int, isActive bool) error {
	return r.DB.Model(&model.Doctor{}).Where("id = ?", id).Update("is_active", isActive).Error
}
