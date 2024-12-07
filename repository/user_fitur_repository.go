package repository

import (
	"calmind/model"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type UserFiturRepository interface {
	GetAllDoctors() ([]model.Doctor, error)
	GetDoctorsByTag(tag string) ([]model.Doctor, error)
	GetDoctorsByStatus(isActive bool) ([]model.Doctor, error)
	SearchDoctors(query string) ([]model.Doctor, error)
	GetDoctorByID(id int) (*model.Doctor, error)
	GetTags() ([]model.Tags, error)
	GetTitles() ([]model.Title, error)
	GetDoctorsByTitle(title string) ([]model.Doctor, error)
}

type UserFiturRepositoryImpl struct {
	DB *gorm.DB
}

func NewUserFiturRepository(db *gorm.DB) UserFiturRepository {
	return &UserFiturRepositoryImpl{DB: db}
}

// Mendapatkan semua dokter yang memenuhi kriteria umum
func (r *UserFiturRepositoryImpl) GetAllDoctors() ([]model.Doctor, error) {
	var doctors []model.Doctor
	err := r.DB.
		Preload("Tags").
		Preload("Title").
		Where("price > 0 AND experience > 0 AND is_verified = true AND is_active = true").
		Find(&doctors).Error

	// Logging data untuk debugging
	if err != nil {
		fmt.Printf("Error fetching doctors: %v\n", err)
		return nil, err
	}
	if len(doctors) == 0 {
		fmt.Println("No doctors found.")
		return doctors, nil
	}

	return doctors, nil
}

// Mendapatkan dokter berdasarkan Tag
func (r *UserFiturRepositoryImpl) GetDoctorsByTag(tag string) ([]model.Doctor, error) {
	var tags model.Tags
	if err := r.DB.Where("name = ?", tag).First(&tags).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("tag '%s' not found", tag)
		}
		return nil, fmt.Errorf("failed to fetch tag: %v", err)
	}

	var doctors []model.Doctor
	err := r.DB.
		Joins("JOIN doctor_tags ON doctors.id = doctor_tags.doctor_id").
		Joins("JOIN tags ON doctor_tags.tag_id = tags.id").
		Where("tags.name = ?", tag).
		Preload("Tags").
		Preload("Title").
		Find(&doctors).Error

	return doctors, err
}

// Mendapatkan dokter berdasarkan status (aktif atau tidak)
func (r *UserFiturRepositoryImpl) GetDoctorsByStatus(isActive bool) ([]model.Doctor, error) {
	var doctors []model.Doctor
	err := r.DB.Where("is_active = ?", isActive).Find(&doctors).Error
	return doctors, err
}

// Mencari dokter berdasarkan kueri di beberapa kolom
func (r *UserFiturRepositoryImpl) SearchDoctors(query string) ([]model.Doctor, error) {
	var doctors []model.Doctor
	err := r.DB.
		Where("username LIKE ? OR title LIKE ? OR about LIKE ?", "%"+query+"%", "%"+query+"%", "%"+query+"%").
		Find(&doctors).Error
	return doctors, err
}

// Mendapatkan dokter berdasarkan ID, termasuk informasi Tags
func (r *UserFiturRepositoryImpl) GetDoctorByID(id int) (*model.Doctor, error) {
	var doctor model.Doctor
	err := r.DB.Preload("Tags").Preload("Title").Where("id = ?", id).First(&doctor).Error
	if err != nil {
		return nil, err
	}
	return &doctor, nil
}

// Mendapatkan daftar semua Tags
func (r *UserFiturRepositoryImpl) GetTags() ([]model.Tags, error) {
	var tags []model.Tags
	err := r.DB.Find(&tags).Error
	return tags, err
}

func (r *UserFiturRepositoryImpl) GetTitles() ([]model.Title, error) {
	var titles []model.Title
	err := r.DB.Find(&titles).Error
	return titles, err
}

// Mendapatkan dokter berdasarkan Title
func (r *UserFiturRepositoryImpl) GetDoctorsByTitle(title string) ([]model.Doctor, error) {
	// Validasi apakah title ada di database
	var existingTitle model.Title
	if err := r.DB.Where("name = ?", title).First(&existingTitle).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("title '%s' not found", title)
		}
		return nil, fmt.Errorf("failed to fetch title: %v", err)
	}

	// Jika title ditemukan, ambil dokter dengan title tersebut
	var doctors []model.Doctor
	err := r.DB.
		Where("title_id = ?", existingTitle.ID).
		Preload("Title").
		Preload("Tags").
		Find(&doctors).Error

	if err != nil {
		return nil, fmt.Errorf("failed to fetch doctors for title '%s': %v", title, err)
	}

	return doctors, nil
}
