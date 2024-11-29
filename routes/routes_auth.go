package routes

import (
	"calmind/controller"

	"github.com/labstack/echo/v4"
)

// Routes untuk User
func UserAuthRoutes(e *echo.Echo, authController *controller.AuthController) {
	e.POST("/user/register", authController.RegisterUser) // Daftar User
	e.POST("/user/login", authController.LoginUser)       // Login User
	e.GET("/user/logout", authController.LogoutUser)      // Logout User
}

func UserProfil(e *echo.Group, profilController *controller.ProfilController) {
	e.GET("/profile", profilController.GetProfile) // Profil User
}

// Routes untuk Admin
func AdminAuthRoutes(e *echo.Echo, authController *controller.AdminAuthController) {
	e.POST("/admin/login", authController.LoginAdmin)  // Login Admin
	e.GET("/admin/logout", authController.LogoutAdmin) // Logout Admin
}

func AdminManagementRoutes(e *echo.Group, adminManagement *controller.AdminManagementController) {
	e.GET("/admin/allusers", adminManagement.GetAllUsers)    // Ambil Semua Data User
	e.DELETE("/admin/users/:id", adminManagement.DeleteUser) // Hapus User berdasarkan ID
}

// Routes untuk Doctor
func DoctorAuthRoutes(e *echo.Echo, authController *controller.DoctorAuthController) {
	e.POST("/doctor/register", authController.RegisterDoctor)
	e.POST("/doctor/login", authController.LoginDoctor)
	e.GET("/doctor/logout", authController.LogoutDoctor)
}
