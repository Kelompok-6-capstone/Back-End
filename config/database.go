package config

import (
	"calmind/model"
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ConfigDB struct {
	Host     string
	User     string
	Password string
	Port     string
	Name     string
}

func InitDB() (*gorm.DB, error) {
	configDB := ConfigDB{
		Host:     os.Getenv("DATABASE_HOST"),
		User:     os.Getenv("DATABASE_USER"),
		Password: os.Getenv("DATABASE_PASSWORD"),
		Port:     os.Getenv("DATABASE_PORT"),
		Name:     os.Getenv("DATABASE_NAME"),
	}

	// Validasi konfigurasi database
	if configDB.Host == "" || configDB.User == "" || configDB.Password == "" || configDB.Port == "" || configDB.Name == "" {
		return nil, fmt.Errorf("konfigurasi database tidak lengkap, periksa file .env Anda")
	}

	// Format DSN
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local",
		configDB.User,
		configDB.Password,
		configDB.Host,
		configDB.Port,
		configDB.Name)

	// Buka koneksi ke database
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("gagal membuka koneksi ke database: %w", err)
	}

	// Nonaktifkan pemeriksaan foreign key sementara (jika perlu)
	if err := db.Exec("SET FOREIGN_KEY_CHECKS=0;").Error; err != nil {
		return nil, fmt.Errorf("gagal menonaktifkan foreign key checks: %w", err)
	}

	// Migrasi model
	models := []interface{}{
		&model.User{},
		&model.Admin{},
		&model.Doctor{},
		&model.Otp{},
		&model.Consultation{},
		&model.Tags{},
		&model.Artikel{},
		&model.ChatLog{},
		&model.Rekomendasi{},
	}

	for _, model := range models {
		if err := db.AutoMigrate(model); err != nil {
			return nil, fmt.Errorf("gagal melakukan migrasi untuk model %T: %w", model, err)
		}
	}

	// Aktifkan kembali foreign key checks
	if err := db.Exec("SET FOREIGN_KEY_CHECKS=1;").Error; err != nil {
		return nil, fmt.Errorf("gagal mengaktifkan kembali foreign key checks: %w", err)
	}

	// Seed Titles
	if err := SeedTitles(db); err != nil {
		fmt.Println("Error seeding titles:", err)
	} else {
		fmt.Println("Titles seeded successfully.")
	}

	// Seed Tags (Specialties)
	if err := SeedSpecialties(db); err != nil {
		fmt.Println("Error seeding specialties:", err)
	} else {
		fmt.Println("Specialties seeded successfully.")
	}

	return db, nil
}
