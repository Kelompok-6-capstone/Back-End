package controller

import (
	"calmind/helper"
	"calmind/model"
	"calmind/service"
	usecase "calmind/usecase/profile"
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

	// Update URL avatar dan delete_url di database
	user := model.User{
		Avatar:    imageURL,
		DeleteURL: deleteURL, // Pastikan kolom ini tersedia di model dan database
	}
	_, err = c.ProfilUsecase.UpdateUserProfile(userID, &user)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengupdate avatar: "+err.Error())
	}

	return helper.JSONSuccessResponse(ctx, map[string]string{
		"message":   "Avatar berhasil diupload",
		"avatarUrl": imageURL,
	})
}

// Delete Avatar
func (c *ProfilController) DeleteAvatar(ctx echo.Context) error {
	claims, ok := ctx.Get("user").(*service.JwtCustomClaims)
	if !ok {
		return helper.JSONErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized")
	}
	userID := claims.UserID

	// Ambil data user untuk mendapatkan avatar URL dan delete_url
	user, err := c.ProfilUsecase.GetUserProfile(userID)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil profil user: "+err.Error())
	}

	// Hapus avatar dari ImgBB
	if user.DeleteURL != "" {
		err := helper.DeleteFromImgBB(user.DeleteURL)
		if err != nil {
			return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal menghapus avatar dari ImgBB: "+err.Error())
		}
	}

	// Update avatar menjadi kosong di database
	user.Avatar = ""
	user.DeleteURL = ""
	_, err = c.ProfilUsecase.UpdateUserProfile(userID, user)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengupdate avatar di database: "+err.Error())
	}

	return helper.JSONSuccessResponse(ctx, "Avatar berhasil dihapus")
}
