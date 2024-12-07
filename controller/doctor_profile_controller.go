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

// Mendapatkan profil dokter
func (c *DoctorProfileController) GetProfile(ctx echo.Context) error {
	claims, _ := ctx.Get("doctor").(*service.JwtCustomClaims)
	doctor, err := c.DoctorProfileUsecase.GetDoctorProfile(claims.UserID)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil profil dokter: "+err.Error())
	}

	// Format response lengkap dengan Title dan Tags
	type TagsResponse struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	type DoctorResponse struct {
		ID          int            `json:"id"`
		Avatar      string         `json:"avatar"`
		Username    string         `json:"username"`
		Email       string         `json:"email"`
		NoHp        string         `json:"no_hp"`
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

	var tagsResponse []TagsResponse
	for _, tag := range doctor.Tags {
		tagsResponse = append(tagsResponse, TagsResponse{
			ID:   tag.ID,
			Name: tag.Name,
		})
	}

	doctorProfile := DoctorResponse{
		ID:          doctor.ID,
		Avatar:      doctor.Avatar,
		Username:    doctor.Username,
		Email:       doctor.Email,
		NoHp:        doctor.NoHp,
		DateOfBirth: doctor.DateOfBirth,
		Address:     doctor.Address,
		Schedule:    doctor.Schedule,
		Title:       doctor.Title.Name,
		Price:       doctor.Price,
		Experience:  doctor.Experience,
		STRNumber:   doctor.STRNumber,
		About:       doctor.About,
		IsActive:    doctor.IsActive,
		Tags:        tagsResponse,
	}

	return helper.JSONSuccessResponse(ctx, doctorProfile)
}

// Memperbarui profil dokter
func (c *DoctorProfileController) UpdateProfile(ctx echo.Context) error {
	claims, _ := ctx.Get("doctor").(*service.JwtCustomClaims)

	var doctor model.Doctor
	if err := ctx.Bind(&doctor); err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Gagal mendapatkan data: "+err.Error())
	}

	_, err := c.DoctorProfileUsecase.UpdateDoctorProfile(claims.UserID, &doctor)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengupdate profil: "+err.Error())
	}

	return helper.JSONSuccessResponse(ctx, "Berhasil update profil")
}

// Mengatur status aktif dokter
func (c *DoctorProfileController) SetActiveStatus(ctx echo.Context) error {
	claims, _ := ctx.Get("doctor").(*service.JwtCustomClaims)

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
