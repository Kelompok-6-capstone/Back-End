package controller

import (
	"calmind/helper"
	"calmind/model"
	usecase "calmind/usecase/authentikasi"

	"net/http"

	"github.com/labstack/echo/v4"
)

type AuthController struct {
	AuthUsecase usecase.UserUsecase
}

func NewAuthController(authUsecase usecase.UserUsecase) *AuthController {
	return &AuthController{
		AuthUsecase: authUsecase,
	}
}

func (c *AuthController) RegisterUser(ctx echo.Context) error {
	var user model.User

	if err := ctx.Bind(&user); err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Gagal mendapatkan data: "+err.Error())
	}

	err := c.AuthUsecase.Register(&user)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal registrasi pengguna: "+err.Error())
	}

	return helper.JSONSuccessResponse(ctx, "Berhasil registrasi pengguna, silakan verifikasi email Anda.")
}

func (c *AuthController) LoginUser(ctx echo.Context) error {
	var request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := ctx.Bind(&request); err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Gagal mendapatkan data: "+err.Error())
	}

	if request.Email == "" || request.Password == "" {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Email dan password wajib diisi.")
	}

	token, err := c.AuthUsecase.Login(request.Email, request.Password)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusUnauthorized, "Login gagal: "+err.Error())
	}

	return helper.JSONSuccessResponse(ctx, map[string]string{
		"token": token,
	})
}

func (c *AuthController) VerifyOtp(ctx echo.Context) error {
	var otp struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	}

	if err := ctx.Bind(&otp); err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Input tidak valid: "+err.Error())
	}

	if otp.Email == "" || otp.Code == "" {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Email dan kode OTP wajib diisi.")
	}

	err := c.AuthUsecase.VerifyOtp(otp.Email, otp.Code)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusUnauthorized, "Verifikasi OTP gagal: "+err.Error())
	}

	return helper.JSONSuccessResponse(ctx, "OTP berhasil diverifikasi.")
}

func (c *AuthController) LogoutUser(ctx echo.Context) error {
	return helper.JSONSuccessResponse(ctx, "Logout berhasil.")
}

func (c *AuthController) ResendOtp(ctx echo.Context) error {
	var request struct {
		Email string `json:"email"`
	}

	if err := ctx.Bind(&request); err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Input tidak valid: "+err.Error())
	}

	if request.Email == "" {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Email wajib diisi.")
	}

	err := c.AuthUsecase.ResendOtp(request.Email)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengirim ulang OTP: "+err.Error())
	}

	return helper.JSONSuccessResponse(ctx, "Kode OTP berhasil dikirim ulang.")
}
