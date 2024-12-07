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
	// Memuat file .env
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

	// Inisialisasi Repositori, Usecase, dan Controller

	// User
	userRepo := repository.NewAuthRepository(DB)
	otpRepo := repository.NewOtpRepository(DB)
	userUsecase := usecase.NewAuthUsecase(userRepo, jwtService, otpRepo, otpService)
	userController := controller.NewAuthController(userUsecase)

	// Admin
	adminRepo := repository.NewAdminAuthRepository(DB)
	adminUsecase := usecase.NewAdminAuthUsecase(adminRepo, jwtService)
	adminController := controller.NewAdminAuthController(adminUsecase)

	adminRepoManagement := repository.NewAdminManagementRepo(DB)
	adminUsecaseManagement := usecase.NewAdminManagementUsecase(adminRepoManagement)
	adminControllerManagement := controller.NewAdminManagementController(adminUsecaseManagement)

	// Dokter
	doctorRepo := repository.NewDoctorAuthRepository(DB)
	doctorUsecase := usecase.NewDoctorAuthUsecase(doctorRepo, jwtService, otpRepo, otpService)
	doctorController := controller.NewDoctorAuthController(doctorUsecase)

	// Profil
	userProfilRepo := repository.NewUserProfilRepository(DB)
	userProfilUsecase := usecase.NewUserProfileUseCase(userProfilRepo)
	userProfilController := controller.NewProfilController(userProfilUsecase)

	doctorProfilRepo := repository.NewDoctorProfilRepository(DB)
	doctorProfilUsecase := usecase.NewDoctorProfileUseCase(doctorProfilRepo)
	doctorProfilController := controller.NewDoctorProfileController(doctorProfilUsecase)

	// Fitur User
	userFiturRepo := repository.NewUserFiturRepository(DB)
	userFiturUsecase := usecase.NewUserFiturUsecase(userFiturRepo)
	userFiturController := controller.NewUserFiturController(userFiturUsecase)

	// Konsultasi
	consultationRepo := repository.NewConsultationRepository(DB)
	consultationUsecase := usecase.NewConsultationUsecase(consultationRepo)
	consultationController := controller.NewConsultationController(consultationUsecase)

	// Artikel
	artikelRepo := repository.NewArtikelRepository(DB)
	artikelUsecase := usecase.NewArtikelUsecase(artikelRepo)
	artikelController := controller.NewArtikelController(artikelUsecase)

	// Pembayaran
	paymentRepo := repository.NewPaymentRepository(DB)
	paymentUsecase := usecase.NewPaymentUsecase(paymentRepo)
	paymentController := controller.NewPaymentController(paymentUsecase, consultationUsecase)

	// Middleware
	jwtMiddleware := middlewares.NewJWTMiddleware(jwtSecret)

	// Echo Instance
	e := echo.New()

	// Konfigurasi Middleware
	e.Static("/uploads", "uploads")
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:5174", "http://localhost:5173", "http://127.0.0.1:5500", "https://jovial-mooncake-23a3d0.netlify.app"},
		AllowMethods:     []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS},
		AllowHeaders:     []string{echo.HeaderAuthorization, echo.HeaderContentType},
		AllowCredentials: true,
	}))

	// Routes untuk User
	routes.UserAuthRoutes(e, userController)

	// Routes untuk Admin
	routes.AdminAuthRoutes(e, adminController)

	// Routes untuk Dokter
	routes.DoctorAuthRoutes(e, doctorController)

	// Protected Routes
	userGroup := e.Group("/user", jwtMiddleware.HandlerUser)
	routes.UserProfil(userGroup, userProfilController, userFiturController, consultationController, artikelController)

	adminGroup := e.Group("/admin", jwtMiddleware.HandlerAdmin)
	routes.AdminManagementRoutes(adminGroup, adminControllerManagement, consultationController, artikelController)

	doctorGroup := e.Group("/doctor", jwtMiddleware.HandlerDoctor)
	routes.DoctorProfil(doctorGroup, doctorProfilController, artikelController, consultationController)

	// Jalankan Server
	log.Fatal(e.Start(":8000"))
}
