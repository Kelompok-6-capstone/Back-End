package controller

import (
	"calmind/helper"
	"calmind/model"
	"calmind/service"
	"calmind/usecase"
	"net/http"

	"github.com/labstack/echo/v4"
)

type DoctorProfileController struct {
	DoctorProfileUsecase usecase.DoctorProfileUseCase
}

func NewDoctorProfileController(DoctorProfileUsecase usecase.DoctorProfileUseCase) *DoctorProfileController {
	return &DoctorProfileController{DoctorProfileUsecase: DoctorProfileUsecase}
}

func (c *DoctorProfileController) GetProfile(ctx echo.Context) error {
	claims, ok := ctx.Get("doctor").(*service.JwtCustomClaims)
	if !ok || claims == nil {
		return helper.JSONErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized access")
	}

	doctor, err := c.DoctorProfileUsecase.GetDoctorProfile(claims.UserID)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil profil dokter: "+err.Error())
	}

	// Format response lengkap dengan Title dan Tags
	type TagsResponse struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	type TitleResponse struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	type DoctorResponse struct {
		ID           int            `json:"id"`
		Username     string         `json:"username"`
		NoHp         string         `json:"no_hp"`
		Email        string         `json:"email"`
		Password     string         `json:"password"`
		Role         string         `json:"role"`
		Avatar       string         `json:"avatar"`
		DateOfBirth  string         `json:"date_of_birth"`
		Address      string         `json:"address"`
		Schedule     string         `json:"schedule"`
		IsVerified   bool           `json:"is_verified"`
		IsActive     bool           `json:"is_active"`
		Price        float64        `json:"price"`
		Experience   int            `json:"experience"`
		STRNumber    string         `json:"str_number"`
		About        string         `json:"about"`
		JenisKelamin string         `json:"jenis_kelamin"`
		Title        TitleResponse  `json:"title"`
		Tags         []TagsResponse `json:"tags"`
		CreatedAt    string         `json:"created_at"`
		UpdatedAt    string         `json:"updated_at"`
	}

	var tagsResponse []TagsResponse
	for _, tag := range doctor.Tags {
		tagsResponse = append(tagsResponse, TagsResponse{
			ID:   tag.ID,
			Name: tag.Name,
		})
	}

	doctorProfile := DoctorResponse{
		ID:           doctor.ID,
		Username:     doctor.Username,
		NoHp:         doctor.NoHp,
		Email:        doctor.Email,
		Password:     doctor.Password,
		Role:         doctor.Role,
		Avatar:       doctor.Avatar,
		DateOfBirth:  doctor.DateOfBirth,
		Address:      doctor.Address,
		Schedule:     doctor.Schedule,
		IsVerified:   doctor.IsVerified,
		IsActive:     doctor.IsActive,
		Price:        doctor.Price,
		Experience:   doctor.Experience,
		STRNumber:    doctor.STRNumber,
		About:        doctor.About,
		JenisKelamin: doctor.JenisKelamin,
		Title: TitleResponse{
			ID:   doctor.Title.ID,
			Name: doctor.Title.Name,
		},
		Tags:      tagsResponse,
		CreatedAt: doctor.CreatedAt.Format("2006-01-02T15:04:05.000Z"),
		UpdatedAt: doctor.UpdatedAt.Format("2006-01-02T15:04:05.000Z"),
	}

	return helper.JSONSuccessResponse(ctx, doctorProfile)
}

// UpdateProfile updates the profile of the logged-in doctor
func (c *DoctorProfileController) UpdateProfile(ctx echo.Context) error {
	claims, ok := ctx.Get("doctor").(*service.JwtCustomClaims)
	if !ok || claims == nil {
		return helper.JSONErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized access")
	}

	var request struct {
		Username    string       `json:"username"`
		NoHp        string       `json:"no_hp"`
		Avatar      string       `json:"avatar"`
		DateOfBirth string       `json:"date_of_birth"`
		Address     string       `json:"address"`
		Schedule    string       `json:"schedule"`
		Title       string       `json:"title"`
		Experience  int          `json:"experience"`
		STRNumber   string       `json:"str_number"`
		About       string       `json:"about"`
		Tags        []model.Tags `json:"tags"`
	}

	if err := ctx.Bind(&request); err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Gagal memproses data: "+err.Error())
	}

	doctor := &model.Doctor{
		Username:    request.Username,
		NoHp:        request.NoHp,
		Avatar:      request.Avatar,
		DateOfBirth: request.DateOfBirth,
		Address:     request.Address,
		Schedule:    request.Schedule,
		Title:       model.Title{Name: request.Title},
		Experience:  request.Experience,
		STRNumber:   request.STRNumber,
		About:       request.About,
		Tags:        request.Tags,
	}

	_, err := c.DoctorProfileUsecase.UpdateDoctorProfile(claims.UserID, doctor)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengupdate profil dokter: "+err.Error())
	}

	return helper.JSONSuccessResponse(ctx, "berhasil updated")
}

// SetActiveStatus updates the active status of the logged-in doctor
func (c *DoctorProfileController) SetActiveStatus(ctx echo.Context) error {
	claims, ok := ctx.Get("doctor").(*service.JwtCustomClaims)
	if !ok || claims == nil {
		return helper.JSONErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized access")
	}

	var statusRequest struct {
		IsActive bool `json:"is_active"`
	}

	if err := ctx.Bind(&statusRequest); err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Gagal memproses data: "+err.Error())
	}

	err := c.DoctorProfileUsecase.SetDoctorActiveStatus(claims.UserID, statusRequest.IsActive)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengubah status aktif dokter: "+err.Error())
	}

	message := "Anda telah mengubah status ke tidak aktif"
	if statusRequest.IsActive {
		message = "Anda telah mengubah status ke aktif"
	}

	return helper.JSONSuccessResponse(ctx, map[string]string{
		"message": message,
	})
}
