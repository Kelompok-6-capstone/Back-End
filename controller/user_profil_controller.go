package controller

import (
	"calmind/helper"
	"calmind/model"
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

	type UserProfileResponse struct {
		ID       int    `json:"id"`
		Username string `json:"username"`
		NoHp     string `json:"no_hp"`
		Email    string `json:"email"`
		Avatar   string `json:"avatar"`
		FullName string `json:"full_name"`
		Bio      string `json:"bio"`
	}

	userProfile := UserProfileResponse{
		ID:       user.ID,
		Username: user.Username,
		NoHp:     user.NoHp,
		Email:    user.Email,
		Avatar:   user.Avatar,
		Bio:      user.Bio,
	}

	return helper.JSONSuccessResponse(ctx, userProfile)
}

func (c *ProfilController) UpdateProfile(ctx echo.Context) error {
	claims, _ := ctx.Get("user").(*service.JwtCustomClaims)
	var user model.User

	if err := ctx.Bind(&user); err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Gagal mendapatkan data: "+err.Error())
	}

	user.Email = ""
	user.Password = ""

	_, err := c.ProfilUsecase.UpdateUserProfile(claims.UserID, &user)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengupdate profil: "+err.Error())
	}

	return helper.JSONSuccessResponse(ctx, "Berhasul update profil")
}
