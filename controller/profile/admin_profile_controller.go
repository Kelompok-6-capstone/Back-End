package controller

import (
	"calmind/helper"
	"calmind/service"
	usecase "calmind/usecase/profile"
	"fmt"
	"strings"

	"net/http"
	"path/filepath"

	"github.com/labstack/echo/v4"
)

type AdminController struct {
	AdminUsecase usecase.AdminProfileUseCase
}

func NewAdminController(adminUsecase usecase.AdminProfileUseCase) *AdminController {
	return &AdminController{AdminUsecase: adminUsecase}
}

// Get Admin Profile
func (c *AdminController) GetAdminProfile(ctx echo.Context) error {
	claims, _ := ctx.Get("admin").(*service.JwtCustomClaims)

	admin, err := c.AdminUsecase.GetAdminProfile(claims.UserID)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil profil admin: "+err.Error())
	}

	type AdminProfileResponse struct {
		ID       int    `json:"id"`
		Avatar   string `json:"avatar"`
		Username string `json:"username"`
	}

	adminProfile := AdminProfileResponse{
		ID:       admin.ID,
		Avatar:   admin.Avatar,
		Username: admin.Username,
	}

	return helper.JSONSuccessResponse(ctx, adminProfile)
}
func (c *AdminController) UploadAdminAvatar(ctx echo.Context) error {
	claims, ok := ctx.Get("admin").(*service.JwtCustomClaims)
	if !ok {
		return helper.JSONErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized")
	}
	adminID := claims.UserID

	// Ambil file dari form
	file, err := ctx.FormFile("avatar")
	if err != nil || file.Size > 5*1024*1024 {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "File tidak valid atau ukuran melebihi 5 MB")
	}

	ext := filepath.Ext(file.Filename)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Format file tidak didukung")
	}

	src, err := file.Open()
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal membuka file")
	}
	defer src.Close()

	// Upload ke Cloudinary
	fileName := fmt.Sprintf("admin_%d_avatar%s", adminID, ext)
	avatarURL, publicID, err := helper.UploadFileToCloudinary(src, fileName)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal upload ke Cloudinary")
	}

	// Update database
	if err := c.AdminUsecase.UploadAdminAvatar(adminID, avatarURL, publicID); err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengupdate database: "+err.Error())
	}

	return helper.JSONSuccessResponse(ctx, map[string]string{
		"message":   "Avatar berhasil diupload",
		"avatarUrl": avatarURL,
	})
}

func (c *AdminController) DeleteAdminAvatar(ctx echo.Context) error {
	claims, ok := ctx.Get("admin").(*service.JwtCustomClaims)
	if !ok {
		return helper.JSONErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized")
	}
	adminID := claims.UserID

	// Ambil data admin
	admin, err := c.AdminUsecase.GetAdminProfile(adminID)
	if err != nil || admin.Avatar == "" {
		return helper.JSONErrorResponse(ctx, http.StatusNotFound, "Avatar tidak ditemukan")
	}

	// Ekstrak public_id dari URL
	parts := strings.Split(admin.Avatar, "/")
	publicID := strings.TrimSuffix(parts[len(parts)-1], filepath.Ext(admin.Avatar))

	// Hapus dari Cloudinary
	if err := helper.DeleteFileFromCloudinary(publicID); err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal menghapus file dari Cloudinary")
	}

	// Update database
	if err := c.AdminUsecase.DeleteAdminAvatar(adminID); err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengupdate database: "+err.Error())
	}

	return helper.JSONSuccessResponse(ctx, "Avatar berhasil dihapus")
}
