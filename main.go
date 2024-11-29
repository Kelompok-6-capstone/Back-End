package main

import (
	"calmind/config"
	"calmind/controller"
	"calmind/middleware"
	"calmind/repository"
	"calmind/routes"
	"calmind/service"
	"calmind/usecase"
	"log"

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

	// Repositori, usecase, dan controller untuk User
	userRepo := repository.NewAuthRepository(DB)
	userUsecase := usecase.NewAuthUsecase(userRepo, jwtService)
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
	doctorUsecaseManagement := usecase.NewDoctorAuthUsecase(doctorRepoManagement, jwtService)
	doctorControllerManagement := controller.NewDoctorAuthController(doctorUsecaseManagement)

	//	Repositori, usecase, dan controller untuk Profil User
	userProfilRepo := repository.NewUserProfilRepository(DB)
	userProfilUsecase := usecase.NewUserProfilUsecaseImpl(userProfilRepo)
	userProfilController := controller.NewProfilController(userProfilUsecase)

	// Middleware
	jwtMiddleware := middleware.NewJWTMiddleware(jwtSecret)

	// Echo instance
	e := echo.New()

	// Routes untuk User
	userGroup := e.Group("/user")

	routes.UserAuthRoutes(userGroup, userController, userProfilController)
	routes.AdminAuthRoutes(e, adminController)

	// Routes untuk Admin
	adminGroup := e.Group("")
	adminGroup.Use(jwtMiddleware.HandlerAdmin)
	routes.AdminManagementRoutes(adminGroup, adminControllerManagement)

	// Routes untuk Dokter
	routes.DoctorAuthRoutes(e, doctorControllerManagement)

	// Mulai server
	log.Fatal(e.Start(":8000"))
}
