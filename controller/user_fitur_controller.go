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
	ID          int            `json:"id"`
	Username    string         `json:"username"`
	Avatar      string         `json:"avatar"`
	DateOfBirth string         `json:"date_of_birth"`
	Address     string         `json:"address"`
	Schedule    string         `json:"schedule"`
	Title       string         `json:"title"`
	Price       float64        `json:"price"`
	Experience  int            `json:"experience"`
	STRNumber   string         `json:"str_number"`
	About       string         `json:"about"`
	IsActive    bool           `json:"is_active"`
	Tags        []TagsResponse `json:"tags"`
}

type TagsResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type TitleResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Endpoint untuk mendapatkan daftar semua dokter
func (c *UserFiturController) GetDoctors(ctx echo.Context) error {
	doctors, err := c.UserFiturUsecase.GetAllDoctors()
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal mendapatkan daftar dokter: "+err.Error())
	}

	var doctorList []DoctorResponse
	for _, doctor := range doctors {
		doctorList = append(doctorList, DoctorResponse{
			ID:         doctor.ID,
			Username:   doctor.Username,
			Title:      doctor.Title.Name, // Ambil Name dari objek Title
			Experience: doctor.Experience,
			Price:      doctor.Price,
			Avatar:     doctor.Avatar,
		})
	}

	return helper.JSONSuccessResponse(ctx, doctorList)
}

// Endpoint untuk mendapatkan daftar dokter berdasarkan tag
func (c *UserFiturController) GetDoctorsByTag(ctx echo.Context) error {
	tag := ctx.QueryParam("tag")
	if tag == "" {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Tag tidak boleh kosong")
	}

	doctors, err := c.UserFiturUsecase.GetDoctorsByTag(tag)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal mendapatkan dokter berdasarkan tag: "+err.Error())
	}

	var doctorList []DoctorResponse
	for _, doctor := range doctors {
		doctorList = append(doctorList, DoctorResponse{
			ID:         doctor.ID,
			Username:   doctor.Username,
			Title:      doctor.Title.Name,
			Experience: doctor.Experience,
			Price:      doctor.Price,
			Avatar:     doctor.Avatar,
		})
	}

	return helper.JSONSuccessResponse(ctx, doctorList)
}

// Endpoint untuk mendapatkan daftar dokter berdasarkan status
func (c *UserFiturController) GetDoctorsByStatus(ctx echo.Context) error {
	status := ctx.QueryParam("status")
	if status != "active" && status != "inactive" {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Status tidak valid. Gunakan 'active' atau 'inactive'.")
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
			Title:      doctor.Title.Name,
			Experience: doctor.Experience,
			Price:      doctor.Price,
			Avatar:     doctor.Avatar,
		})
	}

	return helper.JSONSuccessResponse(ctx, doctorList)
}

// Endpoint untuk mencari dokter berdasarkan query
func (c *UserFiturController) SearchDoctors(ctx echo.Context) error {
	query := ctx.QueryParam("query")
	if query == "" {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Query tidak boleh kosong")
	}

	doctors, err := c.UserFiturUsecase.SearchDoctors(query)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal mencari dokter: "+err.Error())
	}

	var doctorList []DoctorResponse
	for _, doctor := range doctors {
		doctorList = append(doctorList, DoctorResponse{
			ID:         doctor.ID,
			Username:   doctor.Username,
			Title:      doctor.Title.Name,
			Experience: doctor.Experience,
			Price:      doctor.Price,
			Avatar:     doctor.Avatar,
		})
	}

	return helper.JSONSuccessResponse(ctx, doctorList)
}

// Endpoint untuk mendapatkan detail dokter berdasarkan ID
func (c *UserFiturController) GetDoctorDetail(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "ID dokter tidak valid")
	}

	doctor, err := c.UserFiturUsecase.GetDoctorByID(id)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusNotFound, "Dokter tidak ditemukan")
	}

	var tagsResponse []TagsResponse
	for _, tag := range doctor.Tags {
		tagsResponse = append(tagsResponse, TagsResponse{
			ID:   tag.ID,
			Name: tag.Name,
		})
	}

	doctorDetail := DoctorDetailResponse{
		ID:          doctor.ID,
		Username:    doctor.Username,
		Avatar:      doctor.Avatar,
		DateOfBirth: doctor.DateOfBirth,
		Address:     doctor.Address,
		Schedule:    doctor.Schedule,
		Title:       doctor.Title.Name, // Ambil Name dari objek Title
		Price:       doctor.Price,
		Experience:  doctor.Experience,
		STRNumber:   doctor.STRNumber,
		About:       doctor.About,
		IsActive:    doctor.IsActive,
		Tags:        tagsResponse,
	}

	return helper.JSONSuccessResponse(ctx, doctorDetail)
}

// Endpoint untuk mendapatkan semua tags
func (c *UserFiturController) GetAllTags(ctx echo.Context) error {
	tags, err := c.UserFiturUsecase.GetAllTags()
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusNotFound, "Tags tidak ditemukan")
	}

	return helper.JSONSuccessResponse(ctx, tags)
}

// Endpoint untuk mendapatkan semua titles
func (c *UserFiturController) GetAllTitles(ctx echo.Context) error {
	titles, err := c.UserFiturUsecase.GetAllTitles()
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusNotFound, "Titles tidak ditemukan: "+err.Error())
	}

	var titleList []TitleResponse
	for _, title := range titles {
		titleList = append(titleList, TitleResponse{
			ID:   title.ID,
			Name: title.Name,
		})
	}

	return helper.JSONSuccessResponse(ctx, titleList)
}

// Endpoint untuk mendapatkan daftar dokter berdasarkan title
func (c *UserFiturController) GetDoctorsByTitle(ctx echo.Context) error {
	title := ctx.QueryParam("title")
	if title == "" {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Title tidak boleh kosong")
	}

	doctors, err := c.UserFiturUsecase.GetDoctorsByTitle(title)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal mendapatkan dokter berdasarkan title: "+err.Error())
	}

	var doctorList []DoctorResponse
	for _, doctor := range doctors {
		doctorList = append(doctorList, DoctorResponse{
			ID:         doctor.ID,
			Username:   doctor.Username,
			Title:      doctor.Title.Name,
			Experience: doctor.Experience,
			Price:      doctor.Price,
			Avatar:     doctor.Avatar,
		})
	}

	return helper.JSONSuccessResponse(ctx, doctorList)
}
