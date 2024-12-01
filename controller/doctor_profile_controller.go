package controller

import (
	"calmind/helper"
	"calmind/model"
	"calmind/service"
	"calmind/usecase"
	"log"
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
		log.Println("Failed to fetch doctor profile:", err)
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil profil: "+err.Error())
	}

	type DoctorProfileResponse struct {
		ID       int    `json:"id"`
		Username string `json:"username"`
		Avatar   string `json:"avatar"`
		Email    string `json:"email"`
		NoHp     string `json:"no_hp"`
		Birth    string `json:"date_of_birth"`
		Address  string `json:"address"`
		Schedule string `json:"schedule"`
	}

	doctorProfile := DoctorProfileResponse{
		ID:       doctor.ID,
		Username: doctor.Username,
		Email:    doctor.Email,
		NoHp:     doctor.NoHp,
		Avatar:   doctor.Avatar,
		Birth:    doctor.DateOfBirth,
		Address:  doctor.Address,
		Schedule: doctor.Schedule,
	}

	return helper.JSONSuccessResponse(ctx, doctorProfile)
}

func (c *DoctorProfileController) UpdateProfile(ctx echo.Context) error {
	claims, _ := ctx.Get("doctor").(*service.JwtCustomClaims)
	var doctor model.Doctor

	if err := ctx.Bind(&doctor); err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Gagal mendapatkan data: "+err.Error())
	}

	// Email dan Password tidak boleh diperbarui di endpoint ini
	doctor.Email = ""
	doctor.Password = ""

	_, err := c.DoctorProfileUsecase.UpdateDoctorProfile(claims.UserID, &doctor)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengupdate profil: "+err.Error())
	}

	return helper.JSONSuccessResponse(ctx, "Berhasil update profil")
}
