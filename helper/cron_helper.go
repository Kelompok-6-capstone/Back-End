package helper

import (
	"fmt"

	"github.com/robfig/cron/v3"
)

type ExpiredJob struct {
	MarkExpiredFunc func() error // Fungsi yang dipanggil untuk menandai konsultasi expired
}

func StartExpiredConsultationJob(markExpiredFunc func() error) {
	job := ExpiredJob{MarkExpiredFunc: markExpiredFunc}

	c := cron.New()

	// Jalankan setiap 1 menit
	c.AddFunc("@every 1m", func() {
		if err := job.MarkExpiredFunc(); err != nil {
			fmt.Println("Error marking expired consultations:", err)
		} else {
			fmt.Println("Successfully marked expired consultations.")
		}
	})

	c.Start()
}
