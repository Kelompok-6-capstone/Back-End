package main

import (
	"calmind/config"
	"calmind/controller"
	"calmind/helper"
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
	consultationRepo := repository.NewConsultationRepositoryImpl(DB)
	consultationUsecase := usecase.NewConsultationUsecaseImpl(consultationRepo)
	consultationController := controller.NewConsultationController(consultationUsecase)

	helper.StartExpiredConsultationJob(consultationUsecase.MarkExpiredConsultations)

	//    Repositori, usecase, dan controller untuk Consultasi
	artikelonRepo := repository.NewArtikelRepository(DB)
	artikelUsecase := usecase.NewArtikelUsecase(artikelonRepo)
	artikelController := controller.NewArtikelController(artikelUsecase)

	//    Repositori, usecase, dan controller untuk chatbot ai
	chatbotRepo := repository.NewChatLogRepository(DB)
	chatbotUsecase := usecase.NewChatbotUsecase(chatbotRepo)
	chatbotController := controller.NewChatbotController(chatbotUsecase)

	// Middleware
	jwtMiddleware := middlewares.NewJWTMiddleware(jwtSecret)

	e := echo.New()
	e.Static("/uploads", "uploads")
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:5174", "http://localhost:5173", "http://127.0.0.1:5500", "https://jovial-mooncake-23a3d0.netlify.app"},
		AllowMethods:     []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS},
		AllowHeaders:     []string{echo.HeaderAuthorization, echo.HeaderContentType},
		AllowCredentials: false,
	}))

	// routes auth
	routes.UserAuthRoutes(e, userController)               // user
	routes.AdminAuthRoutes(e, adminController)             // admin
	routes.DoctorAuthRoutes(e, doctorControllerManagement) // dokter

	// Group User
	userGroup := e.Group("/user", jwtMiddleware.HandlerUser)
	routes.UserProfil(userGroup, userProfilController, userFiturController, consultationController, artikelController)
	routes.UserChatbotRoutes(userGroup, chatbotController)

	// Group Admin
	adminGroup := e.Group("/admin", jwtMiddleware.HandlerAdmin)
	routes.AdminManagementRoutes(adminGroup, adminControllerManagement, artikelController, consultationController)

	// Group Doctor
	doctorGroup := e.Group("/doctor", jwtMiddleware.HandlerDoctor)
	routes.DoctorProfil(doctorGroup, doctorProfilController, artikelController, consultationController, userFiturController)

	// Mulai server
	log.Fatal(e.Start(":8000"))
}
