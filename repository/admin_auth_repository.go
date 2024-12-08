package repository

import (
	"calmind/model"

	"gorm.io/gorm"
)

type AdminAuthRepository interface {
	GetByEmail(email string) (*model.Admin, error)
}

type AdminAuthRepositoryImpl struct {
	Database *gorm.DB
}

func NewAdminAuthRepository(db *gorm.DB) AdminAuthRepository {
	return &AdminAuthRepositoryImpl{Database: db}
}

func (ar *AdminAuthRepositoryImpl) GetByEmail(email string) (*model.Admin, error) {
	var admin model.Admin
	err := ar.Database.Where("email = ?", email).First(&admin).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}
