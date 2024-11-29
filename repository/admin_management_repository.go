package repository

import (
	"calmind/model"

	"gorm.io/gorm"
)

type AdminManagementRepo interface {
	FindAllUsers() ([]*model.User, error)
	DeleteUsers(id int) (*model.User, error)
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
