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

// Membuat konsultasi baru (User)
func (c *ConsultationController) CreateConsultation(ctx echo.Context) error {
	claims, ok := ctx.Get("user").(*service.JwtCustomClaims)
	if !ok || claims == nil {
		return helper.JSONErrorResponse(ctx, http.StatusUnauthorized, "Pengguna tidak memiliki akses.")
	}

	var request struct {
		DoctorID    int    `json:"doctor_id"`
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	if err := ctx.Bind(&request); err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Input tidak valid. Mohon cek kembali data Anda.")
	}

	paymentURL, consultation, err := c.ConsultationUsecase.CreateConsultation(claims.UserID, request.DoctorID, request.Title, request.Description, claims.Email)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal membuat konsultasi: "+err.Error())
	}

	return helper.JSONSuccessResponse(ctx, map[string]interface{}{
		"pesan":        "Konsultasi berhasil dibuat. Silakan lanjutkan pembayaran melalui tautan berikut.",
		"konsultasi":   consultation,
		"tautan_bayar": paymentURL,
	})
}

// Menyetujui pembayaran (Admin)
func (c *ConsultationController) ApprovePayment(ctx echo.Context) error {
	claims, ok := ctx.Get("admin").(*service.JwtCustomClaims)
	if !ok || claims == nil {
		return helper.JSONErrorResponse(ctx, http.StatusUnauthorized, "Admin tidak memiliki akses.")
	}

	consultationID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "ID konsultasi tidak valid.")
	}

	err = c.ConsultationUsecase.ApprovePayment(claims.UserID, consultationID)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal menyetujui pembayaran: "+err.Error())
	}

	return helper.JSONSuccessResponse(ctx, "Pembayaran berhasil disetujui.")
}

// Melihat daftar konsultasi (User)
func (c *ConsultationController) GetUserConsultations(ctx echo.Context) error {
	claims, ok := ctx.Get("user").(*service.JwtCustomClaims)
	if !ok || claims == nil {
		return helper.JSONErrorResponse(ctx, http.StatusUnauthorized, "Pengguna tidak memiliki akses.")
	}

	consultations, err := c.ConsultationUsecase.GetUserConsultations(claims.UserID)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil daftar konsultasi.")
	}

	return helper.JSONSuccessResponse(ctx, consultations)
}

// Melihat detail konsultasi (User)
func (c *ConsultationController) GetUserConsultationDetails(ctx echo.Context) error {
	claims, ok := ctx.Get("user").(*service.JwtCustomClaims)
	if !ok || claims == nil {
		return helper.JSONErrorResponse(ctx, http.StatusUnauthorized, "Pengguna tidak memiliki akses.")
	}

	consultationID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "ID konsultasi tidak valid.")
	}

	consultation, err := c.ConsultationUsecase.GetConsultationByID(consultationID)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusNotFound, "Konsultasi tidak ditemukan.")
	}

	if consultation.UserID != claims.UserID {
		return helper.JSONErrorResponse(ctx, http.StatusForbidden, "Anda tidak memiliki akses ke konsultasi ini.")
	}

	return helper.JSONSuccessResponse(ctx, consultation)
}

// Melihat daftar konsultasi (Dokter)
func (c *ConsultationController) GetConsultationsForDoctor(ctx echo.Context) error {
	claims, ok := ctx.Get("doctor").(*service.JwtCustomClaims)
	if !ok || claims == nil {
		return helper.JSONErrorResponse(ctx, http.StatusUnauthorized, "Dokter tidak memiliki akses.")
	}

	consultations, err := c.ConsultationUsecase.GetConsultationsForDoctor(claims.UserID)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil daftar konsultasi.")
	}

	return helper.JSONSuccessResponse(ctx, consultations)
}

// Melihat detail konsultasi (Dokter)
func (c *ConsultationController) ViewConsultationDetails(ctx echo.Context) error {
	claims, ok := ctx.Get("doctor").(*service.JwtCustomClaims)
	if !ok || claims == nil {
		return helper.JSONErrorResponse(ctx, http.StatusUnauthorized, "Dokter tidak memiliki akses.")
	}

	consultationID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "ID konsultasi tidak valid.")
	}

	consultation, err := c.ConsultationUsecase.ViewConsultationDetails(claims.UserID, consultationID)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil detail konsultasi: "+err.Error())
	}

	return helper.JSONSuccessResponse(ctx, consultation)
}

// Menambahkan rekomendasi (Dokter)
func (c *ConsultationController) AddRecommendation(ctx echo.Context) error {
	claims, ok := ctx.Get("doctor").(*service.JwtCustomClaims)
	if !ok || claims == nil {
		return helper.JSONErrorResponse(ctx, http.StatusUnauthorized, "Dokter tidak memiliki akses.")
	}

	consultationID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "ID konsultasi tidak valid.")
	}

	var request struct {
		Recommendation string `json:"rekomendasi"`
	}

	if err := ctx.Bind(&request); err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Input tidak valid.")
	}

	err = c.ConsultationUsecase.AddRecommendation(claims.UserID, consultationID, request.Recommendation)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal menambahkan rekomendasi: "+err.Error())
	}

	return helper.JSONSuccessResponse(ctx, "Rekomendasi berhasil ditambahkan.")
}

// Melihat detail konsultasi (Admin)
func (c *ConsultationController) ViewConsultationDetailsForAdmin(ctx echo.Context) error {
	claims, ok := ctx.Get("admin").(*service.JwtCustomClaims)
	if !ok || claims == nil {
		return helper.JSONErrorResponse(ctx, http.StatusUnauthorized, "Admin tidak memiliki akses.")
	}

	consultationID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "ID konsultasi tidak valid.")
	}

	consultation, err := c.ConsultationUsecase.ViewConsultationDetailsForAdmin(consultationID)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil detail konsultasi: "+err.Error())
	}

	return helper.JSONSuccessResponse(ctx, consultation)
}

// Melihat daftar konsultasi pending (Admin)
func (c *ConsultationController) GetPendingConsultations(ctx echo.Context) error {
	claims, ok := ctx.Get("admin").(*service.JwtCustomClaims)
	if !ok || claims == nil {
		return helper.JSONErrorResponse(ctx, http.StatusUnauthorized, "Admin tidak memiliki akses.")
	}

	// Panggil metode yang benar dari usecase
	consultations, err := c.ConsultationUsecase.GetPendingConsultations()
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil daftar konsultasi pending.")
	}

	return helper.JSONSuccessResponse(ctx, consultations)
}

// Melihat daftar pembayaran pending (Admin)
func (c *ConsultationController) GetPendingPayments(ctx echo.Context) error {
	claims, ok := ctx.Get("admin").(*service.JwtCustomClaims)
	if !ok || claims == nil {
		return helper.JSONErrorResponse(ctx, http.StatusUnauthorized, "Admin tidak memiliki akses.")
	}

	pendingPayments, err := c.ConsultationUsecase.GetPendingPaymentsForAdmin()
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil daftar pembayaran pending.")
	}

	return helper.JSONSuccessResponse(ctx, pendingPayments)
}

// Melihat detail pembayaran (Admin)
func (c *ConsultationController) GetPaymentDetails(ctx echo.Context) error {
	claims, ok := ctx.Get("admin").(*service.JwtCustomClaims)
	if !ok || claims == nil {
		return helper.JSONErrorResponse(ctx, http.StatusUnauthorized, "Admin tidak memiliki akses.")
	}

	consultationID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "ID konsultasi tidak valid.")
	}

	// Verifikasi status pembayaran sebelum mengambil detail
	_, err = c.ConsultationUsecase.VerifyPayment(consultationID)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal memverifikasi status pembayaran: "+err.Error())
	}

	paymentDetails, err := c.ConsultationUsecase.GetPaymentDetails(consultationID)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil detail pembayaran.")
	}

	return helper.JSONSuccessResponse(ctx, paymentDetails)
}
