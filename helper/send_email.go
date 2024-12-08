package helper

import (
	"os"

	"gopkg.in/gomail.v2"
)

func SendEmail(email, otpCode string) error {
	// Informasi pengirim
	from := os.Getenv("SMTP_EMAIL")
	password := os.Getenv("SMTP_PASSWORD")

	// Subjek dan isi email
	subject := "Verifikasi Akun Anda - OTP Calmind"
	body := `
		<!DOCTYPE html>
		<html>
		<head>
			<title>Verifikasi Akun Anda</title>
		</head>
		<body style="font-family: Arial, sans-serif; line-height: 1.6;">
			<h2>Halo,</h2>
			<p>Terima kasih telah mendaftar di <strong>Calmind</strong>. Untuk menyelesaikan proses pendaftaran, silakan gunakan kode OTP berikut:</p>
			<h1 style="text-align: center; color: #4CAF50;">` + otpCode + `</h1>
			<p>Kode ini hanya berlaku selama <strong>5 menit</strong>. Jika Anda tidak meminta kode ini, silakan abaikan email ini.</p>
			<p>Terima kasih,</p>
			<p><strong>Tim Calmind</strong></p>
		</body>
		</html>
	`

	// Konfigurasi SMTP Gmail
	smtpHost := "smtp.gmail.com"
	smtpPort := 587

	// Membuat pesan email
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", email) // Pastikan parameter `email` digunakan di sini
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	// Membuat koneksi ke server SMTP Gmail
	d := gomail.NewDialer(smtpHost, smtpPort, from, password)

	// Kirim email
	err := d.DialAndSend(m)
	return err
}
