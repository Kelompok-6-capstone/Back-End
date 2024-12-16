package controller

import (
	"calmind/helper"
	"calmind/model"
	"calmind/service"
	usecase "calmind/usecase/artikel"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

type AdminResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

type ArtikelResponse struct {
	ID        int           `json:"id"`
	Judul     string        `json:"judul"`
	Gambar    string        `json:"gambar"`
	Isi       string        `json:"isi"`
	CreatedAt string        `json:"created_at"`
	UpdatedAt string        `json:"updated_at"`
	Admin     AdminResponse `json:"admin"`
}

type ArtikelController struct {
	Usecase usecase.ArtikelUsecase
}

func NewArtikelController(usecase usecase.ArtikelUsecase) *ArtikelController {
	return &ArtikelController{Usecase: usecase}
}

// CreateArtikel - Membuat artikel baru
func (c *ArtikelController) CreateArtikel(ctx echo.Context) error {
	var artikel model.Artikel

	claims, ok := ctx.Get("admin").(*service.JwtCustomClaims)
	if !ok || claims == nil {
		return helper.JSONErrorResponse(ctx, http.StatusUnauthorized, "Klaim JWT tidak valid atau tidak ditemukan")
	}

	// Validasi input
	if err := ctx.Bind(&artikel); err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Invalid input: "+err.Error())
	}

	// Memanggil usecase untuk membuat artikel
	err := c.Usecase.CreateArtikel(claims.UserID, &artikel)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	// Berhasil
	return helper.JSONSuccessResponse(ctx, "berhasil update artikel")
}

// GetAllArtikel - Mengambil semua artikel
func (c *ArtikelController) GetAllArtikel(ctx echo.Context) error {
	artikels, err := c.Usecase.GetAllArtikel()
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	// Transform artikels ke DTO
	var responses []ArtikelResponse
	for _, artikel := range artikels {
		responses = append(responses, ArtikelResponse{
			ID:        artikel.ID,
			Judul:     artikel.Judul,
			Gambar:    artikel.Gambar,
			Isi:       artikel.Isi,
			CreatedAt: artikel.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: artikel.UpdatedAt.Format("2006-01-02 15:04:05"),
			Admin: AdminResponse{
				ID:       artikel.Admin.ID,
				Username: artikel.Admin.Username,
			},
		})
	}

	return helper.JSONSuccessResponse(ctx, responses)
}

// GetArtikelByID - Mengambil artikel berdasarkan ID
func (c *ArtikelController) GetArtikelByID(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Invalid artikel ID")
	}

	artikel, err := c.Usecase.GetArtikelByID(id)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusNotFound, err.Error())
	}

	response := ArtikelResponse{
		ID:        artikel.ID,
		Judul:     artikel.Judul,
		Gambar:    artikel.Gambar,
		Isi:       artikel.Isi,
		CreatedAt: artikel.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: artikel.UpdatedAt.Format("2006-01-02 15:04:05"),
		Admin: AdminResponse{
			ID:       artikel.Admin.ID,
			Username: artikel.Admin.Username,
		},
	}

	return helper.JSONSuccessResponse(ctx, response)
}

func (c *ArtikelController) UpdateArtikel(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Invalid artikel ID")
	}

	var artikel model.Artikel
	if err := ctx.Bind(&artikel); err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Invalid input: "+err.Error())
	}

	// Validasi agar admin_id tidak diubah melalui payload
	existingArtikel, err := c.Usecase.GetArtikelByID(id)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusNotFound, "Artikel not found")
	}

	if existingArtikel.AdminID != artikel.AdminID && artikel.AdminID != 0 {
		return helper.JSONErrorResponse(ctx, http.StatusForbidden, "Cannot change Admin ID")
	}

	artikel.ID = id
	artikel.AdminID = existingArtikel.AdminID // Tetapkan admin_id dari data asli
	err = c.Usecase.UpdateArtikel(&artikel)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	return helper.JSONSuccessResponse(ctx, map[string]string{"message": "Artikel updated successfully"})
}

// DeleteArtikel - Menghapus artikel berdasarkan ID
func (c *ArtikelController) DeleteArtikel(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Invalid artikel ID")
	}

	claims, ok := ctx.Get("admin").(*service.JwtCustomClaims)
	if !ok || claims == nil {
		return helper.JSONErrorResponse(ctx, http.StatusUnauthorized, "Klaim JWT tidak valid atau tidak ditemukan")
	}

	err = c.Usecase.DeleteArtikel(claims.UserID, id)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	return helper.JSONSuccessResponse(ctx, map[string]string{"message": "Artikel deleted successfully"})
}

func (c *ArtikelController) SearchArtikel(ctx echo.Context) error {
	query := ctx.QueryParam("query")
	if query == "" {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Query tidak boleh kosong")
	}

	// Panggil usecase untuk mencari artikel
	artikels, err := c.Usecase.SearchArtikel(query)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal mencari artikel: "+err.Error())
	}

	// Transform hasil pencarian artikel ke dalam DTO
	var responses []ArtikelResponse
	for _, artikel := range artikels {
		responses = append(responses, ArtikelResponse{
			ID:        artikel.ID,
			Judul:     artikel.Judul,
			Gambar:    artikel.Gambar,
			Isi:       artikel.Isi,
			CreatedAt: artikel.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: artikel.UpdatedAt.Format("2006-01-02 15:04:05"),
			Admin: AdminResponse{
				ID:       artikel.Admin.ID,
				Username: artikel.Admin.Username,
			},
		})
	}

	// Kirim respons sukses dengan data artikel yang diformat
	return helper.JSONSuccessResponse(ctx, responses)
}
func (c *ArtikelController) UploadArtikelImage(ctx echo.Context) error {
	// Ambil klaim admin dari JWT
	claims, ok := ctx.Get("admin").(*service.JwtCustomClaims)
	if !ok {
		return helper.JSONErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized")
	}
	adminID := claims.UserID

	// Ambil file dari form input
	file, err := ctx.FormFile("gambar")
	if err != nil || file.Size > 5*1024*1024 {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "File tidak valid atau ukuran melebihi 5 MB")
	}

	// Validasi ekstensi file
	ext := filepath.Ext(file.Filename)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Format file tidak didukung. Gunakan .jpg, .jpeg, atau .png")
	}

	// Buka file untuk diunggah
	src, err := file.Open()
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal membuka file")
	}
	defer src.Close()

	// Penamaan unik untuk file
	fileName := fmt.Sprintf("artikel_admin_%d_%s", adminID, file.Filename)

	// Upload ke Cloudinary
	imageURL, _, err := helper.UploadFileToCloudinary(src, fileName)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengunggah gambar ke Cloudinary")
	}

	// Ambil artikel ID dari parameter
	artikelID, err := strconv.Atoi(ctx.QueryParam("artikel_id"))
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Artikel ID tidak valid")
	}

	// Ambil data artikel dari database
	artikel, err := c.Usecase.GetArtikelByID(artikelID)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusNotFound, "Artikel tidak ditemukan")
	}

	// Validasi admin yang mengubah artikel
	if artikel.AdminID != adminID {
		return helper.JSONErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized")
	}

	// Update data artikel dengan URL gambar
	artikel.Gambar = imageURL
	err = c.Usecase.UpdateArtikel(artikel)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal menyimpan data artikel")
	}

	// Kirim respons sukses
	return helper.JSONSuccessResponse(ctx, map[string]string{
		"message":  "Gambar artikel berhasil diunggah",
		"imageUrl": imageURL,
	})
}

func (c *ArtikelController) DeleteArtikelImage(ctx echo.Context) error {
	// Ambil klaim admin dari JWT
	claims, ok := ctx.Get("admin").(*service.JwtCustomClaims)
	if !ok {
		return helper.JSONErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized")
	}
	adminID := claims.UserID

	// Ambil artikel ID dari parameter
	id, err := strconv.Atoi(ctx.QueryParam("artikel_id"))
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Artikel ID tidak valid")
	}

	// Ambil artikel berdasarkan ID
	artikel, err := c.Usecase.GetArtikelByID(id)
	if err != nil || artikel.AdminID != adminID {
		return helper.JSONErrorResponse(ctx, http.StatusForbidden, "Akses ditolak atau artikel tidak ditemukan")
	}

	// Ekstrak public_id dari URL gambar
	parts := strings.Split(artikel.Gambar, "/")
	publicID := strings.TrimSuffix(parts[len(parts)-1], filepath.Ext(artikel.Gambar))

	// Hapus dari Cloudinary
	if err := helper.DeleteFileFromCloudinary(publicID); err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal menghapus gambar di Cloudinary")
	}

	// Kosongkan URL gambar di database
	artikel.Gambar = ""
	err = c.Usecase.UpdateArtikel(artikel)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal menghapus URL gambar dari database")
	}

	// Kirim respons sukses
	return helper.JSONSuccessResponse(ctx, "Gambar berhasil dihapus")
}
