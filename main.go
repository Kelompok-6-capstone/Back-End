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

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

	//    Repositori, usecase, dan controller untuk Consultasi
	artikelonRepo := repository.NewArtikelRepository(DB)
	artikelUsecase := usecase.NewArtikelUsecase(artikelonRepo)
	artikelController := controller.NewArtikelController(artikelUsecase)

	// Middleware
	jwtMiddleware := middlewares.NewJWTMiddleware(jwtSecret)

	// Echo instance
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{
			"http://localhost:3000",
			"http://localhost:5173",
			"https://jovial-mooncake-23a3d0.netlify.app", // Tambahkan domain frontend Anda
		},
		AllowMethods:  []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
		AllowHeaders:  []string{echo.HeaderAuthorization, echo.HeaderContentType},
		ExposeHeaders: []string{"Authorization"},
		// AllowCredentials: true, // Izinkan pengiriman cookie atau token di header
	}))

	// routes auth
	routes.UserAuthRoutes(e, userController)               // user
	routes.AdminAuthRoutes(e, adminController)             // admin
	routes.DoctorAuthRoutes(e, doctorControllerManagement) // dokter

	// Group untuk user, dengan middleware yang memastikan hanya user yang login dapat mengaksesnya
	userGroup := e.Group("/user", jwtMiddleware.HandlerUser)
	adminGroup := e.Group("/admin", jwtMiddleware.HandlerAdmin)
	doctorGroup := e.Group("/doctor", jwtMiddleware.HandlerDoctor)

	// Routing group auth
	routes.UserProfil(userGroup, userProfilController, userFiturController, consultationController, artikelController) // Profil User
	routes.AdminManagementRoutes(adminGroup, adminControllerManagement, artikelController)                             // Admin management
	routes.DoctorProfil(doctorGroup, doctorProfilController, artikelController, consultationController)                // Doctor Profile

	// Mulai server
	log.Fatal(e.Start(":8000"))
}
