package usecase

import (
	"calmind/model"
	"calmind/repository"
)

type AdminManagementUsecase interface {
	GetAllUsers() ([]*model.User, error)
	DeleteUsers(id int) (*model.User, error)
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
