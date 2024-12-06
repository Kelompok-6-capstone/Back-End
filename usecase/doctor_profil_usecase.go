package usecase

import (
	"calmind/model"
	"calmind/repository"
	"errors"
)

type DoctorProfileUseCase interface {
	GetDoctorProfile(doctorID int) (*model.Doctor, error)
	UpdateDoctorProfile(doctorID int, doctor *model.Doctor) (*model.Doctor, error)
	SetDoctorActiveStatus(doctorID int, isActive bool) error
}

type doctorProfileUseCaseImpl struct {
	DoctorProfileRepo repository.DoctorProfilRepository
}

func NewDoctorProfileUseCase(repo repository.DoctorProfilRepository) DoctorProfileUseCase {
	return &doctorProfileUseCaseImpl{DoctorProfileRepo: repo}
}

func (u *doctorProfileUseCaseImpl) GetDoctorProfile(doctorID int) (*model.Doctor, error) {
	doctor, err := u.DoctorProfileRepo.GetByID(doctorID)
	if err != nil {
		return nil, errors.New("doctor not found")
	}
	return doctor, nil
}

func (u *doctorProfileUseCaseImpl) UpdateDoctorProfile(doctorID int, doctor *model.Doctor) (*model.Doctor, error) {
	// Memperbarui Tags jika disediakan
	if len(doctor.Tags) > 0 {
		var tags []model.Tags
		for _, tag := range doctor.Tags {
			dbTag, err := u.DoctorProfileRepo.GetTagByID(tag.ID)
			if err != nil {
				return nil, errors.New("invalid tag provided")
			}
			tags = append(tags, *dbTag)
		}

		if err := u.DoctorProfileRepo.UpdateTags(doctorID, tags); err != nil {
			return nil, errors.New("failed to update tags")
		}
	}

	// Memperbarui Title jika disediakan
	if doctor.TitleID != 0 {
		_, err := u.DoctorProfileRepo.GetDoctorTitleByID(doctor.TitleID)
		if err != nil {
			return nil, errors.New("invalid title provided")
		}

		if err := u.DoctorProfileRepo.UpdateDoctorTitle(doctorID, doctor.TitleID); err != nil {
			return nil, errors.New("failed to update title")
		}
	}

	// Memperbarui profil dokter lainnya
	updatedDoctor, err := u.DoctorProfileRepo.UpdateByID(doctorID, doctor)
	if err != nil {
		return nil, errors.New("failed to update doctor profile")
	}
	return updatedDoctor, nil
}

func (u *doctorProfileUseCaseImpl) SetDoctorActiveStatus(doctorID int, isActive bool) error {
	err := u.DoctorProfileRepo.UpdateDoctorActiveStatus(doctorID, isActive)
	if err != nil {
		return errors.New("failed to update doctor active status")
	}
	return nil
}
