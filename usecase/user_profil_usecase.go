package usecase

import (
	"calmind/model"
	"calmind/repository"
	"errors"
)

type UserProfilUsecase interface {
	GetUserProfile(userID int) (*model.User, error)
	UpdateUserProfile(userID int, user *model.User) (*model.User, error)
}

type UserProfilUsecaseImpl struct {
	UserProfilRepo repository.UserProfilRepository
}

func NewUserProfilUsecaseImpl(repo repository.UserProfilRepository) UserProfilUsecase {
	return &UserProfilUsecaseImpl{UserProfilRepo: repo}
}

func (u *UserProfilUsecaseImpl) GetUserProfile(userID int) (*model.User, error) {
	user, err := u.UserProfilRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (u *UserProfilUsecaseImpl) UpdateUserProfile(userID int, user *model.User) (*model.User, error) {
	// Perform the update
	updatedUser, err := u.UserProfilRepo.UpdateByID(userID, user)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return updatedUser, nil
}
