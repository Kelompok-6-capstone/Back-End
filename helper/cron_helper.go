package helper

import (
	"log"
	"time"

	"github.com/robfig/cron/v3"
)

type ExpiredJob struct {
	MarkExpiredFunc func() error // Fungsi yang dipanggil untuk menandai konsultasi expired
}

func StartExpiredConsultationJob(markExpiredFunc func() error) {
	c := cron.New()

	// Jalankan setiap 1 menit
	_, err := c.AddFunc("@every 1m", func() {
		log.Printf("Checking for expired consultations at %v", time.Now())
		if err := markExpiredFunc(); err != nil {
			log.Printf("Error checking expired consultations: %v", err)
		}
	})
	if err != nil {
		log.Printf("Failed to add cron job: %v", err)
		return
	}

	c.Start()
}
