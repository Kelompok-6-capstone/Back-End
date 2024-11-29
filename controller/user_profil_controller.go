package controller

import (
	"calmind/helper"
	"calmind/model"
	"calmind/usecase"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ProfilController struct {
	ProfilUsecase usecase.UserProfilUsecase
}

func NewProfilController(ProfilUsecase usecase.UserProfilUsecase) *ProfilController {
	return &ProfilController{ProfilUsecase: ProfilUsecase}
}
func (c *ProfilController) GetProfile(ctx echo.Context) error {
	userClaims := ctx.Get("user").(*model.User)

	// Mendapatkan data profil user berdasarkan ID
	user, err := c.ProfilUsecase.GetUserProfile(userClaims.ID)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil profil: "+err.Error())
	}

	return helper.JSONSuccessResponse(ctx, user)
}
