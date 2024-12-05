package usecase

import (
	"calmind/model"
	"calmind/repository"
)

type UserFiturUsecase interface {
	GetAllDoctors() ([]model.Doctor, error)
	GetDoctorsBySpecialty(specialty string) ([]model.Doctor, error)
	GetDoctorsByStatus(isActive bool) ([]model.Doctor, error)
	SearchDoctors(query string) ([]model.Doctor, error)
	GetDoctorByID(id int) (*model.Doctor, error)
	GetAllSpesialis() ([]model.Specialty, error)
}

type UserFiturUsecaseImpl struct {
	DoctorRepo repository.UserFiturRepository
}

func NewUserFiturUsecase(repo repository.UserFiturRepository) UserFiturUsecase {
	return &UserFiturUsecaseImpl{DoctorRepo: repo}
}

func (u *UserFiturUsecaseImpl) GetAllDoctors() ([]model.Doctor, error) {
	return u.DoctorRepo.GetAllDoctors()
}

func (u *UserFiturUsecaseImpl) GetDoctorsBySpecialty(specialty string) ([]model.Doctor, error) {
	return u.DoctorRepo.GetDoctorsBySpecialty(specialty)
}

func (u *UserFiturUsecaseImpl) GetDoctorsByStatus(isActive bool) ([]model.Doctor, error) {
	return u.DoctorRepo.GetDoctorsByStatus(isActive)
}

func (u *UserFiturUsecaseImpl) SearchDoctors(query string) ([]model.Doctor, error) {
	return u.DoctorRepo.SearchDoctors(query)
}

func (u *UserFiturUsecaseImpl) GetDoctorByID(id int) (*model.Doctor, error) {
	return u.DoctorRepo.GetDoctorByID(id)
}
func (u *UserFiturUsecaseImpl) GetAllSpesialis() ([]model.Specialty, error) {
	return u.DoctorRepo.GetSpesialis()
}
