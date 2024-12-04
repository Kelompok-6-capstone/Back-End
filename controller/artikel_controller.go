package controller

import (
	"calmind/helper"
	"calmind/model"
	"calmind/service"
	"calmind/usecase"
	"net/http"
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

// UpdateArtikel - Memperbarui artikel berdasarkan ID
func (c *ArtikelController) UpdateArtikel(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Invalid artikel ID")
	}

	var artikel model.Artikel
	if err := ctx.Bind(&artikel); err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Invalid input: "+err.Error())
	}

	artikel.ID = id
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
