package controller

import (
	"calmind/helper"
	"calmind/model"
	"calmind/service"
	"calmind/usecase"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type ConsultationController struct {
	ConsultationUsecase usecase.ConsultationUsecase
}

func NewConsultationController(consultationUsecase usecase.ConsultationUsecase) *ConsultationController {
	return &ConsultationController{ConsultationUsecase: consultationUsecase}
}

// Membuat konsultasi
func (c *ConsultationController) CreateConsultation(ctx echo.Context) error {
	claims, _ := ctx.Get("user").(*service.JwtCustomClaims)

	var request struct {
		DoctorID int    `json:"doctor_id"`
		Message  string `json:"message"`
	}

	if err := ctx.Bind(&request); err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Input tidak valid")
	}

	// Ambil harga dokter
	doctor, err := c.ConsultationUsecase.GetDoctorByID(request.DoctorID)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil data dokter")
	}

	consultation := model.Consultation{
		UserID:     claims.UserID,
		DoctorID:   request.DoctorID,
		Message:    request.Message,
		IsPaid:     false,
		IsApproved: false,
		StartTime:  time.Now(),
		Duration:   120, // Default duration: 2 hours
		Status:     "pending",
	}

	err = c.ConsultationUsecase.CreateConsultation(&consultation)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal membuat konsultasi")
	}

	return helper.JSONSuccessResponse(ctx, map[string]interface{}{
		"message":    "Konsultasi berhasil dibuat. Lanjutkan ke pembayaran.",
		"price":      doctor.Price,
		"duration":   consultation.Duration,
		"start_time": consultation.StartTime,
		"status":     consultation.Status,
	})
}

// Membayar konsultasi
func (c *ConsultationController) PayConsultation(ctx echo.Context) error {
	consultationID, _ := strconv.Atoi(ctx.Param("id"))

	err := c.ConsultationUsecase.UpdatePaymentStatus(consultationID, true)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal memproses pembayaran")
	}

	return helper.JSONSuccessResponse(ctx, "Pembayaran berhasil dilakukan. Menunggu persetujuan admin.")
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
		if consultation.IsApproved && consultation.IsPaid {
			user := map[string]interface{}{
				"id":        consultation.User.ID,
				"nama":      consultation.User.Username,
				"usia":      consultation.User.Tgl_lahir,
				"pekerjaan": consultation.User.Pekerjaan,
			}
			responseData = append(responseData, map[string]interface{}{
				"user":        user,
				"message":     consultation.Message,
				"start_time":  consultation.StartTime,
				"duration":    consultation.Duration,
				"rekomendasi": consultation.Rekomendasi,
			})
		}
	}

	return helper.JSONSuccessResponse(ctx, responseData)
}

// Mendapatkan detail konsultasi pasien
func (c *ConsultationController) GetConsultationDetailByID(ctx echo.Context) error {
	consultationID, _ := strconv.Atoi(ctx.Param("id"))

	consultation, err := c.ConsultationUsecase.GetConsultationByID(consultationID)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Konsultasi tidak ditemukan")
	}

	// Validasi jika konsultasi sudah kedaluwarsa
	now := time.Now()
	endTime := consultation.StartTime.Add(time.Duration(consultation.Duration) * time.Minute)
	if now.After(endTime) {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Konsultasi sudah kedaluwarsa")
	}

	responseData := map[string]interface{}{
		"user": map[string]interface{}{
			"nama":      consultation.User.Username,
			"usia":      consultation.User.Tgl_lahir,
			"pekerjaan": consultation.User.Pekerjaan,
		},
		"message":     consultation.Message,
		"rekomendasi": consultation.Rekomendasi,
	}

	return helper.JSONSuccessResponse(ctx, responseData)
}

// Mendapatkan konsultasi yang menunggu persetujuan (admin)
func (c *ConsultationController) GetPendingApprovals(ctx echo.Context) error {
	consultations, err := c.ConsultationUsecase.GetPendingApprovals()
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal mendapatkan konsultasi yang menunggu persetujuan")
	}

	var responseData []map[string]interface{}
	for _, consultation := range consultations {
		responseData = append(responseData, map[string]interface{}{
			"consultation_id": consultation.ID,
			"user_email":      consultation.User.Email,
			"doctor_id":       consultation.DoctorID,
			"status":          consultation.Status,
			"start_time":      consultation.StartTime,
		})
	}

	return helper.JSONSuccessResponse(ctx, responseData)
}

// Menyetujui pembayaran (admin)
func (c *ConsultationController) ApprovePayment(ctx echo.Context) error {
	consultationID, _ := strconv.Atoi(ctx.Param("id"))

	err := c.ConsultationUsecase.UpdateApprovalStatus(consultationID, true)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal menyetujui pembayaran")
	}

	return helper.JSONSuccessResponse(ctx, "Pembayaran berhasil disetujui")
}

// Menandai konsultasi yang sudah kedaluwarsa
func (c *ConsultationController) MarkExpiredConsultations(ctx echo.Context) error {
	err := c.ConsultationUsecase.MarkExpiredConsultations()
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal menandai konsultasi yang kedaluwarsa")
	}

	return helper.JSONSuccessResponse(ctx, "Konsultasi yang kedaluwarsa berhasil diperbarui")
}
