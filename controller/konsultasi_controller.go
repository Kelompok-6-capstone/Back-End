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

func NewConsultationController(consultationUsecase usecase.ConsultationUsecase) *ConsultationController {
	return &ConsultationController{ConsultationUsecase: consultationUsecase}
}

// Membuat konsultasi baru (dari user)
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

	// Call usecase to create consultation
	err := c.ConsultationUsecase.CreateConsultation(&consultation)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal membuat konsultasi")
	}

	// Success response
	return helper.JSONSuccessResponse(ctx, "Konsultasi berhasil dibuat")
}

// Mendapatkan daftar konsultasi (untuk dokter)
func (c *ConsultationController) GetConsultationsAllDoctor(ctx echo.Context) error {
	claims, ok := ctx.Get("doctor").(*service.JwtCustomClaims)
	if !ok || claims == nil {
		return helper.JSONErrorResponse(ctx, http.StatusUnauthorized, "Doctor is not authorized")
	}

	doctorID := claims.UserID
	consultations, err := c.ConsultationUsecase.GetConsultationsAllDoctor(doctorID)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil daftar konsultasi")
	}

	var responseData []map[string]interface{}
	for _, consultation := range consultations {
		user := map[string]interface{}{
			"id":   consultation.User.ID,
			"nama": consultation.User.Username,
			"usia": consultation.User.Tgl_lahir,
			"pekerjaan":  consultation.User.Pekerjaan,
		}

		responseData = append(responseData, map[string]interface{}{"user": user})
	}

	return helper.JSONSuccessResponse(ctx, responseData)
}

// Mendapatkan detail konsultasi pasien
func (c *ConsultationController) GetConsultationDetailByID(ctx echo.Context) error {
	consultationID := ctx.Param("id")

	consultation, err := c.ConsultationUsecase.GetConsultationByID(consultationID)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Konsultasi tidak ditemukan")
	}

	responseData := map[string]interface{}{
			"user": map[string]interface{}{
			"nama": consultation.User.Username,
			"usia": consultation.User.Tgl_lahir,
			"pekerjaan":  consultation.User.Pekerjaan,
		},
		"message": consultation.Message,
		"rekomendasi": consultation.Rekomendasi,
	}

	return helper.JSONSuccessResponse(ctx, responseData)
}

// Memberikan rekomendasi kepada pasien
// Memberikan rekomendasi kepada pasien
func (c *ConsultationController) GiveRecommendation(ctx echo.Context) error {
	claims, _ := ctx.Get("doctor").(*service.JwtCustomClaims)
	doctorID := claims.UserID

	var request struct {
		Rekomendasi string `json:"rekomendasi"`
	}

	consultationID := ctx.Param("id")
	if err := ctx.Bind(&request); err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Input tidak valid")
	}
	consultation, err := c.ConsultationUsecase.GetConsultationByID(consultationID)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusNotFound, "Konsultasi tidak ditemukan")
	}
	if consultation.DoctorID != doctorID {
		return helper.JSONErrorResponse(ctx, http.StatusForbidden, "Anda tidak memiliki akses ke konsultasi ini")
	}
	err = c.ConsultationUsecase.UpdateRecommendation(consultation.ID, request.Rekomendasi)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal memberikan rekomendasi")
	}

	return helper.JSONSuccessResponse(ctx, "Rekomendasi berhasil diberikan")
}

