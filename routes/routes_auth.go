package routes

import (
	"calmind/controller"

	"github.com/labstack/echo/v4"
)

// Routes untuk User
func UserAuthRoutes(e *echo.Group, authController *controller.AuthController) {
	e.POST("/register", authController.RegisterUser) // Daftar User
	e.POST("/login", authController.LoginUser)       // Login User
	e.GET("/logout", authController.LogoutUser)      // Logout User
}

// Routes untuk Admin
func AdminAuthRoutes(e *echo.Echo, authController *controller.AdminAuthController) {
	e.POST("/admin/login", authController.LoginAdmin)  // Login Admin
	e.GET("/admin/logout", authController.LogoutAdmin) // Logout Admin
}
func AdminManagementRoutes(e *echo.Group, adminmanagement *controller.AdminManagementController) {
	e.GET("admin/allusers", adminmanagement.GetAllUsers)    // Ambil Semua Data User
	e.DELETE("admin/users/:id", adminmanagement.DeleteUser) // Hapus User
}

func DoctorAuthRoutes(e *echo.Group, authController *controller.DoctorAuthController) {
	e.POST("/register", authController.RegisterDoctor)
	e.POST("/login", authController.LoginDoctor)
	e.GET("/logout", authController.LogoutDoctor)
}
