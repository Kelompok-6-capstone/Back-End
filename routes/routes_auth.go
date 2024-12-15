package routes

import (
	controller_artikel "calmind/controller/artikel"
	controller_auth "calmind/controller/authentikasi"
	controller_konsultasi "calmind/controller/konsultasi"
	controller_profil "calmind/controller/profile"
	controller_user_fitur "calmind/controller/user_fitur"

	"github.com/labstack/echo/v4"
)

// Routes untuk User
func UserAuthRoutes(e *echo.Echo, authController *controller_auth.AuthController) {
	e.POST("/user/register", authController.RegisterUser) // Daftar User
	e.POST("/user/login", authController.LoginUser)       // Login User
	e.POST("/user/verify-otp", authController.VerifyOtp)  // Verifikasi OTP
	e.POST("/user/resend-otp", authController.ResendOtp)  // Kirim ulang OTP
	e.GET("/user/logout", authController.LogoutUser)      // Logout User
}

func UserProfil(
	e *echo.Group,
	profilController *controller_profil.ProfilController,
	fitur *controller_user_fitur.UserFiturController,
	konsultasi *controller_konsultasi.ConsultationController,
	artikelController *controller_artikel.ArtikelController,
) {
	// Endpoint untuk profil pengguna
	e.GET("/profile", profilController.GetProfile)            // Melihat profil pengguna
	e.PUT("/profile", profilController.UpdateProfile)         // Memperbarui profil pengguna
	e.POST("/upload-avatar", profilController.UploadAvatar)   // Upload avatar
	e.DELETE("/delete-avatar", profilController.DeleteAvatar) // Hapus avatar

	// Endpoint untuk fitur dokter
	e.GET("/doctors", fitur.GetDoctors)                // Mendapatkan daftar semua dokter
	e.GET("/doctors/tag", fitur.GetDoctorsByTag)       // Mendapatkan dokter berdasarkan tag
	e.GET("/doctors/status", fitur.GetDoctorsByStatus) // Mendapatkan dokter berdasarkan status
	e.GET("/doctors/search", fitur.SearchDoctors)      // Mencari dokter berdasarkan query
	e.GET("/doctors/:id", fitur.GetDoctorDetail)       // Mendapatkan detail dokter berdasarkan ID
	e.GET("/tags", fitur.GetAllTags)                   // Mendapatkan semua tag (bidang keahlian)

	// Endpoint untuk title
	e.GET("/titles", fitur.GetAllTitles)             // Mendapatkan semua title
	e.GET("/doctors/title", fitur.GetDoctorsByTitle) // Mendapatkan dokter berdasarkan title

	e.POST("/consultations", konsultasi.CreateConsultation)            // Membuat konsultasi
	e.GET("/consultations", konsultasi.GetUserConsultations)           // Mendapatkan semua konsultasi user
	e.GET("/consultations/:id", konsultasi.GetUserConsultationDetails) // Mendapatkan detail konsultasi user

	// Endpoint untuk artikel
	e.GET("/artikel", artikelController.GetAllArtikel)      // Mendapatkan semua artikel
	e.GET("/artikel/:id", artikelController.GetArtikelByID) // Mendapatkan detail artikel berdasarkan ID
	e.GET("/artikel/search", artikelController.SearchArtikel)
}
