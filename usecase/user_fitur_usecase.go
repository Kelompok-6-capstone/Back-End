package usecase

import (
	"calmind/model"
	"calmind/repository"
)

type UserFiturUsecase interface {
	GetAllDoctors() ([]model.Doctor, error)
	GetDoctorsByTag(tag string) ([]model.Doctor, error)
	GetDoctorsByStatus(isActive bool) ([]model.Doctor, error)
	SearchDoctors(query string) ([]model.Doctor, error)
	GetDoctorByID(id int) (*model.Doctor, error)
	GetAllTags() ([]model.Tags, error)
	GetAllTitles() ([]model.Title, error)
	GetDoctorsByTitle(title string) ([]model.Doctor, error)
}

type UserFiturUsecaseImpl struct {
	DoctorRepo repository.UserFiturRepository
}

func NewUserFiturUsecase(repo repository.UserFiturRepository) UserFiturUsecase {
	return &UserFiturUsecaseImpl{DoctorRepo: repo}
}

// Mendapatkan semua dokter
func (u *UserFiturUsecaseImpl) GetAllDoctors() ([]model.Doctor, error) {
	doctors, err := u.DoctorRepo.GetAllDoctors()
	if err != nil {
		return nil, err
	}
	return doctors, nil
}

// Mendapatkan dokter berdasarkan tag
func (u *UserFiturUsecaseImpl) GetDoctorsByTag(tag string) ([]model.Doctor, error) {
	return u.DoctorRepo.GetDoctorsByTag(tag)
}

// Mendapatkan dokter berdasarkan status (aktif/tidak aktif)
func (u *UserFiturUsecaseImpl) GetDoctorsByStatus(isActive bool) ([]model.Doctor, error) {
	return u.DoctorRepo.GetDoctorsByStatus(isActive)
}

// Mencari dokter berdasarkan query
func (u *UserFiturUsecaseImpl) SearchDoctors(query string) ([]model.Doctor, error) {
	return u.DoctorRepo.SearchDoctors(query)
}

// Mendapatkan dokter berdasarkan ID
func (u *UserFiturUsecaseImpl) GetDoctorByID(id int) (*model.Doctor, error) {
	return u.DoctorRepo.GetDoctorByID(id)
}

// Mendapatkan semua tags
func (u *UserFiturUsecaseImpl) GetAllTags() ([]model.Tags, error) {
	return u.DoctorRepo.GetTags()
}

// Mendapatkan semua titles
func (u *UserFiturUsecaseImpl) GetAllTitles() ([]model.Title, error) {
	return u.DoctorRepo.GetTitles()
}

// Mendapatkan dokter berdasarkan title
func (u *UserFiturUsecaseImpl) GetDoctorsByTitle(title string) ([]model.Doctor, error) {
	return u.DoctorRepo.GetDoctorsByTitle(title)
}
