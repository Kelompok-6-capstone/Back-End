package controller

import (
	"calmind/helper"
	"calmind/model"
	"calmind/service"
	"calmind/usecase"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
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

	// Konfigurasi S3 dengan kredensial otomatis dari IAM Role
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"), // Ganti sesuai region bucket S3 Anda
	})
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal membuat sesi AWS: "+err.Error())
	}

	// Buat client S3
	s3Client := s3.New(sess)

	// Nama bucket S3 dan key file
	bucketName := "calmind-data"
	fileKey := fmt.Sprintf("uploads/artikel/%d_%s", adminID, file.Filename)

	// Upload file ke S3
	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileKey),
		Body:   src,
		ACL:    aws.String("public-read"), // Set akses publik (opsional, sesuaikan kebutuhan Anda)
	})
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengupload file ke S3: "+err.Error())
	}

	// URL gambar di S3
	imageURL := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucketName, fileKey)

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

	// Periksa apakah artikel memiliki gambar
	if artikel.Gambar == "" {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Artikel tidak memiliki gambar untuk dihapus")
	}

	// Ekstrak key gambar dari URL
	fileKey := artikel.Gambar[strings.LastIndex(artikel.Gambar, "uploads/"):]

	// Konfigurasi S3 dengan kredensial otomatis dari IAM Role
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"), // Ganti sesuai region bucket S3 Anda
	})
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal membuat sesi AWS: "+err.Error())
	}

	// Buat client S3
	s3Client := s3.New(sess)

	// Hapus file dari S3
	_, err = s3Client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String("calmind-data"), // Nama bucket S3 Anda
		Key:    aws.String(fileKey),        // Path file yang akan dihapus
	})
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal menghapus file dari S3: "+err.Error())
	}

	// Update database untuk menghapus referensi gambar
	artikel.Gambar = ""
	err = c.Usecase.UpdateArtikel(artikel)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal memperbarui data artikel: "+err.Error())
	}

	return helper.JSONSuccessResponse(ctx, map[string]string{
		"message": "Gambar artikel berhasil dihapus",
	})
}
