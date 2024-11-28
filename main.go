package main

import (
	"calmind/config"
	"calmind/controller"
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

	// jwt
	jwtsecret := config.NewJWTConfig()
	jwtservice := service.NewJWTService(jwtsecret)

	// state user
	userAuthRepo := repository.NewAuthRepository(DB)
	userAuthUsecase := usecase.NewAuthUsecase(userAuthRepo, jwtservice)
	userAuthcontroller := controller.NewAuthController(userAuthUsecase)

	// state admin
	adminAuthRepo := repository.NewAdminAuthRepository(DB)
	adminAuthUsecase := usecase.NewAdminAuthUsecase(adminAuthRepo, jwtservice)
	adminAuthcontroller := controller.NewAdminAuthController(adminAuthUsecase)

	// state doctor
	docterAuthRepo := repository.NewDoctorAuthRepository(DB)
	doctorAuthUsecase := usecase.NewDoctorAuthUsecase(docterAuthRepo, jwtservice)
	doctorAuthcontroller := controller.NewDoctorAuthController(doctorAuthUsecase)

	e := echo.New()

	// middleware
	// middlewareJwt := middleware.NewJWTMiddleware(jwtsecret)

	// middleware user
	eAuthUser := e.Group("/user")
	// eAuthUser.Use(middlewareJwt.HandlerUser)

	// middleware admin
	eAuthAdmin := e.Group("/admin")
	// eAuthAdmin.Use(middlewareJwt.HandlerAdmin)

	// middleware admin
	eAuthDoctor := e.Group("/doctor")
	// eAuthAdmin.Use(middlewareJwt.HandlerAdmin)

	// routes
	routes.UserAuthRoutes(eAuthUser, userAuthcontroller)
	routes.AdminAuthRoutes(eAuthAdmin, adminAuthcontroller)
	routes.DoctorAuthRoutes(eAuthDoctor, doctorAuthcontroller)

	e.Start(":8000")
}
