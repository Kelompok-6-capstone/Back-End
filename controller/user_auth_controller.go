package controller

import (
	"calmind/helper"
	"calmind/model"
	"calmind/usecase"
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
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "gagal mendapatkan data: "+err.Error())
	}

	err := c.AuthUsecase.Register(&user)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "gagal register user: "+err.Error())
	}

	return helper.JSONSuccessResponse(ctx, "Berhasil Register User")
}

func (c *AuthController) LoginUser(ctx echo.Context) error {
	var user model.User
	if err := ctx.Bind(&user); err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "gagal mendapatkan data: "+err.Error())
	}

	token, err := c.AuthUsecase.Login(user.Email, user.Password)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusUnauthorized, "Login gagal: "+err.Error())
	}

	cookie := &http.Cookie{
		Name:     "token_doctor",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,                  // Ubah ke true jika menggunakan HTTPS di frontend
		MaxAge:   72 * 60 * 60,          // Masa aktif cookie (72 jam)
		SameSite: http.SameSiteNoneMode, // None untuk mendukung lintas domain
	}
	ctx.SetCookie(cookie)

	return helper.JSONSuccessResponse(ctx, map[string]string{
		"Token": token,
	})
}

func (c *AuthController) VerifyOtp(ctx echo.Context) error {
	var otp model.Otp
	if err := ctx.Bind(&otp); err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Invalid input: "+err.Error())
	}
	err := c.AuthUsecase.VerifyOtp(otp.Email, otp.Code)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusUnauthorized, "OTP verification failed: "+err.Error())
	}
	return helper.JSONSuccessResponse(ctx, "OTP verified successfully")
}

func (c *AuthController) LogoutUser(ctx echo.Context) error {
	cookie := &http.Cookie{
		Name:     "token_user",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		MaxAge:   -1,
	}

	ctx.SetCookie(cookie)

	return helper.JSONSuccessResponse(ctx, "Berhasil Logout")
}
