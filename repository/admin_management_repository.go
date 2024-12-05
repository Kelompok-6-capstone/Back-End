package repository

import (
	"calmind/model"

	"gorm.io/gorm"
)

type AdminManagementRepo interface {
	FindAllUsers() ([]*model.User, error)
	FindAllDocters() ([]*model.Doctor, error)
	DeleteUsers(id int) (*model.User, error)
	DeleteDocter(id int) (*model.Doctor, error)
	FindUserDetail(id int) (*model.User, error)
	FindDocterDetail(id int) (*model.Doctor, error)
}

type AdminManagementRepoImpl struct {
	DB *gorm.DB
}

func NewAdminManagementRepo(db *gorm.DB) AdminManagementRepo {
	return &AdminManagementRepoImpl{DB: db}
}

func (ar *AdminManagementRepoImpl) FindAllUsers() ([]*model.User, error) {
	var users []*model.User
	err := ar.DB.Find(&users).Error
	return users, err
}

func (ar *AdminManagementRepoImpl) DeleteUsers(id int) (*model.User, error) {
	var user model.User
	err := ar.DB.First(&user, id).Error
	if err != nil {
	}

	err = ar.DB.Delete(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (ar *AdminManagementRepoImpl) FindUserDetail(id int) (*model.User, error) {
	var user model.User
	err := ar.DB.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (ar *AdminManagementRepoImpl) FindAllDocters() ([]*model.Doctor, error) {
	var dokter []*model.Doctor
	err := ar.DB.Find(&dokter).Error
	return dokter, err
}

func (ar *AdminManagementRepoImpl) DeleteDocter(id int) (*model.Doctor, error) {
	var dokter model.Doctor
	err := ar.DB.First(&dokter, id).Error
	if err != nil {
	}

	err = ar.DB.Delete(&dokter).Error
	if err != nil {
		return nil, err
	}

	return &dokter, nil
}

func (ar *AdminManagementRepoImpl) FindDocterDetail(id int) (*model.Doctor, error) {
	var doctor model.Doctor
	err := ar.DB.First(&doctor, id).Error
	if err != nil {
		return nil, err
	}
	return &doctor, nil
}
