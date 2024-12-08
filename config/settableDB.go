package config

import (
	"calmind/model"
	"fmt"

	"gorm.io/gorm"
)

// Fungsi untuk melakukan seed data Titles
func SeedTitles(db *gorm.DB) error {
	titles := []model.Title{
		{Name: "Psikiatri Anak dan Remaja"},
		{Name: "Psikiatri Umum"},
		{Name: "Psikiatri Geriatri"},
		{Name: "Psikoterapi"},
		{Name: "Konsultasi Keluarga"},
		{Name: "Neuropsikiatri"},
		{Name: "Psikiatri Komunitas"},
		{Name: "Psikologi Klinis"},
		{Name: "Rehabilitasi Psikiatri"},
		{Name: "Psikologi Pendidikan"},
	}

	fmt.Println("Seeding Titles...")

	for _, title := range titles {
		var existingTitle model.Title
		err := db.Where("name = ?", title.Name).First(&existingTitle).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return fmt.Errorf("error saat memeriksa title %s: %w", title.Name, err)
		}

		// Jika title belum ada, tambahkan ke database
		if existingTitle.ID == 0 {
			if err := db.Create(&title).Error; err != nil {
				return fmt.Errorf("gagal menambahkan title %s: %w", title.Name, err)
			}
			fmt.Printf("Title %s berhasil ditambahkan.\n", title.Name)
		} else {
			fmt.Printf("Title %s sudah ada di database.\n", title.Name)
		}
	}

	fmt.Println("Seeding Titles selesai.")
	return nil
}

// Fungsi untuk melakukan seed data Specialties
func SeedSpecialties(db *gorm.DB) error {
	tags := []model.Tags{
		{Name: "Stress"},
		{Name: "Depresi"},
		{Name: "Trauma"},
		{Name: "Adiksi"},
		{Name: "Gangguan Kecemasan"},
		{Name: "Pengembangan Diri"},
		{Name: "Gangguan Mood"},
		{Name: "Pengasuhan & Anak"},
		{Name: "Pekerjaan"},
		{Name: "Hubungan & Keluarga"},
		{Name: "Identitas Seksual"},
		{Name: "Gangguan Kepribadian"},
	}

	fmt.Println("Seeding Specialties (Tags)...")

	for _, tag := range tags {
		var existingTag model.Tags
		err := db.Where("name = ?", tag.Name).First(&existingTag).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return fmt.Errorf("error saat memeriksa tag %s: %w", tag.Name, err)
		}

		// Jika tag belum ada, tambahkan ke database
		if existingTag.ID == 0 {
			if err := db.Create(&tag).Error; err != nil {
				return fmt.Errorf("gagal menambahkan tag %s: %w", tag.Name, err)
			}
			fmt.Printf("Tag %s berhasil ditambahkan.\n", tag.Name)
		} else {
			fmt.Printf("Tag %s sudah ada di database.\n", tag.Name)
		}
	}

	fmt.Println("Seeding Specialties selesai.")
	return nil
}
