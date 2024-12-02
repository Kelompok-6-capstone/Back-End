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

func UserProfil(e *echo.Group, profilController *controller.ProfilController, fitur *controller.UserFiturController, konsultasi *controller.ConsultationController) {
	e.GET("/profile", profilController.GetProfile)    // Profil User
	e.PUT("/profile", profilController.UpdateProfile) // Profil User
	e.GET("/doctors", fitur.GetDoctors)
	e.GET("/doctors/specialty", fitur.GetDoctorsBySpecialty)
	e.GET("/doctors/status", fitur.GetDoctorsByStatus)
	e.GET("/doctors/search", fitur.SearchDoctors)
	e.GET("/doctors/:id", fitur.GetDoctorDetail)
	e.POST("/consultations", konsultasi.CreateConsultation)

}

// Routes untuk Admin
func AdminAuthRoutes(e *echo.Echo, authController *controller.AdminAuthController) {
	e.POST("/admin/login", authController.LoginAdmin)  // Login Admin
	e.GET("/admin/logout", authController.LogoutAdmin) // Logout Admin
}

func AdminManagementRoutes(e *echo.Group, adminManagement *controller.AdminManagementController) {
	e.GET("/allusers", adminManagement.GetAllUsers)    // Ambil Semua Data User
	e.DELETE("/users/:id", adminManagement.DeleteUser) // Hapus User berdasarkan ID
}

// Routes untuk Doctor
func DoctorAuthRoutes(e *echo.Echo, authController *controller.DoctorAuthController) {
	e.POST("/doctor/register", authController.RegisterDoctor)
	e.POST("/doctor/login", authController.LoginDoctor)
	e.GET("/doctor/logout", authController.LogoutDoctor)
	e.POST("/doctor/verify-otp", authController.VerifyOtp) // verivikasi otp
}

func DoctorProfil(e *echo.Group, profilController *controller.DoctorProfileController) {
	e.GET("/profile", profilController.GetProfile)     // Mendapatkan profil dokter
	e.PUT("/profile", profilController.UpdateProfile)  // Mengupdate profil dokter
	e.PUT("/status", profilController.SetActiveStatus) // Mengubah status aktif/tidak aktif dokter
}
