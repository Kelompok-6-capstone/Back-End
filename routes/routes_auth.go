package routes

import (
	"calmind/controller"

	"github.com/labstack/echo/v4"
)

// Routes untuk User
func UserAuthRoutes(e *echo.Echo, authController *controller.AuthController) {
	e.POST("/user/register", authController.RegisterUser) // Daftar User
	e.POST("/user/login", authController.LoginUser)       // Login User
	e.POST("/user/verify-otp", authController.VerifyOtp)  // verivikasi otp
	e.GET("/user/logout", authController.LogoutUser)      // Logout User
}

func UserProfil(
	e *echo.Group,
	profilController *controller.ProfilController,
	fitur *controller.UserFiturController,
	konsultasi *controller.ConsultationController,
	artikelController *controller.ArtikelController,
) {
	// Endpoint untuk profil pengguna
	e.GET("/profile", profilController.GetProfile)          // Melihat profil pengguna
	e.PUT("/profile", profilController.UpdateProfile)       // Memperbarui profil pengguna
	e.POST("/upload-avatar", profilController.UploadAvatar) // Upload avatar

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

	// Endpoint untuk konsultasi
	e.POST("/consultations", konsultasi.CreateConsultation)          // Membuat konsultasi
	e.GET("/consultation/:id", konsultasi.GetConsultationDetailByID) // Melihat detail konsultasi

	// Endpoint untuk artikel
	e.GET("/artikel", artikelController.GetAllArtikel)      // Mendapatkan semua artikel
	e.GET("/artikel/:id", artikelController.GetArtikelByID) // Mendapatkan detail artikel berdasarkan ID
}

// Routes untuk Admin
func AdminAuthRoutes(e *echo.Echo, authController *controller.AdminAuthController) {
	e.POST("/admin/login", authController.LoginAdmin)  // Login Admin
	e.GET("/admin/logout", authController.LogoutAdmin) // Logout Admin
}

func AdminManagementRoutes(e *echo.Group, adminManagement *controller.AdminManagementController, artikelController *controller.ArtikelController) {
	e.GET("/allusers", adminManagement.GetAllUsers)           // Ambil Semua Data User
	e.GET("/alldocters", adminManagement.GetAllDoctors)       // Ambil Semua Data User
	e.DELETE("/users/:id", adminManagement.DeleteUser)        // Hapus User berdasarkan ID
	e.DELETE("/docters/:id", adminManagement.DeleteDoctor)    // Hapus User berdasarkan ID
	e.POST("/artikel", artikelController.CreateArtikel)       // Tambah artikel
	e.GET("/artikel", artikelController.GetAllArtikel)        // Lihat semua artikel
	e.GET("/artikel/:id", artikelController.GetArtikelByID)   // Lihat detail artikel
	e.PUT("/artikel/:id", artikelController.UpdateArtikel)    // Update artikel
	e.DELETE("/artikel/:id", artikelController.DeleteArtikel) // Hapus artikel
}

// Routes untuk Doctor
func DoctorAuthRoutes(e *echo.Echo, authController *controller.DoctorAuthController) {
	e.POST("/doctor/register", authController.RegisterDoctor)
	e.POST("/doctor/login", authController.LoginDoctor)
	e.GET("/doctor/logout", authController.LogoutDoctor)
	e.POST("/doctor/verify-otp", authController.VerifyOtp) // verivikasi otp
	e.POST("/resend-otp", authController.ResendOtp)        // Kirim ulang OTP

}

func DoctorProfil(e *echo.Group, profilController *controller.DoctorProfileController, artikelController *controller.ArtikelController, consultationController *controller.ConsultationController) {
	e.GET("/profile", profilController.GetProfile)                                       // Mendapatkan profil dokter
	e.PUT("/profile", profilController.UpdateProfile)                                    // Mengupdate profil dokter
	e.PUT("/status", profilController.SetActiveStatus)                                   // Mengubah status aktif/tidak aktif dokter
	e.GET("/consultations", consultationController.GetConsultationsAllDoctor)            // Melihat Daftar pasien
	e.GET("/consultation/:id", consultationController.GetConsultationDetailByID)         // Melihat Detail keluhan
	e.PUT("/consultation/:id/recommendation", consultationController.GiveRecommendation) // Memberikan rekomendasi
	e.GET("/artikel", artikelController.GetAllArtikel)                                   // Lihat semua artikel
	e.GET("/artikel/:id", artikelController.GetArtikelByID)                              // Lihat detail artikel
}
