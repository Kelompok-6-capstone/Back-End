package routes

import (
	"calmind/controller"

	"github.com/labstack/echo/v4"
)

// Routes untuk User
func UserAuthRoutes(e *echo.Echo, authController *controller.AuthController) {
	e.POST("/user/register", authController.RegisterUser) // Daftar User
	e.POST("/user/login", authController.LoginUser)       // Login User
	e.POST("/user/verify-otp", authController.VerifyOtp)  // Verifikasi OTP
	e.POST("/user/resend-otp", authController.ResendOtp)  // Kirim ulang OTP
	e.GET("/user/logout", authController.LogoutUser)      // Logout User
}

func UserProfil(
	e *echo.Group,
	profilController *controller.ProfilController,
	fitur *controller.UserFiturController,
	consultationController *controller.ConsultationController,
	artikelController *controller.ArtikelController,
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
	e.GET("/titles", fitur.GetAllTitles)               // Mendapatkan semua title
	e.GET("/doctors/title", fitur.GetDoctorsByTitle)   // Mendapatkan dokter berdasarkan title

	// Endpoint untuk konsultasi
	e.POST("/consultations/select-schedule", consultationController.SelectSchedule)       // Memilih jadwal konsultasi
	e.POST("/consultations/add-message", consultationController.AddMessageToConsultation) // Menambahkan keluhan ke konsultasi
	e.POST("/consultations/payment", consultationController.CreatePayment)                // Membuat pembayaran konsultasi
	e.GET("/consultations/status/:id", consultationController.CheckConsultationStatus)    // Mengecek status konsultasi

	// Endpoint untuk artikel
	e.GET("/artikel", artikelController.GetAllArtikel)      // Mendapatkan semua artikel
	e.GET("/artikel/:id", artikelController.GetArtikelByID) // Mendapatkan detail artikel berdasarkan ID
}

// Routes untuk Admin
func AdminAuthRoutes(e *echo.Echo, authController *controller.AdminAuthController) {
	e.POST("/admin/login", authController.LoginAdmin)  // Login Admin
	e.GET("/admin/logout", authController.LogoutAdmin) // Logout Admin
}

func AdminManagementRoutes(e *echo.Group, adminManagement *controller.AdminManagementController, consultationController *controller.ConsultationController, artikelController *controller.ArtikelController) {
	e.GET("/allusers", adminManagement.GetAllUsers)                                          // Ambil semua data user
	e.GET("/alldoctors", adminManagement.GetAllDoctors)                                      // Ambil semua data dokter
	e.DELETE("/users/:id", adminManagement.DeleteUser)                                       // Hapus user berdasarkan ID
	e.DELETE("/doctors/:id", adminManagement.DeleteDoctor)                                   // Hapus dokter berdasarkan ID
	e.PUT("/consultations/verify-payment/:payment_id", consultationController.VerifyPayment) // Verifikasi pembayaran konsultasi
	e.POST("/artikel", artikelController.CreateArtikel)                                      // Tambah artikel
	e.GET("/artikel", artikelController.GetAllArtikel)                                       // Lihat semua artikel
	e.GET("/artikel/:id", artikelController.GetArtikelByID)                                  // Lihat detail artikel
	e.PUT("/artikel/:id", artikelController.UpdateArtikel)                                   // Update artikel
	e.DELETE("/artikel/:id", artikelController.DeleteArtikel)                                // Hapus artikel
	e.POST("/artikel/upload-image", artikelController.UploadArtikelImage)                    // Upload image untuk artikel
	e.DELETE("/artikel/delete-image", artikelController.DeleteArtikelImage)                  // Hapus image artikel
}

// Routes untuk Dokter
func DoctorAuthRoutes(e *echo.Echo, authController *controller.DoctorAuthController) {
	e.POST("/doctor/register", authController.RegisterDoctor) // Register Dokter
	e.POST("/doctor/login", authController.LoginDoctor)       // Login Dokter
	e.GET("/doctor/logout", authController.LogoutDoctor)      // Logout Dokter
	e.POST("/doctor/verify-otp", authController.VerifyOtp)    // Verifikasi OTP
	e.POST("/doctor/resend-otp", authController.ResendOtp)    // Kirim ulang OTP
}

func DoctorProfil(e *echo.Group, profilController *controller.DoctorProfileController, artikelController *controller.ArtikelController, consultationController *controller.ConsultationController) {
	e.GET("/profile", profilController.GetProfile)                                        // Mendapatkan profil dokter
	e.PUT("/profile", profilController.UpdateProfile)                                     // Mengupdate profil dokter
	e.PUT("/status", profilController.SetActiveStatus)                                    // Mengubah status aktif/tidak aktif dokter
	e.GET("/consultations", consultationController.GetConsultationsAllDoctor)             // Mendapatkan daftar pasien konsultasi
	e.GET("/consultations/:id", consultationController.GetConsultationDetail)             // Mendapatkan detail konsultasi pasien
	e.PUT("/consultations/:id/recommendation", consultationController.GiveRecommendation) // Memberikan rekomendasi untuk konsultasi
	e.GET("/artikel", artikelController.GetAllArtikel)                                    // Mendapatkan semua artikel
	e.GET("/artikel/:id", artikelController.GetArtikelByID)                               // Mendapatkan detail artikel berdasarkan ID
	e.POST("/upload-image", profilController.UploadAvatar)                                // Upload avatar dokter
	e.DELETE("/delete-image", profilController.DeleteAvatar)                              // Hapus avatar dokter
}
