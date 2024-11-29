package controller

import (
	"calmind/helper"
	"calmind/service"
	"calmind/usecase"
	"log"
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
	claims, _ := ctx.Get("user").(*service.JwtCustomClaims)

	user, err := c.ProfilUsecase.GetUserProfile(claims.UserID)
	if err != nil {
		log.Println("Failed to fetch user profile:", err)
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil profil: "+err.Error())
	}

	user.Password = ""
	user.Role = ""

	return helper.JSONSuccessResponse(ctx, user)
}
