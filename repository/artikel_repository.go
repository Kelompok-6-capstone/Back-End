package repository

import (
	"calmind/model"

	"gorm.io/gorm"
)

type ArtikelRepository interface {
	Create(artikel *model.Artikel) error
	GetAll() ([]model.Artikel, error)
	GetByID(id int) (*model.Artikel, error)
	Update(artikel *model.Artikel) error
	Delete(id int) error
}

type ArtikelRepositoryImpl struct {
	Database *gorm.DB
}

func NewArtikelRepository(db *gorm.DB) ArtikelRepository {
	return &ArtikelRepositoryImpl{Database: db}
}

func (r *ArtikelRepositoryImpl) Create(artikel *model.Artikel) error {
	return r.Database.Create(artikel).Error
}

func (r *ArtikelRepositoryImpl) GetAll() ([]model.Artikel, error) {
	var artikels []model.Artikel
	err := r.Database.Preload("Admin").Find(&artikels).Error
	return artikels, err
}

func (r *ArtikelRepositoryImpl) GetByID(id int) (*model.Artikel, error) {
	var artikel model.Artikel
	err := r.Database.Preload("Admin").First(&artikel, id).Error
	if err != nil {
		return nil, err
	}
	return &artikel, nil
}

func (r *ArtikelRepositoryImpl) Update(artikel *model.Artikel) error {
	return r.Database.Save(artikel).Error
}

func (r *ArtikelRepositoryImpl) Delete(id int) error {
	return r.Database.Delete(&model.Artikel{}, id).Error
}
