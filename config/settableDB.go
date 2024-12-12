package config

import (
	"calmind/model"
	"fmt"

	"gorm.io/gorm"
)

// Fungsi untuk melakukan seed data Titles
func SeedTitles(db *gorm.DB) error {
	titles := []model.Title{
		{Name: "Psikiatri Anak dan Remaja", Gambar: "https://example.com/images/psikiatri-anak-remaja.jpg"},
		{Name: "Psikiatri Umum", Gambar: "https://example.com/images/psikiatri-umum.jpg"},
		{Name: "Psikiatri Geriatri", Gambar: "https://example.com/images/psikiatri-geriatri.jpg"},
		{Name: "Psikoterapi", Gambar: "https://example.com/images/psikoterapi.jpg"},
		{Name: "Konsultasi Keluarga", Gambar: "https://example.com/images/konsultasi-keluarga.jpg"},
		{Name: "Neuropsikiatri", Gambar: "https://example.com/images/neuropsikiatri.jpg"},
		{Name: "Psikiatri Komunitas", Gambar: "https://example.com/images/psikiatri-komunitas.jpg"},
		{Name: "Psikologi Klinis", Gambar: "https://example.com/images/psikologi-klinis.jpg"},
		{Name: "Rehabilitasi Psikiatri", Gambar: "https://example.com/images/rehabilitasi-psikiatri.jpg"},
		{Name: "Psikologi Pendidikan", Gambar: "https://example.com/images/psikologi-pendidikan.jpg"},
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
		{Name: "Stress", Gambar: "https://example.com/images/stress.jpg"},
		{Name: "Depresi", Gambar: "https://example.com/images/depresi.jpg"},
		{Name: "Trauma", Gambar: "https://example.com/images/trauma.jpg"},
		{Name: "Adiksi", Gambar: "https://example.com/images/adiksi.jpg"},
		{Name: "Gangguan Kecemasan", Gambar: "https://example.com/images/kecemasan.jpg"},
		{Name: "Pengembangan Diri", Gambar: "https://example.com/images/pengembangan-diri.jpg"},
		{Name: "Gangguan Mood", Gambar: "https://example.com/images/mood.jpg"},
		{Name: "Pengasuhan & Anak", Gambar: "https://example.com/images/pengasuhan-anak.jpg"},
		{Name: "Pekerjaan", Gambar: "https://example.com/images/pekerjaan.jpg"},
		{Name: "Hubungan & Keluarga", Gambar: "https://example.com/images/hubungan-keluarga.jpg"},
		{Name: "Identitas Seksual", Gambar: "https://example.com/images/identitas-seksual.jpg"},
		{Name: "Gangguan Kepribadian", Gambar: "https://example.com/images/kepribadian.jpg"},
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
