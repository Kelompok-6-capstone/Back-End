package repository

import (
	"calmind/model"

	"gorm.io/gorm"
)

type DoctorProfilRepository interface {
	GetByID(id int) (*model.Doctor, error)
	UpdateByID(id int, doctor *model.Doctor) (*model.Doctor, error)
	UpdateDoctorActiveStatus(id int, isActive bool) error
	GetTagByID(id int) (*model.Tags, error)
	UpdateTags(doctorID int, tags []model.Tags) error
	GetDoctorTitleByID(doctorID int) (*model.Title, error)
	UpdateDoctorTitle(doctorID int, titleID int) error
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
		return nil, err
	}
	return &doctor, nil
}

// Memperbarui profil dokter berdasarkan ID
func (r *DoctorProfilRepositoryImpl) UpdateByID(id int, doctor *model.Doctor) (*model.Doctor, error) {
	var existingDoctor model.Doctor
	err := r.DB.Where("id = ?", id).First(&existingDoctor).Error
	if err != nil {
		return nil, err
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

	// Simpan perubahan
	err = r.DB.Save(&existingDoctor).Error
	if err != nil {
		return nil, err
	}

	return &existingDoctor, nil
}

// Memperbarui status aktif dokter
func (r *DoctorProfilRepositoryImpl) UpdateDoctorActiveStatus(id int, isActive bool) error {
	return r.DB.Model(&model.Doctor{}).Where("id = ?", id).Update("is_active", isActive).Error
}

// Mendapatkan tag berdasarkan ID
func (r *DoctorProfilRepositoryImpl) GetTagByID(id int) (*model.Tags, error) {
	var tag model.Tags
	err := r.DB.Where("id = ?", id).First(&tag).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

// Memperbarui tag yang terkait dengan dokter
func (r *DoctorProfilRepositoryImpl) UpdateTags(doctorID int, tags []model.Tags) error {
	var doctor model.Doctor
	if err := r.DB.Where("id = ?", doctorID).First(&doctor).Error; err != nil {
		return err
	}

	// Perbarui asosiasi tags dengan dokter
	return r.DB.Model(&doctor).Association("Tags").Replace(tags)
}

// Mendapatkan title dokter berdasarkan ID dokter
func (r *DoctorProfilRepositoryImpl) GetDoctorTitleByID(doctorID int) (*model.Title, error) {
	var doctor model.Doctor
	err := r.DB.Preload("Title").Where("id = ?", doctorID).First(&doctor).Error
	if err != nil {
		return nil, err
	}
	return &doctor.Title, nil
}

// Memperbarui title dokter
func (r *DoctorProfilRepositoryImpl) UpdateDoctorTitle(doctorID int, titleID int) error {
	var doctor model.Doctor
	if err := r.DB.Where("id = ?", doctorID).First(&doctor).Error; err != nil {
		return err
	}

	var title model.Title
	if err := r.DB.Where("id = ?", titleID).First(&title).Error; err != nil {
		return err
	}

	doctor.TitleID = title.ID
	return r.DB.Save(&doctor).Error
}
