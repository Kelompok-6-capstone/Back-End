package main

import (
	"calmind/config"
	"calmind/helper"
	"calmind/middlewares"
	repository_management "calmind/repository/admin_management"
	repository_artikel "calmind/repository/artikel"
	repository_authentikasi "calmind/repository/authentikasi"
	repository_chatbot_ai "calmind/repository/chatbot_ai"
	repository_customer_service "calmind/repository/customer_service"
	repository_konsultasi "calmind/repository/konsultasi"
	repository_profile "calmind/repository/profile"
	repository_user_fitur "calmind/repository/user_fitur"

	usecase_management "calmind/usecase/admin_management"
	usecase_artikel "calmind/usecase/artikel"
	usecase_authentikasi "calmind/usecase/authentikasi"
	usecase_chatbot_ai "calmind/usecase/chatbot_ai"
	usecase_customer_service "calmind/usecase/customer_service"
	usecase_konsultasi "calmind/usecase/konsultasi"
	usecase_profile "calmind/usecase/profile"
	usecase_user_fitur "calmind/usecase/user_fitur"

	controller_management "calmind/controller/admin_management"
	controller_artikel "calmind/controller/artikel"
	controller_authentikasi "calmind/controller/authentikasi"
	controller_chatbot_ai "calmind/controller/chatbot_ai"
	controller_customer_service "calmind/controller/customer_service"
	controller_konsultasi "calmind/controller/konsultasi"
	controller_notifikasi "calmind/controller/midtrans_notifikasi"
	controller_profile "calmind/controller/profile"
	controller_user_fitur "calmind/controller/user_fitur"

	"calmind/routes"
	"calmind/service"
	"log"
	"net/http"

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
	userRepo := repository_authentikasi.NewAuthRepository(DB)
	otpRepo := repository_authentikasi.NewOtpRepository(DB)
	userUsecase := usecase_authentikasi.NewAuthUsecase(userRepo, jwtService, otpRepo, otpService)
	userController := controller_authentikasi.NewAuthController(userUsecase)

	// Repositori, usecase, dan controller untuk Admin
	adminRepo := repository_authentikasi.NewAdminAuthRepository(DB)
	adminUsecase := usecase_authentikasi.NewAdminAuthUsecase(adminRepo, jwtService)
	adminController := controller_authentikasi.NewAdminAuthController(adminUsecase)

	// Repositori, usecase, dan controller untuk Admin management
	adminRepoManagement := repository_management.NewAdminManagementRepo(DB)
	adminUsecaseManagement := usecase_management.NewAdminManagementUsecase(adminRepoManagement)
	adminControllerManagement := controller_management.NewAdminManagementController(adminUsecaseManagement)

	// Repositori, usecase, dan controller untuk dokter
	doctorRepoManagement := repository_authentikasi.NewDoctorAuthRepository(DB)
	doctorUsecaseManagement := usecase_authentikasi.NewDoctorAuthUsecase(doctorRepoManagement, jwtService, otpRepo, otpService)
	doctorControllerManagement := controller_authentikasi.NewDoctorAuthController(doctorUsecaseManagement)

	//	Repositori, usecase, dan controller untuk Profil User
	userProfilRepo := repository_profile.NewUserProfilRepository(DB)
	userProfilUsecase := usecase_profile.NewUserProfileUseCase(userProfilRepo)
	userProfilController := controller_profile.NewProfilController(userProfilUsecase)

	//	Repositori, usecase, dan controller untuk Fitur User
	userFiturRepo := repository_user_fitur.NewUserFiturRepository(DB)
	userFiturUsecase := usecase_user_fitur.NewUserFiturUsecase(userFiturRepo)
	userFiturController := controller_user_fitur.NewUserFiturController(userFiturUsecase)

	//	Repositori, usecase, dan controller untuk Profil doctor
	doctorProfilRepo := repository_profile.NewDoctorProfilRepository(DB)
	doctorProfilUsecase := usecase_profile.NewDoctorProfileUseCase(doctorProfilRepo)
	doctorProfilController := controller_profile.NewDoctorProfileController(doctorProfilUsecase)

	//    Repositori, usecase, dan controller untuk Consultasi
	consultationRepo := repository_konsultasi.NewConsultationRepositoryImpl(DB)
	consultationUsecase := usecase_konsultasi.NewConsultationUsecaseImpl(consultationRepo)
	consultationController := controller_konsultasi.NewConsultationController(consultationUsecase)

	helper.StartExpiredConsultationJob(consultationUsecase.MarkExpiredConsultations)

	//    Repositori, usecase, dan controller untuk Consultasi
	artikelonRepo := repository_artikel.NewArtikelRepository(DB)
	artikelUsecase := usecase_artikel.NewArtikelUsecase(artikelonRepo)
	artikelController := controller_artikel.NewArtikelController(artikelUsecase)

	//    Repositori, usecase, dan controller untuk chatbot ai
	chatbotRepo := repository_chatbot_ai.NewChatLogRepository(DB)
	chatbotUsecase := usecase_chatbot_ai.NewChatbotUsecase(chatbotRepo)
	chatbotController := controller_chatbot_ai.NewChatbotController(chatbotUsecase)

	//    Repositori, usecase, dan controller untuk profil admin
	adminprofil := repository_profile.NewAdminProfileRepository(DB)
	adminusecase := usecase_profile.NewAdminProfileUseCase(adminprofil)
	admincontroller := controller_profile.NewAdminController(adminusecase)

	// customer service
	cs := repository_customer_service.NewCustServiceRepository(DB)
	csusecase := usecase_customer_service.NewCustServiceUsecase(cs)
	cscontroller := controller_customer_service.NewCustServiceController(csusecase)

	// notifikasi midtrans
	midtrans_notifikasi := controller_notifikasi.NewMidtransNotificationController(consultationUsecase)

	// Middleware
	jwtMiddleware := middlewares.NewJWTMiddleware(jwtSecret)

	e := echo.New()
	e.Static("/uploads", "uploads")
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"https://jovial-mooncake-23a3d0.netlify.app", "http://localhost:5173", "http://localhost:5174", "http://127.0.0.1:5500"},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowHeaders:     []string{echo.HeaderAuthorization, echo.HeaderContentType, echo.HeaderAccept},
		AllowCredentials: true,
	}))

	e.OPTIONS("/*", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	http.HandleFunc("/ws", helper.HandleConnections)

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
	routes.AdminManagementRoutes(adminGroup, adminControllerManagement, artikelController, consultationController, admincontroller)

	// Group Doctor
	doctorGroup := e.Group("/doctor", jwtMiddleware.HandlerDoctor)
	routes.DoctorProfil(doctorGroup, doctorProfilController, artikelController, consultationController, userFiturController)

	routes.UserCustServiceRoutes(e, cscontroller)

	routes.WebhookRoutes(e, midtrans_notifikasi)

	// Mulai server
	log.Fatal(e.Start(":8000"))
}
