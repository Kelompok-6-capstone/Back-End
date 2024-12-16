package repository

import (
	"calmind/model"
	"fmt"

	"gorm.io/gorm"
)

type AdminProfileRepository interface {
	GetByID(id int) (*model.Admin, error)
	UpdateAvatarByID(id int, avatarURL string) error
	ClearAvatarByID(id int) error
}

type AdminProfileRepositoryImpl struct {
	DB *gorm.DB
}

func NewAdminProfileRepository(db *gorm.DB) AdminProfileRepository {
	return &AdminProfileRepositoryImpl{DB: db}
}

func (r *AdminProfileRepositoryImpl) GetByID(id int) (*model.Admin, error) {
	var admin model.Admin
	err := r.DB.Where("id = ?", id).First(&admin).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func (r *AdminProfileRepositoryImpl) UpdateAvatarByID(adminID int, avatarURL string) error {
	err := r.DB.Model(&model.Admin{}).Where("id = ?", adminID).Update("avatar", avatarURL).Error
	if err != nil {
		return fmt.Errorf("failed to update avatar: %v", err)
	}
	return nil
}

func (r *AdminProfileRepositoryImpl) ClearAvatarByID(id int) error {
	err := r.DB.Model(&model.Admin{}).Where("id = ?", id).Update("avatar", nil).Error
	return err
}
