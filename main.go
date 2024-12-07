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
	// Memuat konfigurasi .env
	if err := godotenv.Load(); err != nil {
		log.Fatal("Gagal memuat file .env")
	}

	// Inisialisasi database
	DB, err := config.InitDB()
	if err != nil {
		log.Fatalf("Gagal menginisialisasi database: %v", err)
	}
	log.Println("Database berhasil diinisialisasi.")

	// Konfigurasi JWT dan OTP Service
	jwtSecret := config.NewJWTConfig()
	jwtService := service.NewJWTService(jwtSecret)
	otpService := service.NewOtpService()
	log.Println("JWT dan OTP Service berhasil dikonfigurasi.")

	// Inisialisasi Repositori
	userRepo := repository.NewAuthRepository(DB)
	otpRepo := repository.NewOtpRepository(DB)
	adminRepo := repository.NewAdminAuthRepository(DB)
	adminManagementRepo := repository.NewAdminManagementRepo(DB)
	doctorAuthRepo := repository.NewDoctorAuthRepository(DB)
	userProfilRepo := repository.NewUserProfilRepository(DB)
	userFiturRepo := repository.NewUserFiturRepository(DB)
	doctorProfilRepo := repository.NewDoctorProfilRepository(DB)
	consultationRepo := repository.NewConsultationRepository(DB)
	artikelRepo := repository.NewArtikelRepository(DB)
	paymentRepo := repository.NewPaymentRepository(DB)

	log.Println("Repositori berhasil diinisialisasi.")

	// Inisialisasi Usecase
	userUsecase := usecase.NewAuthUsecase(userRepo, jwtService, otpRepo, otpService)
	adminUsecase := usecase.NewAdminAuthUsecase(adminRepo, jwtService)
	adminManagementUsecase := usecase.NewAdminManagementUsecase(adminManagementRepo)
	doctorAuthUsecase := usecase.NewDoctorAuthUsecase(doctorAuthRepo, jwtService, otpRepo, otpService)
	userProfilUsecase := usecase.NewUserProfileUseCase(userProfilRepo)
	userFiturUsecase := usecase.NewUserFiturUsecase(userFiturRepo)
	doctorProfilUsecase := usecase.NewDoctorProfileUseCase(doctorProfilRepo)
	consultationUsecase := usecase.NewConsultationUsecase(consultationRepo)
	artikelUsecase := usecase.NewArtikelUsecase(artikelRepo)
	paymentUsecase := usecase.NewPaymentUsecase(paymentRepo)

	log.Println("Usecase berhasil diinisialisasi.")

	// Inisialisasi Controller
	userController := controller.NewAuthController(userUsecase)
	adminController := controller.NewAdminAuthController(adminUsecase)
	adminManagementController := controller.NewAdminManagementController(adminManagementUsecase)
	doctorAuthController := controller.NewDoctorAuthController(doctorAuthUsecase)
	userProfilController := controller.NewProfilController(userProfilUsecase)
	userFiturController := controller.NewUserFiturController(userFiturUsecase)
	doctorProfilController := controller.NewDoctorProfileController(doctorProfilUsecase)
	consultationController := controller.NewConsultationController(consultationUsecase, paymentUsecase)
	artikelController := controller.NewArtikelController(artikelUsecase)

	log.Println("Controller berhasil diinisialisasi.")

	// Middleware JWT
	jwtMiddleware := middlewares.NewJWTMiddleware(jwtSecret)

	// Inisialisasi Echo
	e := echo.New()
	e.Static("/uploads", "uploads")
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:5174", "http://localhost:5173", "http://127.0.0.1:5500", "https://jovial-mooncake-23a3d0.netlify.app"},
		AllowMethods:     []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS},
		AllowHeaders:     []string{echo.HeaderAuthorization, echo.HeaderContentType},
		AllowCredentials: true,
	}))

	log.Println("Echo server berhasil dikonfigurasi.")

	// Konfigurasi Routes
	routes.UserAuthRoutes(e, userController)
	routes.AdminAuthRoutes(e, adminController)
	routes.DoctorAuthRoutes(e, doctorAuthController)

	userGroup := e.Group("/user", jwtMiddleware.HandlerUser)
	routes.UserProfil(userGroup, userProfilController, userFiturController, consultationController, artikelController)

	adminGroup := e.Group("/admin", jwtMiddleware.HandlerAdmin)
	routes.AdminManagementRoutes(adminGroup, adminManagementController, consultationController, artikelController)

	doctorGroup := e.Group("/doctor", jwtMiddleware.HandlerDoctor)
	routes.DoctorProfil(doctorGroup, doctorProfilController, artikelController, consultationController)

	log.Println("Routes berhasil dikonfigurasi.")

	// Menjalankan server
	log.Fatal(e.Start(":8000"))
}
