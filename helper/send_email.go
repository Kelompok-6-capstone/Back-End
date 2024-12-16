package helper

import (
	"fmt"
	"os"

	"gopkg.in/gomail.v2"
)

// SendEmail mengirimkan email verifikasi dengan kode OTP ke alamat email yang diberikan.
func SendEmail(email, otpCode string) error {
	// Informasi pengirim diambil dari variabel lingkungan
	from := os.Getenv("SMTP_EMAIL")
	password := os.Getenv("SMTP_PASSWORD")

	// Pastikan variabel lingkungan SMTP_EMAIL dan SMTP_PASSWORD telah diatur
	if from == "" || password == "" {
		return fmt.Errorf("SMTP credentials are not set in environment variables")
	}

	// Subjek email
	subject := "Verifikasi Akun Anda - OTP Calmind"

	// Isi email dalam format HTML
	body := `
		<!DOCTYPE html>
		<html>
		<head>
			<title>Verifikasi Akun Anda</title>
			<style>
				body {
					font-family: Arial, sans-serif;
					line-height: 1.6;
					color: #333;
					background-color: #f9f9f9;
					padding: 20px;
				}
				.container {
					max-width: 600px;
					margin: 0 auto;
					background: #fff;
					border-radius: 8px;
					box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
					padding: 20px;
				}
				.header {
					text-align: center;
					color: #4CAF50;
					margin-bottom: 20px;
				}
				.otp {
					font-size: 24px;
					font-weight: bold;
					text-align: center;
					color: #4CAF50;
					margin: 20px 0;
				}
				.footer {
					margin-top: 30px;
					text-align: center;
					font-size: 12px;
					color: #aaa;
				}
			</style>
		</head>
		<body>
			<div class="container">
				<h2 class="header">Verifikasi Akun Anda</h2>
				<p>Terima kasih telah mendaftar di <strong>Calmind</strong>. Untuk menyelesaikan proses pendaftaran, silakan gunakan kode OTP berikut:</p>
				<div class="otp">` + otpCode + `</div>
				<p>Kode ini hanya berlaku selama <strong>5 menit</strong>. Jika Anda tidak meminta kode ini, silakan abaikan email ini.</p>
				<p>Terima kasih,</p>
				<p><strong>Tim Calmind</strong></p>
			</div>
			<div class="footer">
				&copy; 2024 Calmind. All rights reserved.
			</div>
		</body>
		</html>
	`

	// Konfigurasi SMTP Gmail
	smtpHost := "smtp.gmail.com"
	smtpPort := 587

	// Membuat pesan email
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", email)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	// Membuat koneksi ke server SMTP Gmail
	d := gomail.NewDialer(smtpHost, smtpPort, from, password)

	// Mengirim email
	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
