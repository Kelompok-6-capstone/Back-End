package routes

import (
	"calmind/controller"

	"github.com/labstack/echo/v4"
)

func AdminAuthRoutes(e *echo.Group, authController *controller.AdminAuthController) {
	e.POST("/login", authController.LoginAdmin)
	e.GET("/logout", authController.LogoutAdmin)
}
func UserAuthRoutes(e *echo.Group, authController *controller.AuthController) {
	e.POST("/register", authController.RegisterUser)
	e.POST("/login", authController.LoginUser)
	e.GET("/logout", authController.LogoutUser)
}

func DoctorAuthRoutes(e *echo.Group, authController *controller.DoctorAuthController) {
	e.POST("/register", authController.RegisterDoctor)
	e.POST("/login", authController.LoginDoctor)
	e.GET("/logout", authController.LogoutDoctor)
}
