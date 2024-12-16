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

func (u *AdminProfileUseCaseImpl) UploadAdminAvatar(adminID int, avatar string, deleteURL string) error {
	err := u.AdminProfileRepo.UpdateAvatarByID(adminID, avatar, deleteURL)
	if err != nil {
		return errors.New("failed to upload avatar")
	}
	return nil
}

func (u *AdminProfileUseCaseImpl) DeleteAdminAvatar(adminID int) error {
	// Ambil data admin berdasarkan ID
	admin, err := u.AdminProfileRepo.GetByID(adminID)
	if err != nil {
		return errors.New("admin not found")
	}

	// Hapus avatar dari Cloudinary jika ada
	if admin.Avatar != "" {
		// Ekstrak public_id dari URL avatar
		parts := strings.Split(admin.Avatar, "/")
		publicID := strings.TrimSuffix(parts[len(parts)-1], filepath.Ext(admin.Avatar))

		// Hapus file dari Cloudinary
		err = helper.DeleteFileFromCloudinary(publicID)
		if err != nil {
			return errors.New("failed to delete avatar from Cloudinary")
		}
	}

	// Kosongkan kolom avatar di database
	err = u.AdminProfileRepo.ClearAvatarByID(adminID)
	if err != nil {
		return errors.New("failed to clear avatar in database")
	}

	return nil
}
