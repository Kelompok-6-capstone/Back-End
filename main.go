package main

import (
	"calmind/config"
	"calmind/controller"
	"calmind/middlewares"
	"calmind/repository"
	"calmind/routes"
	"calmind/service"
	"calmind/usecase"
	"log"

	"net/http"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Gagal memuat file .env")
	}

	// Inisialisasi database
	DB, err := config.InitDB()
	if err != nil {
		log.Fatalf("Gagal menginisialisasi database: %v", err)
	}

	// Konfigurasi JWT
	jwtSecret := config.NewJWTConfig()
	jwtService := service.NewJWTService(jwtSecret)
	otpService := service.NewOtpService()

	// Repositori, usecase, dan controller untuk User
	userRepo := repository.NewAuthRepository(DB)
	otpRepo := repository.NewOtpRepository(DB)
	userUsecase := usecase.NewAuthUsecase(userRepo, jwtService, otpRepo, otpService)
	userController := controller.NewAuthController(userUsecase)

	// Repositori, usecase, dan controller untuk Admin
	adminRepo := repository.NewAdminAuthRepository(DB)
	adminUsecase := usecase.NewAdminAuthUsecase(adminRepo, jwtService)
	adminController := controller.NewAdminAuthController(adminUsecase)

	// Repositori, usecase, dan controller untuk Admin management
	adminRepoManagement := repository.NewAdminManagementRepo(DB)
	adminUsecaseManagement := usecase.NewAdminManagementUsecase(adminRepoManagement)
	adminControllerManagement := controller.NewAdminManagementController(adminUsecaseManagement)

	// Repositori, usecase, dan controller untuk dokter
	doctorRepoManagement := repository.NewDoctorAuthRepository(DB)
	doctorUsecaseManagement := usecase.NewDoctorAuthUsecase(doctorRepoManagement, jwtService, otpRepo, otpService)
	doctorControllerManagement := controller.NewDoctorAuthController(doctorUsecaseManagement)

	//	Repositori, usecase, dan controller untuk Profil User
	userProfilRepo := repository.NewUserProfilRepository(DB)
	userProfilUsecase := usecase.NewUserProfileUseCase(userProfilRepo)
	userProfilController := controller.NewProfilController(userProfilUsecase)

	//	Repositori, usecase, dan controller untuk Fitur User
	userFiturRepo := repository.NewUserFiturRepository(DB)
	userFiturUsecase := usecase.NewUserFiturUsecase(userFiturRepo)
	userFiturController := controller.NewUserFiturController(userFiturUsecase)

	//	Repositori, usecase, dan controller untuk Profil doctor
	doctorProfilRepo := repository.NewDoctorProfilRepository(DB)
	doctorProfilUsecase := usecase.NewDoctorProfileUseCase(doctorProfilRepo)
	doctorProfilController := controller.NewDoctorProfileController(doctorProfilUsecase)

	//    Repositori, usecase, dan controller untuk Consultasi
	consultationRepo := repository.NewConsultationRepository(DB)
	consultationUsecase := usecase.NewConsultationUsecase(consultationRepo)
	consultationController := controller.NewConsultationController(consultationUsecase)

	// Middleware
	jwtMiddleware := middlewares.NewJWTMiddleware(jwtSecret)

	// Echo instance
	e := echo.New()

	// Middleware untuk CORS Dinamis
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			origin := c.Request().Header.Get("Origin")
			if origin != "" {
				c.Response().Header().Set("Access-Control-Allow-Origin", origin) // Izinkan origin spesifik
				c.Response().Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
				c.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
				c.Response().Header().Set("Access-Control-Allow-Credentials", "true")
			}

			// Tangani preflight request
			if c.Request().Method == http.MethodOptions {
				return c.NoContent(http.StatusOK)
			}

			return next(c)
		}
	})

	// routes auth
	routes.UserAuthRoutes(e, userController)               // user
	routes.AdminAuthRoutes(e, adminController)             // admin
	routes.DoctorAuthRoutes(e, doctorControllerManagement) // dokter

	// Group untuk user, dengan middleware yang memastikan hanya user yang login dapat mengaksesnya
	userGroup := e.Group("/user", jwtMiddleware.HandlerUser)
	adminGroup := e.Group("/admin", jwtMiddleware.HandlerAdmin)
	doctorGroup := e.Group("/doctor", jwtMiddleware.HandlerDoctor)

	// Routing group auth
	routes.UserProfil(userGroup, userProfilController, userFiturController, consultationController) // Profil User
	routes.AdminManagementRoutes(adminGroup, adminControllerManagement)                             // Admin management
	routes.DoctorProfil(doctorGroup, doctorProfilController)                                        // Doctor Profile

	// Mulai server
	log.Fatal(e.Start(":8000"))
}
