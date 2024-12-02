package repository

import (
	"calmind/model"

	"gorm.io/gorm"
)

type UserFiturRepository interface {
	GetAllDoctors() ([]model.Doctor, error)
	GetDoctorsBySpecialty(specialty string) ([]model.Doctor, error)
	GetDoctorsByStatus(isActive bool) ([]model.Doctor, error)
	SearchDoctors(query string) ([]model.Doctor, error)
	GetDoctorByID(id int) (*model.Doctor, error)
}

type UserFiturRepositoryImpl struct {
	DB *gorm.DB
}

func NewUserFiturRepository(db *gorm.DB) UserFiturRepository {
	return &UserFiturRepositoryImpl{DB: db}
}

func (r *UserFiturRepositoryImpl) GetAllDoctors() ([]model.Doctor, error) {
	var doctors []model.Doctor
	err := r.DB.Preload("Specialties").
		Where("title != '' AND price > 0 AND experience > 0 AND is_verified = true AND is_active = true").
		Find(&doctors).Error
	return doctors, err
}

func (r *UserFiturRepositoryImpl) GetDoctorsBySpecialty(specialty string) ([]model.Doctor, error) {
	var doctors []model.Doctor
	err := r.DB.Preload("Specialties").
		Joins("JOIN doctor_specialties ON doctors.id = doctor_specialties.doctor_id").
		Joins("JOIN specialties ON doctor_specialties.specialty_id = specialties.id").
		Where("specialties.name = ?", specialty).
		Find(&doctors).Error
	return doctors, err
}

func (r *UserFiturRepositoryImpl) GetDoctorsByStatus(isActive bool) ([]model.Doctor, error) {
	var doctors []model.Doctor
	err := r.DB.Preload("Specialties").Where("is_active = ?", isActive).Find(&doctors).Error
	return doctors, err
}

func (r *UserFiturRepositoryImpl) SearchDoctors(query string) ([]model.Doctor, error) {
	var doctors []model.Doctor
	err := r.DB.Preload("Specialties").
		Where("username LIKE ? OR title LIKE ? OR about LIKE ?", "%"+query+"%", "%"+query+"%", "%"+query+"%").
		Find(&doctors).Error
	return doctors, err
}

func (r *UserFiturRepositoryImpl) GetDoctorByID(id int) (*model.Doctor, error) {
	var doctor model.Doctor
	err := r.DB.Preload("Specialties").Where("id = ?", id).First(&doctor).Error
	return &doctor, err
}
