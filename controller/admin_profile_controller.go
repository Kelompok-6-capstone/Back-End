package controller

import (
	"calmind/helper"
	"calmind/service"
	"calmind/usecase"
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

// Upload Admin Avatar
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

	// Upload gambar ke ImgBB
	imageURL, deleteURL, err := helper.UploadToImgBB(os.Getenv("API_KEY_IMBB"), file.Filename, src)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengunggah avatar ke ImgBB: "+err.Error())
	}

	// Update URL avatar di database
	err = c.AdminUsecase.UploadAdminAvatar(claims.UserID, imageURL, deleteURL)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengupdate avatar admin: "+err.Error())
	}

	return helper.JSONSuccessResponse(ctx, map[string]string{
		"message":   "Avatar berhasil diupload",
		"avatarUrl": imageURL,
	})
}

// Delete Admin Avatar
func (c *AdminController) DeleteAdminAvatar(ctx echo.Context) error {
	claims, _ := ctx.Get("admin").(*service.JwtCustomClaims)

	err := c.AdminUsecase.DeleteAdminAvatar(claims.UserID)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal menghapus avatar: "+err.Error())
	}

	return helper.JSONSuccessResponse(ctx, "Avatar berhasil dihapus")
}
