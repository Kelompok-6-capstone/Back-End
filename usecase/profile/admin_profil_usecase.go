package usecase

import (
	"calmind/helper"
	"calmind/model"
	repository "calmind/repository/profile"
	"errors"
	"path/filepath"
	"strings"
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

func (u *AdminProfileUseCaseImpl) UploadAdminAvatar(adminID int, avatarURL string, _ string) error {
	err := u.AdminProfileRepo.UpdateAvatarByID(adminID, avatarURL)
	if err != nil {
		return errors.New("failed to update avatar in database")
	}
	return nil
}
func (u *AdminProfileUseCaseImpl) DeleteAdminAvatar(adminID int) error {
	admin, err := u.AdminProfileRepo.GetByID(adminID)
	if err != nil {
		return errors.New("admin not found")
	}

	if admin.Avatar != "" {
		// Ekstrak public ID dari Cloudinary URL
		parts := strings.Split(admin.Avatar, "/")
		publicID := strings.TrimSuffix(parts[len(parts)-1], filepath.Ext(admin.Avatar))

		// Hapus file dari Cloudinary
		if err := helper.DeleteFileFromCloudinary(publicID); err != nil {
			return errors.New("failed to delete avatar from Cloudinary")
		}
	}

	// Clear kolom avatar di database
	if err := u.AdminProfileRepo.ClearAvatarByID(adminID); err != nil {
		return errors.New("failed to clear avatar in database")
	}

	return nil
}
