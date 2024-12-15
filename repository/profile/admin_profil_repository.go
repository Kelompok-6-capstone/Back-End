package repository

import (
	"calmind/model"

	"gorm.io/gorm"
)

type AdminProfileRepository interface {
	GetByID(id int) (*model.Admin, error)
	UpdateAvatarByID(id int, avatar string, deleteURL string) error
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

func (r *AdminProfileRepositoryImpl) UpdateAvatarByID(id int, avatar string, deleteURL string) error {
	err := r.DB.Model(&model.Admin{}).Where("id = ?", id).Updates(map[string]interface{}{
		"avatar":     avatar,
		"delete_url": deleteURL,
	}).Error
	return err
}

func (r *AdminProfileRepositoryImpl) ClearAvatarByID(id int) error {
	err := r.DB.Model(&model.Admin{}).Where("id = ?", id).Updates(map[string]interface{}{
		"avatar":     "",
		"delete_url": "",
	}).Error
	return err
}
