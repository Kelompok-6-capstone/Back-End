package usecase

import (
	"calmind/model"
	"calmind/repository"
	"errors"
)

type UserProfileUseCase interface {
	GetUserProfile(UserID int) (*model.User, error)
	UpdateUserProfile(UserID int, User *model.User) (*model.User, error)
}

type UserProfileUseCaseImpl struct {
	UserProfileRepo repository.UserProfilRepository
}

func NewUserProfileUseCase(repo repository.UserProfilRepository) UserProfileUseCase {
	return &UserProfileUseCaseImpl{UserProfileRepo: repo}
}

// GetUserProfile retrieves the profile of a User by their ID
func (u *UserProfileUseCaseImpl) GetUserProfile(UserID int) (*model.User, error) {
	User, err := u.UserProfileRepo.GetByID(UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return User, nil
}

// UpdateUserProfile updates the profile of a User
func (u *UserProfileUseCaseImpl) UpdateUserProfile(UserID int, User *model.User) (*model.User, error) {
	updatedUser, err := u.UserProfileRepo.UpdateByID(UserID, User)
	if err != nil {
		return nil, errors.New("failed to update user profile")
	}
	return updatedUser, nil
}
