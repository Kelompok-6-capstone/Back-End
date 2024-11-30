package repository

import (
	"calmind/model"

	"gorm.io/gorm"
)

type UserProfilRepository interface {
	GetByID(id int) (*model.User, error)
	UpdateByID(id int, user *model.User) (*model.User, error)
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

func (r *userProfilRepositoryImpl) UpdateByID(id int, user *model.User) (*model.User, error) {
	var existingUser model.User
	err := r.DB.Where("id = ?", id).First(&existingUser).Error
	if err != nil {
		return nil, err
	}

	existingUser.Avatar = user.Avatar
	existingUser.Username = user.Username
	existingUser.Email = user.Email
	existingUser.NoHp = user.NoHp
	existingUser.Alamat = user.Alamat
	existingUser.Tgl_lahir = user.Tgl_lahir
	existingUser.JenisKelamin = user.JenisKelamin

	err = r.DB.Save(&existingUser).Error
	if err != nil {
		return nil, err
	}

	return &existingUser, nil
}
