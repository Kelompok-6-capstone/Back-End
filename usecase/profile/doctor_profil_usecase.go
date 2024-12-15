package usecase

import (
	"calmind/model"
	repository "calmind/repository/profile"
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

// Mendapatkan profil dokter berdasarkan ID
func (u *doctorProfileUseCaseImpl) GetDoctorProfile(doctorID int) (*model.Doctor, error) {
	doctor, err := u.DoctorProfileRepo.GetByID(doctorID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch doctor profile: %v", err)
	}
	return doctor, nil
}
func (u *doctorProfileUseCaseImpl) UpdateDoctorProfile(doctorID int, doctor *model.Doctor) (*model.Doctor, error) {
	// Memperbarui Tags jika disediakan
	if len(doctor.Tags) > 0 {
		var tagNames []string
		for _, tag := range doctor.Tags {
			tagNames = append(tagNames, tag.Name)
		}

		err := u.DoctorProfileRepo.UpdateTagsByName(doctorID, tagNames)
		if err != nil {
			return nil, fmt.Errorf("failed to update tags: %v", err)
		}
	}

	// Memperbarui Title jika disediakan
	if doctor.Title.Name != "" {
		err := u.DoctorProfileRepo.UpdateDoctorTitleByName(doctorID, doctor.Title.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to update title: %v", err)
		}
	}

	// Memperbarui profil dokter lainnya
	updatedDoctor, err := u.DoctorProfileRepo.UpdateByID(doctorID, doctor)
	if err != nil {
		return nil, fmt.Errorf("failed to update doctor profile: %v", err)
	}
	return updatedDoctor, nil
}

// Mengatur status aktif dokter
func (u *doctorProfileUseCaseImpl) SetDoctorActiveStatus(doctorID int, isActive bool) error {
	err := u.DoctorProfileRepo.UpdateDoctorActiveStatus(doctorID, isActive)
	if err != nil {
		return fmt.Errorf("failed to update doctor active status: %v", err)
	}
	return nil
}
