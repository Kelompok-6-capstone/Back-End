package routes

import (
	controller_management "calmind/controller/admin_management"
	controller_artikel "calmind/controller/artikel"
	controller_auth "calmind/controller/authentikasi"
	controller_konsultasi "calmind/controller/konsultasi"
	controller_notifikasi "calmind/controller/midtrans_notifikasi"
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
	e.GET("/consultations", consultationController.GetConsultationsForDoctor)             // Mendapatkan semua konsultasi pasien dokter
	e.GET("/consultations/:id", consultationController.ViewConsultationDetails)           // Mendapatkan detail konsultasi tertentu
	e.POST("/consultations/:id/recommendation", consultationController.AddRecommendation) // Menambahkan rekomendasi pada konsultasi
	e.GET("/consultations/search", consultationController.SearchConsultationsByName)      // Melihat pasien sesuai nama dengan search

}

// Routes untuk Admin
func AdminAuthRoutes(e *echo.Echo, authController *controller_auth.AdminAuthController) {
	e.POST("/admin/login", authController.LoginAdmin)  // Login Admin
	e.GET("/admin/logout", authController.LogoutAdmin) // Logout Admin
}

func AdminManagementRoutes(e *echo.Group, adminManagement *controller_management.AdminManagementController, artikelController *controller_artikel.ArtikelController, consultationController *controller_konsultasi.ConsultationController, profil *controller_profil.AdminController) {
	// user
	e.GET("/allusers", adminManagement.GetAllUsers)    // Ambil Semua Data User
	e.DELETE("/users/:id", adminManagement.DeleteUser) // Hapus User berdasarkan ID

	e.GET("/profil", profil.GetAdminProfile)            // Mendapatkan profil dokter
	e.POST("/upload-image", profil.UploadAdminAvatar)   // Upload avatar dokter
	e.DELETE("/delete-image", profil.DeleteAdminAvatar) // Hapus avatar dokter

	// dokter
	e.GET("/alldocters", adminManagement.GetAllDoctors)    // Ambil Semua Data Dokter
	e.DELETE("/docters/:id", adminManagement.DeleteDoctor) // Hapus Dokter berdasarkan ID

	// artikel
	e.POST("/artikel", artikelController.CreateArtikel)                     // Tambah artikel
	e.GET("/artikel", artikelController.GetAllArtikel)                      // Lihat semua artikel
	e.GET("/artikel/:id", artikelController.GetArtikelByID)                 // Lihat detail artikel
	e.PUT("/artikel/:id", artikelController.UpdateArtikel)                  // Update artikel
	e.DELETE("/artikel/:id", artikelController.DeleteArtikel)               // Hapus artikel
	e.POST("/artikel/upload-image", artikelController.UploadArtikelImage)   // Upload image untuk artikel
	e.DELETE("/artikel/delete-image", artikelController.DeleteArtikelImage) // Hapus image artikel

	// konsultasi
	e.GET("/consultations", consultationController.GetAllStatusConsultations)
	e.GET("/consultations/:id", consultationController.ViewPendingConsultation)
	e.GET("/consultations/pending", consultationController.GetPendingConsultations)
	e.GET("/consultations/approve", consultationController.GetAproveConsultations)
	e.PUT("/consultations/:id/approve", consultationController.ApprovePaymentAndConsultation)
}

func WebhookRoutes(e *echo.Echo, notifikasi *controller_notifikasi.MidtransNotificationController) {
	e.POST("/notifications/midtrans", notifikasi.MidtransNotification)
}
