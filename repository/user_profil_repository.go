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

	if user.Username != "" {
		existingUser.Username = user.Username
	}
	if user.NoHp != "" {
		existingUser.NoHp = user.NoHp
	}
	if user.Avatar != "" {
		existingUser.Avatar = user.Avatar
	}
	if user.Bio != "" {
		existingUser.Bio = user.Bio
	}

	err = r.DB.Save(&existingUser).Error
	if err != nil {
		return nil, err
	}

	return &existingUser, nil
}
