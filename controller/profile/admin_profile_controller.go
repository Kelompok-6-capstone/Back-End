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
	claims, _ := ctx.Get("admin").(*service.JwtCustomClaims)

	// Ambil file dari form
	file, err := ctx.FormFile("avatar")
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Gagal mendapatkan file: "+err.Error())
	}

	// Validasi ukuran file (maksimal 10 MB)
	if file.Size > 10*1024*1024 {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Ukuran file maksimal 10 MB")
	}

	// Validasi ekstensi file
	ext := filepath.Ext(file.Filename)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Hanya file dengan format .jpg, .jpeg, atau .png yang diperbolehkan")
	}

	src, err := file.Open()
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal membuka file")
	}
	defer src.Close()

	// Path direktori penyimpanan avatar
	uploadDir := "/app/uploads/admin_avatars"
	err = os.MkdirAll(uploadDir, 0777) // Membuat direktori jika belum ada
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal membuat direktori upload: "+err.Error())
	}

	// Path file avatar
	filePath := fmt.Sprintf("%s/%d%s", uploadDir, claims.UserID, ext)

	// Simpan file avatar
	dst, err := os.Create(filePath)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal membuat file avatar: "+err.Error())
	}
	defer dst.Close()

	// Salin konten file
	if _, err = io.Copy(dst, src); err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal menyimpan file avatar: "+err.Error())
	}

	// URL file avatar yang disimpan
	imageURL := fmt.Sprintf("/uploads/admin_avatars/%d%s", claims.UserID, ext)

	// Update URL avatar di database
	err = c.AdminUsecase.UploadAdminAvatar(claims.UserID, imageURL, "")
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengupdate avatar admin: "+err.Error())
	}

	return helper.JSONSuccessResponse(ctx, map[string]string{
		"message":   "Avatar berhasil diupload",
		"avatarUrl": imageURL,
	})
}

func (c *AdminController) DeleteAdminAvatar(ctx echo.Context) error {
	claims, _ := ctx.Get("admin").(*service.JwtCustomClaims)

	// Ambil data admin dari database untuk mendapatkan URL avatar
	admin, err := c.AdminUsecase.GetAdminProfile(claims.UserID)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil data admin: "+err.Error())
	}

	// Path file avatar
	avatarPath := fmt.Sprintf("/app%s", admin.Avatar)

	// Hapus file avatar jika ada
	if admin.Avatar != "" {
		if _, err := os.Stat(avatarPath); err == nil {
			err := os.Remove(avatarPath)
			if err != nil {
				return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal menghapus avatar: "+err.Error())
			}
		}
	}

	// Update avatar menjadi kosong di database
	err = c.AdminUsecase.DeleteAdminAvatar(claims.UserID)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengupdate avatar admin: "+err.Error())
	}

	return helper.JSONSuccessResponse(ctx, "Avatar berhasil dihapus")
}
