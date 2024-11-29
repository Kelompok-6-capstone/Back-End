package repository

import (
	"calmind/model"

	"gorm.io/gorm"
)

type UserProfilRepository interface {
	GetByID(id int) (*model.User, error)
}

type userProfilRepositoryImpl struct {
	DB *gorm.DB
}

func NewUserProfilRepository(db *gorm.DB) UserProfilRepository {
	return &userProfilRepositoryImpl{DB: db}
}

func (r *userProfilRepositoryImpl) GetByID(id int) (*model.User, error) {
	var user model.User
	err := r.DB.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
