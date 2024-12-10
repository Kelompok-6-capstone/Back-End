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
	SearchArtikel(query string) ([]model.Artikel, error)
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
	return r.db.Model(&model.Artikel{}).Where("id = ?", artikel.ID).
		Updates(map[string]interface{}{
			"judul":      artikel.Judul,
			"gambar":     artikel.Gambar,
			"isi":        artikel.Isi,
			"updated_at": artikel.UpdatedAt,
		}).Error
}

func (r *artikelRepository) Delete(id int) error {
	return r.db.Delete(&model.Artikel{}, id).Error
}

func (r *artikelRepository) SearchArtikel(query string) ([]model.Artikel, error) {
	var artikel []model.Artikel
	err := r.db.
		Where("judul LIKE ?", "%"+query+"%").
		Find(&artikel).Error
	return artikel, err
}
