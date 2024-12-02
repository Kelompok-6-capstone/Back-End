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
	claims, _ := ctx.Get("doctor").(*service.JwtCustomClaims)
	doctor, err := c.DoctorProfileUsecase.GetDoctorProfile(claims.UserID)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil profil dokter: "+err.Error())
	}

	type DoctorResponse struct {
		ID        int    `json:"id"`
		Avatar    string `json:"avatar"`
		Username  string `json:"username"`
		Email     string `json:"email"`
		NoHp      string `json:"no_hp"`
		Tgl_lahir string `json:"date_of_birth"`
		Alamat    string `json:"address"`
		Schedule  string `json:"schedule"`
	}

	doctorProfile := DoctorResponse{
		ID:        doctor.ID,
		Avatar:    doctor.Avatar,
		Username:  doctor.Username,
		Email:     doctor.Email,
		NoHp:      doctor.NoHp,
		Alamat:    doctor.Address,
		Tgl_lahir: doctor.DateOfBirth,
		Schedule:  doctor.Schedule,
	}

	return helper.JSONSuccessResponse(ctx, doctorProfile)
}

// UpdateProfile updates the profile of a doctor
func (c *DoctorProfileController) UpdateProfile(ctx echo.Context) error {
	claims, _ := ctx.Get("doctor").(*service.JwtCustomClaims)
	var doctor model.Doctor

	if err := ctx.Bind(&doctor); err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Gagal mendapatkan data: "+err.Error())
	}

	// Validasi specialties
	for _, specialty := range doctor.Specialties {
		if specialty.Name == "" {
			return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Nama spesialisasi tidak boleh kosong")
		}
	}

	// Update profil dokter
	_, err := c.DoctorProfileUsecase.UpdateDoctorProfile(claims.UserID, &doctor)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengupdate profil: "+err.Error())
	}

	return helper.JSONSuccessResponse(ctx, "Berhasil Update Profil")
}

// SetActiveStatus allows a doctor to change their active/inactive status
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
