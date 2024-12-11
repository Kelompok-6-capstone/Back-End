package repository

import (
	"calmind/model"

	"gorm.io/gorm"
)

type AdminManagementRepo interface {
	FindAllUsers() ([]*model.User, error)
	FindAllDoctors() ([]*model.Doctor, error)
	DeleteUser(id int) (*model.User, error)
	DeleteDoctor(id int) (*model.Doctor, error)
	FindAllUsersWithLastConsultation() ([]*model.User, error)
	FindAllDoctorsWithLastConsultation() ([]*model.Doctor, error)
}

type AdminManagementRepoImpl struct {
	DB *gorm.DB
}

func NewAdminManagementRepo(db *gorm.DB) AdminManagementRepo {
	return &AdminManagementRepoImpl{DB: db}
}

// Find all users
func (ar *AdminManagementRepoImpl) FindAllUsers() ([]*model.User, error) {
	var users []*model.User
	err := ar.DB.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

// Find all doctors
func (ar *AdminManagementRepoImpl) FindAllDoctors() ([]*model.Doctor, error) {
	var doctors []*model.Doctor
	err := ar.DB.Find(&doctors).Error
	if err != nil {
		return nil, err
	}
	return doctors, nil
}

// Find all users with their last consultation
func (ar *AdminManagementRepoImpl) FindAllUsersWithLastConsultation() ([]*model.User, error) {
	var users []*model.User
	err := ar.DB.Preload("Consultations", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at DESC").Limit(1)
	}).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

// Find all doctors with their last consultation
func (ar *AdminManagementRepoImpl) FindAllDoctorsWithLastConsultation() ([]*model.Doctor, error) {
	var doctors []*model.Doctor
	err := ar.DB.Preload("Consultations", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at DESC").Limit(1)
	}).Preload("Title").Preload("Tags").Find(&doctors).Error
	if err != nil {
		return nil, err
	}
	return doctors, nil
}

// Delete user by ID
func (ar *AdminManagementRepoImpl) DeleteUser(id int) (*model.User, error) {
	var user model.User
	err := ar.DB.First(&user, id).Error
	if err != nil {
		return nil, err
	}

	err = ar.DB.Delete(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Delete doctor by ID
func (ar *AdminManagementRepoImpl) DeleteDoctor(id int) (*model.Doctor, error) {
	var doctor model.Doctor
	err := ar.DB.First(&doctor, id).Error
	if err != nil {
		return nil, err
	}

	err = ar.DB.Delete(&doctor).Error
	if err != nil {
		return nil, err
	}

	return &doctor, nil
}
