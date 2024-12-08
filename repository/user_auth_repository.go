package repository

import (
	"calmind/model"
	"errors"

	"gorm.io/gorm"
)

type UserRepository interface {
	GetByUsername(email string) (*model.User, error)
	CreateUser(*model.User) error
	UpdateUserVerificationStatus(email string, isVerified bool) error
}

type userRepositorystate struct {
	DB *gorm.DB
}

func NewAuthRepository(db *gorm.DB) UserRepository {
	return &userRepositorystate{DB: db}
}

func (r *userRepositorystate) CreateUser(user *model.User) error {
	var existingUser model.User
	err := r.DB.Where("email = ?", user.Email).First(&existingUser).Error
	if err == nil {
		return errors.New("email already registered")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	return r.DB.Create(user).Error
}

func (r *userRepositorystate) GetByUsername(email string) (*model.User, error) {
	var user model.User
	err := r.DB.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *userRepositorystate) UpdateUserVerificationStatus(email string, isVerified bool) error {
	return r.DB.Model(&model.User{}).Where("email = ?", email).Update("is_verified", isVerified).Error
}
