package routes

import (
	controller_artikel "calmind/controller/artikel"
	controller_auth "calmind/controller/authentikasi"
	controller_konsultasi "calmind/controller/konsultasi"
	controller_profil "calmind/controller/profile"
	controller_user_fitur "calmind/controller/user_fitur"

	"github.com/labstack/echo/v4"
)

// Routes untuk Doctor
func DoctorAuthRoutes(e *echo.Echo, authController *controller_auth.DoctorAuthController) {
	e.POST("/doctor/register", authController.RegisterDoctor) // Daftar Dokter
	e.POST("/doctor/login", authController.LoginDoctor)       // Login Dokter
	e.GET("/doctor/logout", authController.LogoutDoctor)      // Logout Dokter
	e.POST("/doctor/verify-otp", authController.VerifyOtp)    // Verifikasi OTP
	e.POST("/doctor/resend-otp", authController.ResendOtp)    // Kirim ulang OTP
}

func DoctorProfil(
	e *echo.Group,
	profilController *controller_profil.DoctorProfileController,
	artikelController *controller_artikel.ArtikelController,
	consultationController *controller_konsultasi.ConsultationController,
	fiturController *controller_user_fitur.UserFiturController,
) {
	// Profil Dokter
	e.GET("/profile", profilController.GetProfile)           // Mendapatkan profil dokter
	e.PUT("/profile", profilController.UpdateProfile)        // Mengupdate profil dokter
	e.PUT("/status", profilController.SetActiveStatus)       // Mengubah status aktif/tidak aktif dokter
	e.POST("/upload-image", profilController.UploadAvatar)   // Upload avatar dokter
	e.DELETE("/delete-image", profilController.DeleteAvatar) // Hapus avatar dokter

	// Tags dan Titles
	e.GET("/tags", fiturController.GetAllTags)     // Mendapatkan semua tag (bidang keahlian)
	e.GET("/titles", fiturController.GetAllTitles) // Mendapatkan semua title

	// Konsultasi
	e.GET("/consultations", consultationController.GetAllConsultationsForDoctor)          // Mendapatkan semua konsultasi pasien dokter
	e.GET("/consultations/:id", consultationController.ViewConsultationDetails)           // Mendapatkan detail konsultasi tertentu
	e.POST("/consultations/:id/recommendation", consultationController.AddRecommendation) // Menambahkan rekomendasi pada konsultasi
	e.GET("/consultations/search", consultationController.SearchConsultationsByName)      // Melihat pasien sesuai nama dengan search
}
