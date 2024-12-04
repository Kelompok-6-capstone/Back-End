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

type ArtikelUsecaseImpl struct {
	Repo repository.ArtikelRepository
}

func NewArtikelUsecase(repo repository.ArtikelRepository) ArtikelUsecase {
	return &ArtikelUsecaseImpl{Repo: repo}
}

func (u *ArtikelUsecaseImpl) CreateArtikel(adminID int, artikel *model.Artikel) error {
	artikel.AdminID = adminID
	return u.Repo.Create(artikel)
}

func (u *ArtikelUsecaseImpl) GetAllArtikel() ([]model.Artikel, error) {
	return u.Repo.GetAll()
}

func (u *ArtikelUsecaseImpl) GetArtikelByID(id int) (*model.Artikel, error) {
	return u.Repo.GetByID(id)
}

func (u *ArtikelUsecaseImpl) UpdateArtikel(adminID int, artikel *model.Artikel) error {
	existingArtikel, err := u.Repo.GetByID(artikel.ID)
	if err != nil {
		return err
	}
	if existingArtikel.AdminID != adminID {
		return errors.New("unauthorized to update this artikel")
	}
	return u.Repo.Update(artikel)
}

func (u *ArtikelUsecaseImpl) DeleteArtikel(adminID int, id int) error {
	existingArtikel, err := u.Repo.GetByID(id)
	if err != nil {
		return err
	}
	if existingArtikel.AdminID != adminID {
		return errors.New("unauthorized to delete this artikel")
	}
	return u.Repo.Delete(id)
}
