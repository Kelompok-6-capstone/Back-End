package routes

import (
	controller_management "calmind/controller/admin_management"
	controller_artikel "calmind/controller/artikel"
	controller_auth "calmind/controller/authentikasi"
	controller_konsultasi "calmind/controller/konsultasi"
	controller_profil "calmind/controller/profile"

	"github.com/labstack/echo/v4"
)

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
