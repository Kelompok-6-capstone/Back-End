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

	if configDB.Host == "" || configDB.User == "" || configDB.Password == "" || configDB.Port == "" || configDB.Name == "" {
		return nil, fmt.Errorf("konfigurasi database tidak lengkap, periksa file .env Anda")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		configDB.User,
		configDB.Password,
		configDB.Host,
		configDB.Port,
		configDB.Name,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return nil, fmt.Errorf("gagal membuka koneksi ke database: %w", err)
	}
	fmt.Println("Koneksi database berhasil.")

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
		&model.ChatMessage{},
	}

	for _, model := range models {
		if !db.Migrator().HasTable(model) {
			if err := db.AutoMigrate(model); err != nil {
				return nil, fmt.Errorf("gagal melakukan migrasi untuk model %T: %w", model, err)
			}
			fmt.Printf("Migrasi berhasil untuk model: %T\n", model)
		} else {
			fmt.Printf("Tabel untuk model %T sudah ada, tidak dilakukan migrasi.\n", model)
		}
	}

	if err := SeedTitles(db); err != nil {
		fmt.Println("Error seeding titles:", err)
	}

	if err := SeedSpecialties(db); err != nil {
		fmt.Println("Error seeding specialties:", err)
	}

	return db, nil
}
