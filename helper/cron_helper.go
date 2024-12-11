package helper

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

type ExpiredJob struct {
	MarkExpiredFunc func() error // Fungsi yang dipanggil untuk menandai konsultasi expired
}

func StartExpiredConsultationJob(markExpiredFunc func() error) {
	job := ExpiredJob{MarkExpiredFunc: markExpiredFunc}

	c := cron.New()

	// Jalankan setiap 1 menit
	_, err := c.AddFunc("@every 1m", func() {
		fmt.Printf("Marking expired consultations at %v\n", time.Now())
		if err := job.MarkExpiredFunc(); err != nil {
			fmt.Printf("Error marking expired consultations: %v\n", err)
		}
	})
	if err != nil {
		fmt.Printf("Failed to add cron job: %v\n", err)
		return
	}

	c.Start()
}
