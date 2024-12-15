package usecase

import (
	"calmind/helper"
	"calmind/model"
	repository "calmind/repository/profile"
	"errors"
)

type AdminProfileUseCase interface {
	GetAdminProfile(adminID int) (*model.Admin, error)
	UploadAdminAvatar(adminID int, avatar string, deleteURL string) error
	DeleteAdminAvatar(adminID int) error
}

type AdminProfileUseCaseImpl struct {
	AdminProfileRepo repository.AdminProfileRepository
}

func NewAdminProfileUseCase(repo repository.AdminProfileRepository) AdminProfileUseCase {
	return &AdminProfileUseCaseImpl{AdminProfileRepo: repo}
}

func (u *AdminProfileUseCaseImpl) GetAdminProfile(adminID int) (*model.Admin, error) {
	admin, err := u.AdminProfileRepo.GetByID(adminID)
	if err != nil {
		return nil, errors.New("admin not found")
	}
	return admin, nil
}

func (u *AdminProfileUseCaseImpl) UploadAdminAvatar(adminID int, avatar string, deleteURL string) error {
	err := u.AdminProfileRepo.UpdateAvatarByID(adminID, avatar, deleteURL)
	if err != nil {
		return errors.New("failed to upload avatar")
	}
	return nil
}

func (u *AdminProfileUseCaseImpl) DeleteAdminAvatar(adminID int) error {
	admin, err := u.AdminProfileRepo.GetByID(adminID)
	if err != nil {
		return errors.New("admin not found")
	}

	// Hapus avatar dari ImgBB
	if admin.Avatar != "" {
		err = helper.DeleteFromImgBB(admin.Avatar)
		if err != nil {
			return errors.New("failed to delete avatar from ImgBB")
		}
	}

	err = u.AdminProfileRepo.ClearAvatarByID(adminID)
	if err != nil {
		return errors.New("failed to clear avatar in database")
	}
	return nil
}
