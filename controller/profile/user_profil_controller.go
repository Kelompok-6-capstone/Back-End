package controller

import (
	"calmind/helper"
	"calmind/model"
	"calmind/service"
	usecase "calmind/usecase/profile"
	"fmt"
	"io"
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

	// Buka file
	src, err := file.Open()
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal membuka file")
	}
	defer src.Close()

	// Path direktori penyimpanan avatar
	uploadDir := "/app/uploads/avatars"
	err = os.MkdirAll(uploadDir, 0777) // Pastikan direktori ada
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal membuat direktori upload: "+err.Error())
	}

	// Path file avatar
	filePath := fmt.Sprintf("%s/%d%s", uploadDir, userID, ext)

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

	// Update URL avatar di database
	user := model.User{
		Avatar: fmt.Sprintf("/uploads/avatars/%d%s", userID, ext), // URL relatif ke file
	}
	_, err = c.ProfilUsecase.UpdateUserProfile(userID, &user)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengupdate avatar di database: "+err.Error())
	}

	return helper.JSONSuccessResponse(ctx, map[string]string{
		"message":   "Avatar berhasil diupload",
		"avatarUrl": user.Avatar,
	})
}

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

	// Path file avatar
	avatarPath := fmt.Sprintf("/app%s", user.Avatar) // Pastikan path sesuai lokasi file

	// Hapus file avatar jika ada
	if user.Avatar != "" {
		if _, err := os.Stat(avatarPath); err == nil {
			err := os.Remove(avatarPath)
			if err != nil {
				return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal menghapus avatar: "+err.Error())
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
