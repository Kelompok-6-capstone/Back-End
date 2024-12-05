package repository

import (
	"calmind/model"

	"gorm.io/gorm"
)

type UserProfilRepository interface {
	GetByID(id int) (*model.User, error)
	UpdateByID(id int, User *model.User) (*model.User, error)
}

type UserProfilRepositoryImpl struct {
	DB *gorm.DB
}

func NewUserProfilRepository(db *gorm.DB) UserProfilRepository {
	return &UserProfilRepositoryImpl{DB: db}
}

func (r *UserProfilRepositoryImpl) GetByID(id int) (*model.User, error) {
	var user model.User
	err := r.DB.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserProfilRepositoryImpl) UpdateByID(id int, user *model.User) (*model.User, error) {
	var existingUser model.User
	err := r.DB.Where("id = ?", id).First(&existingUser).Error
	if err != nil {
		return nil, err
	}

	if user.Avatar != "" {
		existingUser.Avatar = user.Avatar
	}
	if user.Username != "" {
		existingUser.Username = user.Username
	}
	if user.NoHp != "" {
		existingUser.NoHp = user.NoHp
	}
	if user.Alamat != "" {
		existingUser.Alamat = user.Alamat
	}
	if user.Tgl_lahir != "" {
		existingUser.Tgl_lahir = user.Tgl_lahir
	}
	if user.JenisKelamin != "" {
		existingUser.JenisKelamin = user.JenisKelamin
	}
	if user.Pekerjaan != "" {
		existingUser.Pekerjaan = user.Pekerjaan
	}


	// Simpan perubahan
	err = r.DB.Save(&existingUser).Error
	if err != nil {
		return nil, err
	}

	return &existingUser, nil
}
