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
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "gagal mendapatkan data: "+err.Error())
	}

	err := c.DoctorUsecase.Register(&doctor)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "gagal register dokter: "+err.Error())
	}

	return helper.JSONSuccessResponse(ctx, "Berhasil Register Dokter")
}

func (c *DoctorAuthController) LoginDoctor(ctx echo.Context) error {
	var doctor model.Doctor
	if err := ctx.Bind(&doctor); err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "gagal mendapatkan data: "+err.Error())
	}

	token, err := c.DoctorUsecase.Login(doctor.Email, doctor.Password)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusUnauthorized, "Login gagal: "+err.Error())
	}

	cookie := &http.Cookie{
		Name:     "token_doctor",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		MaxAge:   72 * 60 * 60,
	}

	ctx.SetCookie(cookie)

	return helper.JSONSuccessResponse(ctx, map[string]string{
		"Token": token,
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
	cookie := &http.Cookie{
		Name:     "token_doctor",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		MaxAge:   -1,
	}

	ctx.SetCookie(cookie)

	return helper.JSONSuccessResponse(ctx, "Berhasil Logout Dokter")
}
