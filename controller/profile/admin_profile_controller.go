package controller

import (
	"calmind/helper"
	"calmind/service"
	usecase "calmind/usecase/profile"
	"fmt"
	"io"

	"net/http"
	"os"
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
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Gagal mendapatkan file: "+err.Error())
	}

	// Validasi ukuran file (maksimal 5 MB)
	if file.Size > 5*1024*1024 {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Ukuran file maksimal 5 MB")
	}

	// Validasi ekstensi file
	ext := filepath.Ext(file.Filename)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Hanya file dengan format .jpg, .jpeg, atau .png yang diperbolehkan")
	}

	// Simpan file di direktori uploads
	uploadDir := "uploads/admin_avatars"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		err = os.MkdirAll(uploadDir, os.ModePerm)
		if err != nil {
			return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal membuat direktori upload")
		}
	}

	// Buat path untuk file avatar
	filePath := fmt.Sprintf("%s/admin_%d_%s", uploadDir, adminID, file.Filename)
	src, err := file.Open()
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal membuka file")
	}
	defer src.Close()

	dst, err := os.Create(filePath)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal menyimpan file")
	}
	defer dst.Close()

	// Salin isi file dari sumber ke tujuan
	if _, err := io.Copy(dst, src); err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal menyimpan file")
	}

	// URL file avatar yang disimpan
	avatarURL := fmt.Sprintf("https://%s/%s", ctx.Request().Host, filePath)

	// Update URL avatar di database
	err = c.AdminUsecase.UploadAdminAvatar(adminID, avatarURL, "")
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengupdate avatar admin: "+err.Error())
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

	// Ambil data admin untuk mendapatkan URL avatar
	admin, err := c.AdminUsecase.GetAdminProfile(adminID)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil data admin: "+err.Error())
	}

	// Hapus file avatar
	if admin.Avatar != "" {
		filePath := "." + admin.Avatar // Tambahkan "." untuk path relatif
		if _, err := os.Stat(filePath); err == nil {
			if err := os.Remove(filePath); err != nil {
				return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal menghapus file avatar: "+err.Error())
			}
		}
	}

	// Update avatar menjadi kosong di database
	admin.Avatar = ""
	err = c.AdminUsecase.DeleteAdminAvatar(adminID)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengupdate avatar admin: "+err.Error())
	}

	return helper.JSONSuccessResponse(ctx, "Avatar berhasil dihapus")
}
