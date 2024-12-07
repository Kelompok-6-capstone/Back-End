package usecase

import (
	"calmind/model"
	"calmind/repository"
	"fmt"
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

// GetDoctorProfile retrieves a doctor's profile by their ID.
func (u *doctorProfileUseCaseImpl) GetDoctorProfile(doctorID int) (*model.Doctor, error) {
	doctor, err := u.DoctorProfileRepo.GetByID(doctorID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch doctor profile: %v", err)
	}
	return doctor, nil
}

// UpdateDoctorProfile updates a doctor's profile, including tags and title if provided.
func (u *doctorProfileUseCaseImpl) UpdateDoctorProfile(doctorID int, doctor *model.Doctor) (*model.Doctor, error) {
	// Validate and update tags if provided
	if len(doctor.Tags) > 0 {
		tagNames := make([]string, 0, len(doctor.Tags))
		for _, tag := range doctor.Tags {
			if tag.Name == "" {
				return nil, fmt.Errorf("tag name cannot be empty")
			}
			tagNames = append(tagNames, tag.Name)
		}

		if err := u.DoctorProfileRepo.UpdateTagsByName(doctorID, tagNames); err != nil {
			return nil, fmt.Errorf("failed to update tags: %v", err)
		}
	}

	// Validate and update title if provided
	if doctor.Title.Name != "" {
		if err := u.DoctorProfileRepo.UpdateDoctorTitleByName(doctorID, doctor.Title.Name); err != nil {
			return nil, fmt.Errorf("failed to update title: %v", err)
		}
	}

	// Update other profile fields
	updatedDoctor, err := u.DoctorProfileRepo.UpdateByID(doctorID, doctor)
	if err != nil {
		return nil, fmt.Errorf("failed to update doctor profile: %v", err)
	}

	return updatedDoctor, nil
}

// SetDoctorActiveStatus updates the active status of a doctor.
func (u *doctorProfileUseCaseImpl) SetDoctorActiveStatus(doctorID int, isActive bool) error {
	err := u.DoctorProfileRepo.UpdateDoctorActiveStatus(doctorID, isActive)
	if err != nil {
		return fmt.Errorf("failed to update doctor active status: %v", err)
	}
	return nil
}
