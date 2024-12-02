package controller

import (
	"calmind/helper"
	"calmind/usecase"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type UserFiturController struct {
	UserFiturUsecase usecase.UserFiturUsecase
}

func NewUserFiturController(UserFiturUsecase usecase.UserFiturUsecase) *UserFiturController {
	return &UserFiturController{UserFiturUsecase: UserFiturUsecase}
}

// Struct untuk respons dokter dalam daftar
type DoctorResponse struct {
	ID         int     `json:"id"`
	Username   string  `json:"username"`
	Title      string  `json:"title"`
	Experience int     `json:"experience"`
	Price      float64 `json:"price"`
	Avatar     string  `json:"avatar"`
}

// Struct untuk respons detail dokter
type DoctorDetailResponse struct {
	ID          int      `json:"id"`
	Username    string   `json:"username"`
	Avatar      string   `json:"avatar"`
	DateOfBirth string   `json:"date_of_birth"`
	Address     string   `json:"address"`
	Schedule    string   `json:"schedule"`
	Title       string   `json:"title"`
	Price       float64  `json:"price"`
	Experience  int      `json:"experience"`
	STRNumber   string   `json:"str_number"`
	About       string   `json:"about"`
	Specialties []string `json:"specialties"`
	IsActive    bool     `json:"is_active"`
}

// Endpoint untuk mendapatkan daftar semua dokter
func (c *UserFiturController) GetDoctors(ctx echo.Context) error {
	doctors, err := c.UserFiturUsecase.GetAllDoctors()
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal mendapatkan daftar dokter: "+err.Error())
	}

	// Mengisi data ke struct
	var doctorList []DoctorResponse
	for _, doctor := range doctors {
		doctorList = append(doctorList, DoctorResponse{
			ID:         doctor.ID,
			Username:   doctor.Username,
			Title:      doctor.Title,
			Experience: doctor.Experience,
			Price:      doctor.Price,
			Avatar:     doctor.Avatar,
		})
	}

	return helper.JSONSuccessResponse(ctx, map[string]interface{}{
		"data":    doctorList,
		"success": true,
	})
}

// Endpoint untuk mendapatkan daftar dokter berdasarkan spesialisasi
func (c *UserFiturController) GetDoctorsBySpecialty(ctx echo.Context) error {
	specialty := ctx.QueryParam("specialty")
	if specialty == "" {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Spesialisasi tidak boleh kosong")
	}

	doctors, err := c.UserFiturUsecase.GetDoctorsBySpecialty(specialty)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal mendapatkan dokter berdasarkan spesialisasi: "+err.Error())
	}

	var doctorList []DoctorResponse
	for _, doctor := range doctors {
		doctorList = append(doctorList, DoctorResponse{
			ID:         doctor.ID,
			Username:   doctor.Username,
			Title:      doctor.Title,
			Experience: doctor.Experience,
			Price:      doctor.Price,
			Avatar:     doctor.Avatar,
		})
	}

	return helper.JSONSuccessResponse(ctx, map[string]interface{}{
		"data":    doctorList,
		"success": true,
	})
}

// Endpoint untuk mendapatkan daftar dokter berdasarkan status
func (c *UserFiturController) GetDoctorsByStatus(ctx echo.Context) error {
	status := ctx.QueryParam("status")
	if status != "active" && status != "inactive" {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Status tidak valid")
	}

	isActive := (status == "active")
	doctors, err := c.UserFiturUsecase.GetDoctorsByStatus(isActive)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal mendapatkan dokter berdasarkan status: "+err.Error())
	}

	var doctorList []DoctorResponse
	for _, doctor := range doctors {
		doctorList = append(doctorList, DoctorResponse{
			ID:         doctor.ID,
			Username:   doctor.Username,
			Title:      doctor.Title,
			Experience: doctor.Experience,
			Price:      doctor.Price,
			Avatar:     doctor.Avatar,
		})
	}

	return helper.JSONSuccessResponse(ctx, map[string]interface{}{
		"data":    doctorList,
		"success": true,
	})
}

// Endpoint untuk mencari dokter berdasarkan query
func (c *UserFiturController) SearchDoctors(ctx echo.Context) error {
	query := ctx.QueryParam("query")
	doctors, err := c.UserFiturUsecase.SearchDoctors(query)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal mencari dokter: "+err.Error())
	}

	var doctorList []DoctorResponse
	for _, doctor := range doctors {
		doctorList = append(doctorList, DoctorResponse{
			ID:         doctor.ID,
			Username:   doctor.Username,
			Title:      doctor.Title,
			Experience: doctor.Experience,
			Price:      doctor.Price,
			Avatar:     doctor.Avatar,
		})
	}

	return helper.JSONSuccessResponse(ctx, map[string]interface{}{
		"data":    doctorList,
		"success": true,
	})
}

// Endpoint untuk mendapatkan detail dokter
func (c *UserFiturController) GetDoctorDetail(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "ID dokter tidak valid")
	}

	doctor, err := c.UserFiturUsecase.GetDoctorByID(id)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusNotFound, "Dokter tidak ditemukan")
	}

	specialties := []string{}
	for _, specialty := range doctor.Specialties {
		specialties = append(specialties, specialty.Name)
	}

	doctorDetail := DoctorDetailResponse{
		ID:          doctor.ID,
		Username:    doctor.Username,
		Avatar:      doctor.Avatar,
		DateOfBirth: doctor.DateOfBirth,
		Address:     doctor.Address,
		Schedule:    doctor.Schedule,
		Title:       doctor.Title,
		Price:       doctor.Price,
		Experience:  doctor.Experience,
		STRNumber:   doctor.STRNumber,
		About:       doctor.About,
		Specialties: specialties,
		IsActive:    doctor.IsActive,
	}

	return helper.JSONSuccessResponse(ctx, doctorDetail)
}
