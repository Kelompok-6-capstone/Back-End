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
	// Handle specialties if provided
	if len(doctor.Specialties) > 0 {
		var specialties []model.Specialty
		for _, specialty := range doctor.Specialties {
			var dbSpecialty model.Specialty
			if err := u.DoctorProfileRepo.GetSpecialtyByID(specialty.ID, &dbSpecialty); err != nil {
				return nil, err
			}
			specialties = append(specialties, dbSpecialty)
		}

		if err := u.DoctorProfileRepo.UpdateSpecialties(doctorID, specialties); err != nil {
			return nil, err
		}
	}

	return u.DoctorProfileRepo.UpdateByID(doctorID, doctor)
}

func (u *doctorProfileUseCaseImpl) SetDoctorActiveStatus(doctorID int, isActive bool) error {
	err := u.DoctorProfileRepo.UpdateDoctorActiveStatus(doctorID, isActive)
	if err != nil {
		return errors.New("failed to update doctor active status")
	}
	return nil
}
