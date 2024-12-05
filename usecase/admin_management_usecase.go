package usecase

import (
	"calmind/model"
	"calmind/repository"
)

type AdminManagementUsecase interface {
	GetAllUsers() ([]*model.User, error)
	DeleteUsers(id int) (*model.User, error)
	GetAllDocter() ([]*model.Doctor, error)
	DeleteDocter(id int) (*model.Doctor, error)
	GetUserDetail(id int) (*model.User, error)
	GetDocterDetail(id int) (*model.Doctor, error)
}

type AdminManagementUsecaseImpl struct {
	Repo repository.AdminManagementRepo
}

func NewAdminManagementUsecase(repo repository.AdminManagementRepo) AdminManagementUsecase {
	return &AdminManagementUsecaseImpl{Repo: repo}
}

func (au *AdminManagementUsecaseImpl) GetAllUsers() ([]*model.User, error) {
	users, err := au.Repo.FindAllUsers()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (au *AdminManagementUsecaseImpl) DeleteUsers(id int) (*model.User, error) {
	user, err := au.Repo.DeleteUsers(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (au *AdminManagementUsecaseImpl) GetUserDetail(id int) (*model.User, error) {
	user, err := au.Repo.FindUserDetail(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (au *AdminManagementUsecaseImpl) GetAllDocter() ([]*model.Doctor, error) {
	doctors, err := au.Repo.FindAllDocters()
	if err != nil {
		return nil, err
	}
	return doctors, nil
}

func (au *AdminManagementUsecaseImpl) DeleteDocter(id int) (*model.Doctor, error) {
	doctor, err := au.Repo.DeleteDocter(id)
	if err != nil {
		return nil, err
	}
	return doctor, nil
}

func (au *AdminManagementUsecaseImpl) GetDocterDetail(id int) (*model.Doctor, error) {
	doctor, err := au.Repo.FindDocterDetail(id)
	if err != nil {
		return nil, err
	}
	return doctor, nil
}
