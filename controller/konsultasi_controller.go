package controller

import (
	"calmind/helper"
	"calmind/model"
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
		"tautan_bayar": paymentURL,
		"dokter": model.DoctorDTO{
			Username: consultation.Doctor.Username,
			Email:    consultation.Doctor.Email,
			Avatar:   consultation.Doctor.Avatar,
			Price:    consultation.Doctor.Price,
		},
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

	consultation, err := c.ConsultationUsecase.GetConsultationByID(consultationID)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusNotFound, "Konsultasi tidak ditemukan.")
	}

	if consultation.PaymentStatus != "settlement" {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Pembayaran belum selesai.")
	}

	err = c.ConsultationUsecase.ApprovePayment(claims.UserID, consultationID)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal menyetujui pembayaran: "+err.Error())
	}

	return helper.JSONSuccessResponse(ctx, "Pembayaran berhasil disetujui dan konsultasi diaktifkan.")
}

func (c *ConsultationController) PaymentNotification(ctx echo.Context) error {
	var notificationPayload map[string]interface{}

	if err := ctx.Bind(&notificationPayload); err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Invalid payload")
	}

	orderID, exists := notificationPayload["order_id"].(string)
	if !exists {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Order ID not found")
	}

	transactionStatusResp, err := c.ConsultationUsecase.VerifyPaymentStatus(orderID)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Failed to verify payment status")
	}

	err = c.ConsultationUsecase.UpdatePaymentStatus(orderID, transactionStatusResp.TransactionStatus)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Failed to update payment status")
	}

	return helper.JSONSuccessResponse(ctx, "Payment status updated successfully")
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

	var consultationsDTO []model.ConsultationDTO
	for _, consultation := range consultations {
		if consultation.Duration == 0 {
			consultation.Duration = 120 // Set default jika 0
		}
		consultationsDTO = append(consultationsDTO, model.ConsultationDTO{
			ID:          consultation.ID,
			Title:       consultation.Title,
			Description: consultation.Description,
			Duration:    consultation.Duration,
			Status:      consultation.Status,
			User: &model.UserDTO{
				Username: consultation.User.Username,
				Email:    consultation.User.Email,
			},
			Doctor: &model.DoctorDTO{
				Username: consultation.Doctor.Username,
				Email:    consultation.Doctor.Email,
				Avatar:   consultation.Doctor.Avatar,
			},
		})
	}

	return helper.JSONSuccessResponse(ctx, consultationsDTO)
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

	consultationDTO := model.ConsultationDTO{
		ID:          consultation.ID,
		Title:       consultation.Title,
		Description: consultation.Description,
		Duration:    consultation.Duration,
		Status:      consultation.Status,
		StartTime:   consultation.StartTime.Format("2006-01-02 15:04:05"),
		User: &model.UserDTO{
			Username: consultation.User.Username,
			Email:    consultation.User.Email,
			Avatar:   consultation.User.Avatar,
		},
		Rekomendasi: []model.RecommendationDTO{},
	}
	for _, recommendation := range consultation.Rekomendasi {
		consultationDTO.Rekomendasi = append(consultationDTO.Rekomendasi, model.RecommendationDTO{
			ID:             recommendation.ID,
			ConsultationID: recommendation.ConsultationID,
			DoctorID:       recommendation.DoctorID,
			Recommendation: recommendation.Rekomendasi,
		})
	}
	return helper.JSONSuccessResponse(ctx, consultationDTO)

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

	var consultationsDTO []model.UserDTO
	for _, consultation := range consultations {
		consultationsDTO = append(consultationsDTO, model.UserDTO{
			Username: consultation.User.Username,
			Email:    consultation.User.Email,
			Avatar:   consultation.User.Avatar,
		})
	}
	return helper.JSONSuccessResponse(ctx, consultationsDTO)

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

	consultationDTO := model.ConsultationDTO{
		ID:          consultation.ID,
		Title:       consultation.Title,
		Description: consultation.Description,
		Duration:    consultation.Duration,
		Status:      consultation.Status,
		StartTime:   consultation.StartTime.Format("2006-01-02 15:04:05"),
		User: &model.UserDTO{
			Username: consultation.User.Username,
			Email:    consultation.User.Email,
			Avatar:   consultation.User.Avatar,
		},
		Rekomendasi: []model.RecommendationDTO{},
	}
	for _, recommendation := range consultation.Rekomendasi {
		consultationDTO.Rekomendasi = append(consultationDTO.Rekomendasi, model.RecommendationDTO{
			ID:             recommendation.ID,
			ConsultationID: recommendation.ConsultationID,
			DoctorID:       recommendation.DoctorID,
			Recommendation: recommendation.Rekomendasi,
		})
	}
	return helper.JSONSuccessResponse(ctx, consultationDTO)

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

	consultationDTO := model.ConsultationDTO{
		ID:          consultation.ID,
		Title:       consultation.Title,
		Description: consultation.Description,
		Status:      consultation.Status,
		User: &model.UserDTO{
			Username: consultation.User.Username,
			Email:    consultation.User.Email,
			Avatar:   consultation.User.Avatar,
		},
		Doctor: &model.DoctorDTO{
			Username: consultation.Doctor.Username,
			Email:    consultation.Doctor.Email,
			Avatar:   consultation.Doctor.Avatar,
		},
	}
	return helper.JSONSuccessResponse(ctx, consultationDTO)

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

	var consultationsDTO []model.ConsultationDTO
	for _, consultation := range consultations {
		consultationsDTO = append(consultationsDTO, model.ConsultationDTO{
			ID:          consultation.ID,
			Title:       consultation.Title,
			Description: consultation.Description,
			Status:      consultation.Status,
			User: &model.UserDTO{
				Username: consultation.User.Username,
				Email:    consultation.User.Email,
				Avatar:   consultation.User.Avatar,
			},
			Doctor: &model.DoctorDTO{
				Username: consultation.Doctor.Username,
				Email:    consultation.Doctor.Email,
				Avatar:   consultation.Doctor.Avatar,
			},
		})
	}
	return helper.JSONSuccessResponse(ctx, consultationsDTO)

}

// Melihat daftar pembayaran pending (Admin)
// Melihat daftar pembayaran pending (Admin)
func (c *ConsultationController) GetPendingPayments(ctx echo.Context) error {
	// Validasi klaim admin
	claims, ok := ctx.Get("admin").(*service.JwtCustomClaims)
	if !ok || claims == nil {
		return helper.JSONErrorResponse(ctx, http.StatusUnauthorized, "Admin tidak memiliki akses.")
	}

	// Ambil daftar pembayaran pending
	pendingPayments, err := c.ConsultationUsecase.GetPendingPaymentsForAdmin()
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil daftar pembayaran pending.")
	}

	// Filter data untuk hanya menampilkan informasi yang diperlukan
	var filteredPayments []map[string]interface{}
	for _, consultation := range pendingPayments {
		filteredPayments = append(filteredPayments, map[string]interface{}{
			"id":            consultation.ID,
			"title":         consultation.Title,
			"description":   consultation.Description,
			"paymentStatus": consultation.PaymentStatus,
			"status":        consultation.Status,
			"user": map[string]interface{}{
				"name":  consultation.User.Username,
				"email": consultation.User.Email,
				"no_hp": consultation.User.NoHp,
			},
			"doctor": map[string]interface{}{
				"name":  consultation.Doctor.Username,
				"email": consultation.Doctor.Email,
				"no_hp": consultation.Doctor.NoHp,
			},
			"startTime": consultation.StartTime.Format("2006-01-02 15:04:05"),
			"createdAt": consultation.CreatedAt.Format("2006-01-02 15:04:05"),
			"updatedAt": consultation.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	// Return data yang sudah difilter
	return helper.JSONSuccessResponse(ctx, filteredPayments)
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

	consultation, err := c.ConsultationUsecase.GetConsultationByID(consultationID)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusNotFound, "Konsultasi tidak ditemukan.")
	}

	paymentDetails, err := c.ConsultationUsecase.GetPaymentDetails(consultationID)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil detail pembayaran.")
	}

	return helper.JSONSuccessResponse(ctx, map[string]interface{}{
		"consultation": model.ConsultationDTO{
			ID:     consultation.ID,
			Status: consultation.Status,
			User: &model.UserDTO{
				Username: consultation.User.Username,
				Email:    consultation.User.Email,
				Avatar:   consultation.User.Avatar,
			},
			Doctor: &model.DoctorDTO{
				Username: consultation.Doctor.Username,
				Email:    consultation.Doctor.Email,
				Avatar:   consultation.Doctor.Avatar,
			},
		},
		"payment": paymentDetails,
	})
}
