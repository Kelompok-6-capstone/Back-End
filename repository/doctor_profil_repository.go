package repository

import (
	"calmind/model"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type DoctorProfilRepository interface {
	GetByID(id int) (*model.Doctor, error)
	UpdateByID(id int, doctor *model.Doctor) (*model.Doctor, error)
	UpdateDoctorActiveStatus(id int, isActive bool) error
	GetTagByID(id int) (*model.Tags, error)
	UpdateTagsByName(doctorID int, tagNames []string) error
	GetDoctorTitleByID(doctorID int) (*model.Title, error)
	UpdateDoctorTitleByName(doctorID int, titleName string) error
}

type DoctorProfilRepositoryImpl struct {
	DB *gorm.DB
}

func NewDoctorProfilRepository(db *gorm.DB) DoctorProfilRepository {
	return &DoctorProfilRepositoryImpl{DB: db}
}

// Mendapatkan dokter berdasarkan ID
func (r *DoctorProfilRepositoryImpl) GetByID(id int) (*model.Doctor, error) {
	var doctor model.Doctor
	err := r.DB.Preload("Tags").Preload("Title").Where("id = ?", id).First(&doctor).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("doctor with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to fetch doctor: %v", err)
	}
	return &doctor, nil
}
func (r *DoctorProfilRepositoryImpl) UpdateByID(id int, doctor *model.Doctor) (*model.Doctor, error) {
	var existingDoctor model.Doctor
	err := r.DB.Where("id = ?", id).First(&existingDoctor).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("doctor with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to fetch doctor: %v", err)
	}

	// Update hanya field yang diberikan
	if doctor.Username != "" {
		existingDoctor.Username = doctor.Username
	}
	if doctor.NoHp != "" {
		existingDoctor.NoHp = doctor.NoHp
	}
	if doctor.Avatar != "" {
		existingDoctor.Avatar = doctor.Avatar
	}
	if doctor.DateOfBirth != "" {
		existingDoctor.DateOfBirth = doctor.DateOfBirth
	}
	if doctor.Address != "" {
		existingDoctor.Address = doctor.Address
	}
	if doctor.Schedule != "" {
		existingDoctor.Schedule = doctor.Schedule
	}
	if doctor.Experience > 0 {
		existingDoctor.Experience = doctor.Experience
	}
	if doctor.STRNumber != "" {
		existingDoctor.STRNumber = doctor.STRNumber
	}
	if doctor.About != "" {
		existingDoctor.About = doctor.About
	}

	// Force update JenisKelamin jika diberikan
	if doctor.JenisKelamin != "" {
		if doctor.JenisKelamin != "Laki-laki" && doctor.JenisKelamin != "Perempuan" {
			return nil, fmt.Errorf("invalid value for JenisKelamin: %s", doctor.JenisKelamin)
		}
		existingDoctor.JenisKelamin = doctor.JenisKelamin
	}

	// Simpan perubahan
	err = r.DB.Save(&existingDoctor).Error
	if err != nil {
		return nil, fmt.Errorf("failed to update doctor profile: %v", err)
	}

	return &existingDoctor, nil
}

// Memperbarui status aktif dokter
func (r *DoctorProfilRepositoryImpl) UpdateDoctorActiveStatus(id int, isActive bool) error {
	err := r.DB.Model(&model.Doctor{}).Where("id = ?", id).UpdateColumn("is_active", isActive).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("doctor with ID %d not found", id)
		}
		return fmt.Errorf("failed to update active status for doctor with ID %d: %v", id, err)
	}
	return nil
}

// Mendapatkan tag berdasarkan ID
func (r *DoctorProfilRepositoryImpl) GetTagByID(id int) (*model.Tags, error) {
	var tag model.Tags
	err := r.DB.Where("id = ?", id).First(&tag).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("tag with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to fetch tag: %v", err)
	}
	return &tag, nil
}

// Memperbarui tag yang terkait dengan dokter berdasarkan nama
func (r *DoctorProfilRepositoryImpl) UpdateTagsByName(doctorID int, tagNames []string) error {
	var tags []model.Tags
	err := r.DB.Where("name IN ?", tagNames).Find(&tags).Error
	if err != nil {
		return fmt.Errorf("failed to fetch tags from database: %v", err)
	}

	if len(tags) != len(tagNames) {
		return errors.New("some tags were not found in the database")
	}

	var doctor model.Doctor
	err = r.DB.Where("id = ?", doctorID).First(&doctor).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("doctor with ID %d not found", doctorID)
		}
		return fmt.Errorf("failed to fetch doctor: %v", err)
	}

	// Perbarui asosiasi Tags
	err = r.DB.Model(&doctor).Association("Tags").Replace(tags)
	if err != nil {
		return fmt.Errorf("failed to update tags: %v", err)
	}

	return nil
}

// Mendapatkan title dokter berdasarkan ID dokter
func (r *DoctorProfilRepositoryImpl) GetDoctorTitleByID(doctorID int) (*model.Title, error) {
	var doctor model.Doctor
	err := r.DB.Preload("Title").Where("id = ?", doctorID).First(&doctor).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("doctor with ID %d not found", doctorID)
		}
		return nil, fmt.Errorf("failed to fetch doctor title: %v", err)
	}
	return &doctor.Title, nil
}

// Memperbarui title dokter berdasarkan nama
func (r *DoctorProfilRepositoryImpl) UpdateDoctorTitleByName(doctorID int, titleName string) error {
	var title model.Title
	err := r.DB.Where("name = ?", titleName).First(&title).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("title '%s' not found", titleName)
		}
		return fmt.Errorf("failed to fetch title: %v", err)
	}

	var doctor model.Doctor
	err = r.DB.Where("id = ?", doctorID).First(&doctor).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("doctor with ID %d not found", doctorID)
		}
		return fmt.Errorf("failed to fetch doctor: %v", err)
	}

	doctor.TitleID = title.ID

	// Simpan perubahan TitleID
	err = r.DB.Save(&doctor).Error
	if err != nil {
		return fmt.Errorf("failed to update doctor title: %v", err)
	}

	return nil
}
