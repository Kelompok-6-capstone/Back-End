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
		ID           int    `json:"id"`
		Avatar       string `json:"avatar"`
		Username     string `json:"username"`
		Email        string `json:"email"`
		NoHp         string `json:"no_hp"`
		Alamat       string `json:"alamat"`
		Tgl_lahir    string `json:"tgl_lahir"`
		JenisKelamin string `json:"jenis_kelamin"`
	}

	userProfile := UserProfileResponse{
		ID:           user.ID,
		Avatar:       user.Avatar,
		Username:     user.Username,
		Email:        user.Email,
		NoHp:         user.NoHp,
		Alamat:       user.Alamat,
		Tgl_lahir:    user.Tgl_lahir,
		JenisKelamin: user.JenisKelamin,
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

	return helper.JSONSuccessResponse(ctx, "Berhasil update profil")
}
