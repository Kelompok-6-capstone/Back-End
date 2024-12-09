package controller

import (
	"calmind/helper"
	"calmind/service"
	"calmind/usecase"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type ConsultationController struct {
	ConsultationUsecase *usecase.ConsultationUsecaseImpl
}

func NewConsultationController(consultationUsecase *usecase.ConsultationUsecaseImpl) *ConsultationController {
	return &ConsultationController{ConsultationUsecase: consultationUsecase}
}

// 1. CreateConsultation
func (c *ConsultationController) CreateConsultation(ctx echo.Context) error {
	claims, ok := ctx.Get("user").(*service.JwtCustomClaims)
	if !ok || claims == nil {
		return helper.JSONErrorResponse(ctx, http.StatusUnauthorized, "Pengguna tidak diizinkan")
	}

	var request struct {
		DoctorID    int    `json:"doctor_id"`
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	if err := ctx.Bind(&request); err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Input tidak valid")
	}

	paymentURL, consultation, err := c.ConsultationUsecase.CreateConsultation(claims.UserID, request.DoctorID, request.Title, request.Description, claims.Email)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal membuat konsultasi dan link pembayaran: "+err.Error())
	}

	response := map[string]interface{}{
		"message": "Konsultasi berhasil dibuat. Silakan lanjutkan ke pembayaran.",
		"consultation": map[string]interface{}{
			"id":          consultation.ID,
			"title":       consultation.Title,
			"description": consultation.Description,
			"status":      consultation.Status,
			"total_price": consultation.Doctor.Price,
		},
		"payment_url": paymentURL,
	}

	return helper.JSONSuccessResponse(ctx, response)
}

// 2. ApprovePayment
func (c *ConsultationController) ApprovePayment(ctx echo.Context) error {
	consultationID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "ID konsultasi tidak valid")
	}

	paymentStatus, err := c.ConsultationUsecase.VerifyPayment(consultationID)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal memverifikasi pembayaran: "+err.Error())
	}

	if paymentStatus != "settlement" {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Pembayaran belum selesai")
	}

	err = c.ConsultationUsecase.ApprovePayment(consultationID)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal menyetujui pembayaran")
	}

	return helper.JSONSuccessResponse(ctx, "Pembayaran berhasil disetujui.")
}

// 3. GetUserConsultations
func (c *ConsultationController) GetUserConsultations(ctx echo.Context) error {
	claims, ok := ctx.Get("user").(*service.JwtCustomClaims)
	if !ok || claims == nil {
		return helper.JSONErrorResponse(ctx, http.StatusUnauthorized, "Pengguna tidak diizinkan")
	}

	consultations, err := c.ConsultationUsecase.GetUserConsultations(claims.UserID)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil daftar konsultasi")
	}

	return helper.JSONSuccessResponse(ctx, consultations)
}

// 4. GetUserConsultationDetails
func (c *ConsultationController) GetUserConsultationDetails(ctx echo.Context) error {
	claims, ok := ctx.Get("user").(*service.JwtCustomClaims)
	if !ok || claims == nil {
		return helper.JSONErrorResponse(ctx, http.StatusUnauthorized, "Pengguna tidak diizinkan")
	}

	consultationID, _ := strconv.Atoi(ctx.Param("id"))
	consultation, err := c.ConsultationUsecase.GetConsultationByID(consultationID)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusNotFound, "Konsultasi tidak ditemukan")
	}

	if consultation.UserID != claims.UserID {
		return helper.JSONErrorResponse(ctx, http.StatusForbidden, "Anda tidak memiliki akses ke konsultasi ini")
	}

	return helper.JSONSuccessResponse(ctx, consultation)
}

// 5. GetConsultationsForDoctor
func (c *ConsultationController) GetConsultationsForDoctor(ctx echo.Context) error {
	claims, ok := ctx.Get("doctor").(*service.JwtCustomClaims)
	if !ok || claims == nil {
		return helper.JSONErrorResponse(ctx, http.StatusUnauthorized, "Dokter tidak diizinkan")
	}

	consultations, err := c.ConsultationUsecase.GetConsultationsForDoctor(claims.UserID)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil daftar konsultasi")
	}

	return helper.JSONSuccessResponse(ctx, consultations)
}

// 6. GetPendingPayments
func (c *ConsultationController) GetPendingPayments(ctx echo.Context) error {
	consultations, err := c.ConsultationUsecase.GetPendingPaymentsForAdmin()
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil daftar pembayaran yang pending")
	}

	return helper.JSONSuccessResponse(ctx, consultations)
}

// 7. GetPaymentDetails
func (c *ConsultationController) GetPaymentDetails(ctx echo.Context) error {
	consultationID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "ID konsultasi tidak valid")
	}

	paymentDetails, err := c.ConsultationUsecase.GetPaymentDetails(consultationID)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil detail pembayaran")
	}

	return helper.JSONSuccessResponse(ctx, paymentDetails)
}

func (c *ConsultationController) ViewConsultationDetailsForAdmin(ctx echo.Context) error {
	consultationID, _ := strconv.Atoi(ctx.Param("id"))

	consultation, err := c.ConsultationUsecase.ViewConsultationDetailsForAdmin(consultationID)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil detail konsultasi")
	}

	return helper.JSONSuccessResponse(ctx, consultation)
}

func (c *ConsultationController) ViewConsultationDetails(ctx echo.Context) error {
	claims, ok := ctx.Get("doctor").(*service.JwtCustomClaims)
	if !ok || claims == nil {
		return helper.JSONErrorResponse(ctx, http.StatusUnauthorized, "Dokter tidak diizinkan")
	}

	consultationID, _ := strconv.Atoi(ctx.Param("id"))
	consultation, err := c.ConsultationUsecase.ViewConsultationDetails(claims.UserID, consultationID)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil detail konsultasi")
	}

	return helper.JSONSuccessResponse(ctx, consultation)
}

// 9. AddRecommendation
func (c *ConsultationController) AddRecommendation(ctx echo.Context) error {
	claims, ok := ctx.Get("doctor").(*service.JwtCustomClaims)
	if !ok || claims == nil {
		return helper.JSONErrorResponse(ctx, http.StatusUnauthorized, "Dokter tidak diizinkan")
	}

	consultationID, _ := strconv.Atoi(ctx.Param("id"))
	var request struct {
		Recommendation string `json:"recommendation"`
	}

	if err := ctx.Bind(&request); err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Input tidak valid")
	}

	err := c.ConsultationUsecase.AddRecommendation(claims.UserID, consultationID, request.Recommendation)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal menambahkan rekomendasi")
	}

	return helper.JSONSuccessResponse(ctx, "Rekomendasi berhasil ditambahkan.")
}
