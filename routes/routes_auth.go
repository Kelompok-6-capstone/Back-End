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

// Routes untuk Profil dan Fitur User
func UserProfilRoutes(
	e *echo.Group,
	profilController *controller.ProfilController,
	fiturController *controller.UserFiturController,
	konsultasiController *controller.ConsultationController,
	artikelController *controller.ArtikelController,
	chatController *controller.ChatController,
	custServiceController *controller.CustServiceController,
) {
	// Profil pengguna
	e.GET("/profile", profilController.GetProfile)            // Lihat profil pengguna
	e.PUT("/profile", profilController.UpdateProfile)         // Perbarui profil pengguna
	e.POST("/upload-avatar", profilController.UploadAvatar)   // Upload avatar
	e.DELETE("/delete-avatar", profilController.DeleteAvatar) // Hapus avatar

	// Fitur dokter
	e.GET("/doctors", fiturController.GetDoctors)                // Semua dokter
	e.GET("/doctors/tag", fiturController.GetDoctorsByTag)       // Dokter berdasarkan tag
	e.GET("/doctors/status", fiturController.GetDoctorsByStatus) // Dokter berdasarkan status
	e.GET("/doctors/search", fiturController.SearchDoctors)      // Cari dokter berdasarkan query
	e.GET("/doctors/:id", fiturController.GetDoctorDetail)       // Detail dokter
	e.GET("/tags", fiturController.GetAllTags)                   // Semua tag
	e.GET("/titles", fiturController.GetAllTitles)               // Semua title

	// Konsultasi
	e.POST("/consultations", konsultasiController.CreateConsultation)            // Buat konsultasi baru
	e.GET("/consultations", konsultasiController.GetUserConsultations)           // Semua konsultasi pengguna
	e.GET("/consultations/:id", konsultasiController.GetUserConsultationDetails) // Detail konsultasi

	// Artikel
	e.GET("/artikel", artikelController.GetAllArtikel)      // Semua artikel
	e.GET("/artikel/:id", artikelController.GetArtikelByID) // Detail artikel
	e.GET("/artikel/search", artikelController.SearchArtikel)

	// Chat
	e.POST("/chat/send", chatController.SendChat)                  // Kirim pesan
	e.GET("/chat/history/:room_id", chatController.GetChatHistory) // Riwayat chat berdasarkan room

	// cs
	e.POST("/customer-service", custServiceController.GetResponse)
	e.GET("/customer-service", custServiceController.GetQuestion)
}

// Routes untuk Auth Dokter
func DoctorAuthRoutes(e *echo.Echo, authController *controller.DoctorAuthController) {
	e.POST("/doctor/register", authController.RegisterDoctor) // Daftar Dokter
	e.POST("/doctor/login", authController.LoginDoctor)       // Login Dokter
	e.GET("/doctor/logout", authController.LogoutDoctor)      // Logout Dokter
	e.POST("/doctor/verify-otp", authController.VerifyOtp)    // Verifikasi OTP
	e.POST("/doctor/resend-otp", authController.ResendOtp)    // Kirim ulang OTP
}

// Routes untuk Profil dan Fitur Dokter
func DoctorProfilRoutes(
	e *echo.Group,
	profilController *controller.DoctorProfileController,
	consultationController *controller.ConsultationController,
	userFiturController *controller.UserFiturController,
	chatController *controller.ChatController,
	custServiceController *controller.CustServiceController,
) {
	// Endpoint untuk Profil Dokter
	e.GET("/profile", profilController.GetProfile)
	e.PUT("/profile", profilController.UpdateProfile)
	e.PUT("/status", profilController.SetActiveStatus)
	e.POST("/upload-image", profilController.UploadAvatar)
	e.DELETE("/delete-image", profilController.DeleteAvatar)

	// Konsultasi
	e.GET("/consultations", consultationController.GetConsultationsForDoctor)
	e.GET("/consultations/:id", consultationController.ViewConsultationDetails)

	// Chat
	e.POST("/chat/send", chatController.SendChat)
	e.GET("/chat/history/:room_id", chatController.GetChatHistory)

	// cs
	e.POST("/customer-service", custServiceController.GetResponse)
	e.GET("/customer-service", custServiceController.GetQuestion)
}

// Routes untuk Admin
func AdminAuthRoutes(e *echo.Echo, authController *controller.AdminAuthController) {
	e.POST("/admin/login", authController.LoginAdmin)  // Login Admin
	e.GET("/admin/logout", authController.LogoutAdmin) // Logout Admin
}

// Routes untuk Manajemen Admin
func AdminManagementRoutes(
	e *echo.Group,
	adminManagementController *controller.AdminManagementController,
	artikelController *controller.ArtikelController,
	konsultasiController *controller.ConsultationController,
) {
	// User management
	e.GET("/allusers", adminManagementController.GetAllUsers)    // Semua user
	e.DELETE("/users/:id", adminManagementController.DeleteUser) // Hapus user berdasarkan ID

	// Doctor management
	e.GET("/alldoctors", adminManagementController.GetAllDoctors)    // Semua dokter
	e.DELETE("/doctors/:id", adminManagementController.DeleteDoctor) // Hapus dokter berdasarkan ID

	// Artikel management
	e.POST("/artikel", artikelController.CreateArtikel)                     // Tambah artikel
	e.GET("/artikel", artikelController.GetAllArtikel)                      // Semua artikel
	e.GET("/artikel/:id", artikelController.GetArtikelByID)                 // Detail artikel
	e.PUT("/artikel/:id", artikelController.UpdateArtikel)                  // Perbarui artikel
	e.DELETE("/artikel/:id", artikelController.DeleteArtikel)               // Hapus artikel
	e.POST("/artikel/upload-image", artikelController.UploadArtikelImage)   // Upload gambar artikel
	e.DELETE("/artikel/delete-image", artikelController.DeleteArtikelImage) // Hapus gambar artikel

	// Konsultasi management
	e.GET("/consultations", konsultasiController.GetAllStatusConsultations)                 // Semua konsultasi
	e.GET("/consultations/:id", konsultasiController.ViewPendingConsultation)               // Konsultasi tertunda
	e.GET("/consultations/pending", konsultasiController.GetPendingConsultations)           // Konsultasi pending
	e.GET("/consultations/approved", konsultasiController.GetAproveConsultations)           // Konsultasi approved
	e.PUT("/consultations/:id/approve", konsultasiController.ApprovePaymentAndConsultation) // Approve konsultasi
}
