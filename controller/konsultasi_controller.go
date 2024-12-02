package controller

import (
	"calmind/helper"
	"calmind/model"
	"calmind/service"
	"calmind/usecase"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ConsultationController struct {
	ConsultationUsecase usecase.ConsultationUsecase
}

func NewConsultationController(usecase usecase.ConsultationUsecase) *ConsultationController {
	return &ConsultationController{ConsultationUsecase: usecase}
}

// Buat konsultasi
func (c *ConsultationController) CreateConsultation(ctx echo.Context) error {
	claims, _ := ctx.Get("user").(*service.JwtCustomClaims)

	var request struct {
		DoctorID int    `json:"doctor_id"`
		Message  string `json:"message"`
	}

	if err := ctx.Bind(&request); err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Input tidak valid")
	}

	consultation := model.Consultation{
		UserID:   claims.UserID,
		DoctorID: request.DoctorID,
		Message:  request.Message,
	}

	err := c.ConsultationUsecase.CreateConsultation(&consultation)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal membuat konsultasi")
	}

	return helper.JSONSuccessResponse(ctx, "Konsultasi berhasil dibuat")
}
