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
	"strconv"

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

	go func() {
		helper.Broadcast <- artikel // Ganti dari `helper.broadcast` ke `helper.Broadcast`
	}()

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
	claims, ok := ctx.Get("admin").(*service.JwtCustomClaims)
	if !ok {
		return helper.JSONErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized")
	}
	adminID := claims.UserID

	// Ambil file dari form
	file, err := ctx.FormFile("gambar")
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Gagal mendapatkan file: "+err.Error())
	}

	// Validasi ukuran file (maksimal 5 MB)
	if file.Size > 10*1024*1024 {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Ukuran file maksimal 5 MB")
	}

	// Validasi ekstensi file
	ext := filepath.Ext(file.Filename)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Hanya file dengan format .jpg, .jpeg, atau .png yang diperbolehkan")
	}

	// Simpan file di direktori uploads
	uploadDir := "uploads/artikel"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		err = os.MkdirAll(uploadDir, os.ModePerm)
		if err != nil {
			return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal membuat direktori upload")
		}
	}

	filePath := fmt.Sprintf("%s/%d_%s", uploadDir, adminID, file.Filename)
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

	// Buat URL gambar
	imageURL := fmt.Sprintf("https://%s/uploads/artikel/%d_%s", ctx.Request().Host, adminID, file.Filename)

	return helper.JSONSuccessResponse(ctx, map[string]string{
		"message":  "Gambar artikel berhasil diupload",
		"imageUrl": imageURL,
	})
}
func (c *ArtikelController) DeleteArtikelImage(ctx echo.Context) error {
	// Ambil ID artikel dari parameter
	artikelIDStr := ctx.QueryParam("artikel_id")
	if artikelIDStr == "" {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Artikel ID diperlukan")
	}

	artikelID, err := strconv.Atoi(artikelIDStr)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Artikel ID tidak valid")
	}

	// Ambil artikel dari database untuk mendapatkan URL gambar
	artikel, err := c.Usecase.GetArtikelByID(artikelID)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusNotFound, "Artikel tidak ditemukan: "+err.Error())
	}

	// Hapus file gambar artikel
	if artikel.Gambar != "" {
		filePath := "." + artikel.Gambar // Path relatif ke file
		if _, err := os.Stat(filePath); err == nil {
			if err := os.Remove(filePath); err != nil {
				return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal menghapus file gambar artikel: "+err.Error())
			}
		}
	}

	// Update URL gambar menjadi kosong di database
	artikel.Gambar = ""
	err = c.Usecase.UpdateArtikel(artikel) // Menggunakan variabel `artikel` secara langsung
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengupdate gambar artikel di database: "+err.Error())
	}

	return helper.JSONSuccessResponse(ctx, "Gambar artikel berhasil dihapus")
}
