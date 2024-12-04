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

type artikelRepository struct {
	db *gorm.DB
}

func NewArtikelRepository(db *gorm.DB) ArtikelRepository {
	return &artikelRepository{db: db}
}

func (r *artikelRepository) Create(artikel *model.Artikel) error {
	return r.db.Create(artikel).Error
}

func (r *artikelRepository) GetAll() ([]model.Artikel, error) {
	var artikels []model.Artikel
	err := r.db.Preload("Admin").Find(&artikels).Error
	return artikels, err
}

func (r *artikelRepository) GetByID(id int) (*model.Artikel, error) {
	var artikel model.Artikel
	err := r.db.Preload("Admin").First(&artikel, id).Error
	return &artikel, err
}

func (r *artikelRepository) Update(artikel *model.Artikel) error {
	return r.db.Save(artikel).Error
}

func (r *artikelRepository) Delete(id int) error {
	return r.db.Delete(&model.Artikel{}, id).Error
}
