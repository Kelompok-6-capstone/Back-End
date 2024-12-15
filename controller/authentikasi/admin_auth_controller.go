package controller

import (
	"calmind/helper"
	"calmind/model"
	usecase "calmind/usecase/authentikasi"

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
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Gagal mendapatkan data: "+err.Error())
	}

	token, err := c.AdminUsecase.LoginAdmin(user.Email, user.Password)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusUnauthorized, "Login gagal: "+err.Error())
	}

	// Kembalikan token dalam respons
	return helper.JSONSuccessResponse(ctx, map[string]string{
		"token": token,
	})
}

func (c *AdminAuthController) LogoutAdmin(ctx echo.Context) error {
	return helper.JSONSuccessResponse(ctx, "Logout berhasil")
}
