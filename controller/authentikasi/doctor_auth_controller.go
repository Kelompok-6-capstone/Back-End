package controller

import (
	"calmind/helper"
	"calmind/model"
	usecase "calmind/usecase/authentikasi"

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
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal registrasi pengguna: "+err.Error())
	}

	return helper.JSONSuccessResponse(ctx, "Berhasil registrasi pengguna, silakan verifikasi email Anda.")
}

func (c *DoctorAuthController) LoginDoctor(ctx echo.Context) error {
	var doctor model.Doctor
	if err := ctx.Bind(&doctor); err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Gagal mendapatkan data: "+err.Error())
	}

	if doctor.Email == "" || doctor.Password == "" {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Email dan password wajib diisi.")
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
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Input tidak valid: "+err.Error())
	}

	if otp.Email == "" || otp.Code == "" {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Email dan kode OTP wajib diisi.")
	}

	err := c.DoctorUsecase.VerifyOtp(otp.Email, otp.Code)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusUnauthorized, "Verifikasi OTP gagal: "+err.Error())
	}

	return helper.JSONSuccessResponse(ctx, "OTP berhasil diverifikasi.")
}

func (c *DoctorAuthController) LogoutDoctor(ctx echo.Context) error {
	return helper.JSONSuccessResponse(ctx, "Logout berhasil")
}

func (c *DoctorAuthController) ResendOtp(ctx echo.Context) error {
	var request struct {
		Email string `json:"email"`
	}

	if err := ctx.Bind(&request); err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Input tidak valid: "+err.Error())
	}

	if request.Email == "" {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Email wajib diisi.")
	}

	err := c.DoctorUsecase.ResendOtp(request.Email)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengirim ulang OTP: "+err.Error())
	}

	return helper.JSONSuccessResponse(ctx, "Kode OTP berhasil dikirim ulang.")
}
