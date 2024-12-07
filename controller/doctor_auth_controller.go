package controller

import (
	"calmind/helper"
	"calmind/model"
	"calmind/usecase"
	"net/http"

	"github.com/labstack/echo/v4"
)

type DoctorAuthController struct {
	DoctorUsecase usecase.DoctorUsecase
}

func NewDoctorAuthController(doctorUsecase usecase.DoctorUsecase) *DoctorAuthController {
	return &DoctorAuthController{
		DoctorUsecase: doctorUsecase,
	}
}

func (c *DoctorAuthController) RegisterDoctor(ctx echo.Context) error {
	var doctor model.Doctor
	if err := ctx.Bind(&doctor); err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Gagal mendapatkan data: "+err.Error())
	}

	err := c.DoctorUsecase.Register(&doctor)
	if err != nil {
		if err.Error() == "invalid title_id" {
			return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "title_id tidak valid. Silakan pilih title yang tersedia.")
		}
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal register dokter: "+err.Error())
	}

	return helper.JSONSuccessResponse(ctx, "Berhasil Register Dokter")
}

func (c *DoctorAuthController) LoginDoctor(ctx echo.Context) error {
	var doctor model.Doctor
	if err := ctx.Bind(&doctor); err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Gagal mendapatkan data: "+err.Error())
	}

	token, err := c.DoctorUsecase.Login(doctor.Email, doctor.Password)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusUnauthorized, "Login gagal: "+err.Error())
	}

	// Kembalikan token dalam respons
	return helper.JSONSuccessResponse(ctx, map[string]string{
		"token": token,
	})
}

func (c *DoctorAuthController) VerifyOtp(ctx echo.Context) error {
	var otp model.Otp
	if err := ctx.Bind(&otp); err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Invalid input: "+err.Error())
	}

	err := c.DoctorUsecase.VerifyOtp(otp.Email, otp.Code)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusUnauthorized, "OTP verification failed: "+err.Error())
	}

	return helper.JSONSuccessResponse(ctx, "OTP verified successfully")
}

func (c *DoctorAuthController) LogoutDoctor(ctx echo.Context) error {
	return helper.JSONSuccessResponse(ctx, "Logout berhasil")
}

func (c *DoctorAuthController) ResendOtp(ctx echo.Context) error {
	var request struct {
		Email string `json:"email"`
	}

	// Bind data input
	if err := ctx.Bind(&request); err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Invalid input: "+err.Error())
	}

	// Validasi input email
	if request.Email == "" {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Email is required")
	}

	// Resend OTP
	err := c.DoctorUsecase.ResendOtp(request.Email)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Failed to resend OTP: "+err.Error())
	}

	return helper.JSONSuccessResponse(ctx, "OTP sent successfully")
}
