package controller

import (
	"calmind/helper"
	"calmind/model"
	"calmind/service"
	usecase "calmind/usecase/profile"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v4"
)

type ProfilController struct {
	ProfilUsecase usecase.UserProfileUseCase
}

func NewProfilController(ProfilUsecase usecase.UserProfileUseCase) *ProfilController {
	return &ProfilController{ProfilUsecase: ProfilUsecase}
}

func (c *ProfilController) GetProfile(ctx echo.Context) error {
	claims, _ := ctx.Get("user").(*service.JwtCustomClaims)

	user, err := c.ProfilUsecase.GetUserProfile(claims.UserID)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil profil: "+err.Error())
	}

	type UserProfileResponse struct {
		ID            int    `json:"id"`
		Avatar        string `json:"avatar"`
		Username      string `json:"username"`
		Email         string `json:"email"`
		NoHp          string `json:"no_hp"`
		Alamat        string `json:"alamat"`
		Tgl_lahir     string `json:"tgl_lahir"`
		Jenis_kelamin string `json:"jenis_kelamin"`
		Pekerjaan     string `json:"pekerjaan"`
	}

	userProfile := UserProfileResponse{
		ID:            user.ID,
		Avatar:        user.Avatar,
		Username:      user.Username,
		Email:         user.Email,
		NoHp:          user.NoHp,
		Alamat:        user.Alamat,
		Tgl_lahir:     user.TglLahir,
		Jenis_kelamin: user.JenisKelamin,
		Pekerjaan:     user.Pekerjaan,
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

	return helper.JSONSuccessResponse(ctx, "berhasil update profil")
}
func (c *ProfilController) UploadAvatar(ctx echo.Context) error {
	claims, ok := ctx.Get("user").(*service.JwtCustomClaims)
	if !ok {
		return helper.JSONErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized")
	}
	userID := claims.UserID

	// Ambil file dari form
	file, err := ctx.FormFile("avatar")
	if err != nil || file.Size > 5*1024*1024 {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "File tidak valid atau ukuran melebihi 5 MB")
	}

	ext := filepath.Ext(file.Filename)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Format file tidak didukung")
	}

	// Buka file untuk upload
	src, err := file.Open()
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal membuka file")
	}
	defer src.Close()

	// Upload ke Cloudinary
	fileName := fmt.Sprintf("user_%d_avatar%s", userID, ext)
	avatarURL, publicID, err := helper.UploadFileToCloudinary(src, fileName)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal upload ke Cloudinary")
	}

	// Update avatar di database
	user := model.User{Avatar: avatarURL}
	_, err = c.ProfilUsecase.UpdateUserProfile(userID, &user)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengupdate avatar: "+err.Error())
	}

	return helper.JSONSuccessResponse(ctx, map[string]string{
		"message":   "Avatar berhasil diupload",
		"avatarUrl": avatarURL,
		"publicID":  publicID,
	})
}

func (c *ProfilController) DeleteAvatar(ctx echo.Context) error {
	claims, ok := ctx.Get("user").(*service.JwtCustomClaims)
	if !ok {
		return helper.JSONErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized")
	}
	userID := claims.UserID

	// Ambil profil user
	user, err := c.ProfilUsecase.GetUserProfile(userID)
	if err != nil || user.Avatar == "" {
		return helper.JSONErrorResponse(ctx, http.StatusNotFound, "Avatar tidak ditemukan")
	}

	// Ekstrak public_id dari URL
	parts := strings.Split(user.Avatar, "/")
	publicID := strings.TrimSuffix(parts[len(parts)-1], filepath.Ext(user.Avatar))

	// Hapus dari Cloudinary
	if err := helper.DeleteFileFromCloudinary(publicID); err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal menghapus avatar dari Cloudinary")
	}

	// Update avatar di database menjadi kosong
	user.Avatar = ""
	_, err = c.ProfilUsecase.UpdateUserProfile(userID, user)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengupdate database: "+err.Error())
	}

	return helper.JSONSuccessResponse(ctx, "Avatar berhasil dihapus")
}
