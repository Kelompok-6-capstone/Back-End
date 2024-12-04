package usecase

import (
	"calmind/model"
	"calmind/repository"
	"errors"
)

type ArtikelUsecase interface {
	CreateArtikel(adminID int, artikel *model.Artikel) error
	GetAllArtikel() ([]model.Artikel, error)
	GetArtikelByID(id int) (*model.Artikel, error)
	UpdateArtikel(adminID int, artikel *model.Artikel) error
	DeleteArtikel(adminID int, id int) error
}

type artikelUsecase struct {
	repo repository.ArtikelRepository
}

func NewArtikelUsecase(repo repository.ArtikelRepository) ArtikelUsecase {
	return &artikelUsecase{repo: repo}
}

func (u *artikelUsecase) CreateArtikel(adminID int, artikel *model.Artikel) error {
	artikel.AdminID = adminID
	return u.repo.Create(artikel)
}

func (u *artikelUsecase) GetAllArtikel() ([]model.Artikel, error) {
	return u.repo.GetAll()
}

func (u *artikelUsecase) GetArtikelByID(id int) (*model.Artikel, error) {
	return u.repo.GetByID(id)
}

func (u *artikelUsecase) UpdateArtikel(adminID int, artikel *model.Artikel) error {
	existingArtikel, err := u.repo.GetByID(artikel.ID)
	if err != nil {
		return err
	}
	if existingArtikel.AdminID != adminID {
		return errors.New("unauthorized to update this artikel")
	}
	return u.repo.Update(artikel)
}

func (u *artikelUsecase) DeleteArtikel(adminID int, id int) error {
	existingArtikel, err := u.repo.GetByID(id)
	if err != nil {
		return err
	}
	if existingArtikel.AdminID != adminID {
		return errors.New("unauthorized to delete this artikel")
	}
	return u.repo.Delete(id)
}
