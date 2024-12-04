package controller

import (
	"calmind/helper"
	"calmind/model"
	"calmind/usecase"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AdminAuthController struct {
	AdminUsecase usecase.AdminAuthUsecase
}

func NewAdminAuthController(usecase usecase.AdminAuthUsecase) *AdminAuthController {
	return &AdminAuthController{AdminUsecase: usecase}
}

func (c *AdminAuthController) LoginAdmin(ctx echo.Context) error {
	var user model.User
	if err := ctx.Bind(&user); err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "gagal mendapatkan data: "+err.Error())
	}

	token, err := c.AdminUsecase.LoginAdmin(user.Email, user.Password)
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

func (c *AdminAuthController) LogoutAdmin(ctx echo.Context) error {
	cookie := &http.Cookie{
		Name:     "token_admin",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		MaxAge:   -1,
	}

	ctx.SetCookie(cookie)

	return helper.JSONSuccessResponse(ctx, "Berhasil Logout")
}
