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

// GetDoctorProfile retrieves the profile of a doctor by their ID
func (u *doctorProfileUseCaseImpl) GetDoctorProfile(doctorID int) (*model.Doctor, error) {
	doctor, err := u.DoctorProfileRepo.GetByID(doctorID)
	if err != nil {
		return nil, errors.New("doctor not found")
	}
	return doctor, nil
}

// UpdateDoctorProfile updates the profile of a doctor
func (u *doctorProfileUseCaseImpl) UpdateDoctorProfile(doctorID int, doctor *model.Doctor) (*model.Doctor, error) {
	updatedDoctor, err := u.DoctorProfileRepo.UpdateByID(doctorID, doctor)
	if err != nil {
		return nil, errors.New("failed to update doctor profile")
	}
	return updatedDoctor, nil
}

// SetDoctorActiveStatus updates the active status of a doctor
func (u *doctorProfileUseCaseImpl) SetDoctorActiveStatus(doctorID int, isActive bool) error {
	err := u.DoctorProfileRepo.UpdateDoctorActiveStatus(doctorID, isActive)
	if err != nil {
		return errors.New("failed to update doctor active status")
	}
	return nil
}
