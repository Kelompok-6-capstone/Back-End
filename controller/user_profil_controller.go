package controller

import (
	"calmind/helper"
	"calmind/model"
	"calmind/service"
	"calmind/usecase"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

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
		Tgl_lahir:     user.Tgl_lahir,
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

// Upload Avatar
func (c *ProfilController) UploadAvatar(ctx echo.Context) error {
	claims, ok := ctx.Get("user").(*service.JwtCustomClaims)
	if !ok {
		return helper.JSONErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized")
	}
	userID := claims.UserID

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
	uploadDir := "uploads"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		err = os.MkdirAll(uploadDir, os.ModePerm)
		if err != nil {
			return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal membuat direktori upload")
		}
	}

	filePath := fmt.Sprintf("%s/%d_%s", uploadDir, userID, file.Filename)
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

	if _, err := helper.CopyFile(src, dst); err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal menyimpan file")
	}

	// Update URL avatar di database
	avatarURL := fmt.Sprintf("http://%s/uploads/%d_%s", ctx.Request().Host, userID, file.Filename)
	user := model.User{
		Avatar: avatarURL,
	}
	_, err = c.ProfilUsecase.UpdateUserProfile(userID, &user)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengupdate avatar: "+err.Error())
	}

	return helper.JSONSuccessResponse(ctx, map[string]string{
		"message":   "Avatar berhasil diupload",
		"avatarUrl": avatarURL,
	})
}

// Delete Avatar
func (c *ProfilController) DeleteAvatar(ctx echo.Context) error {
	claims, ok := ctx.Get("user").(*service.JwtCustomClaims)
	if !ok {
		return helper.JSONErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized")
	}
	userID := claims.UserID

	// Ambil data user untuk mendapatkan avatar URL
	user, err := c.ProfilUsecase.GetUserProfile(userID)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil profil user: "+err.Error())
	}

	// Hapus file avatar
	if user.Avatar != "" {
		filePath := "." + user.Avatar // Tambahkan "." untuk path relatif
		if _, err := os.Stat(filePath); err == nil {
			if err := os.Remove(filePath); err != nil {
				return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal menghapus file avatar: "+err.Error())
			}
		}
	}

	// Update avatar menjadi kosong di database
	user.Avatar = ""
	_, err = c.ProfilUsecase.UpdateUserProfile(userID, user)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengupdate avatar di database: "+err.Error())
	}

	return helper.JSONSuccessResponse(ctx, "Avatar berhasil dihapus")
}
