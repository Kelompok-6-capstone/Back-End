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

// DTO untuk respons daftar konsultasi
type ConsultationListResponse struct {
	ID        int      `json:"id"`
	Title     string   `json:"title"`
	User      UserData `json:"user"`
	Status    string   `json:"status"`
	Duration  int      `json:"duration"`
	StartTime string   `json:"start_time,omitempty"`
	CreatedAt string   `json:"created_at"`
}

// DTO untuk respons detail konsultasi
type ConsultationDetailResponse struct {
	ID          int                  `json:"id"`
	Title       string               `json:"title"`
	Description string               `json:"description"`
	Status      string               `json:"status"`
	Duration    int                  `json:"duration"`
	StartTime   string               `json:"start_time,omitempty"`
	CreatedAt   string               `json:"created_at"`
	UpdatedAt   string               `json:"updated_at"`
	User        UserData             `json:"user,omitempty"`
	Doctor      DoctorData           `json:"doctor,omitempty"`
	Rekomendasi []RecommendationData `json:"rekomendasi,omitempty"`
}

// DTO untuk rekomendasi
type RecommendationData struct {
	ID          int    `json:"id"`
	Rekomendasi string `json:"rekomendasi"`
}

// DTO untuk data pengguna
type UserData struct {
	Name      string `json:"name"`
	Avatar    string `json:"avatar"`
	Pekerjaan string `json:"pekerjaan"`
}

// DTO untuk data dokter
type DoctorData struct {
	Name       string `json:"name"`
	Avatar     string `json:"avatar"`
	Title      string `json:"title"`
	Experience int    `json:"experience"`
	About      string `json:"about"`
}

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

	// Validasi input
	if request.DoctorID <= 0 || request.Title == "" || request.Description == "" {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Semua input harus diisi")
	}

	// Buat konsultasi
	paymentURL, consultation, err := c.ConsultationUsecase.CreateConsultation(claims.UserID, request.DoctorID, request.Title, request.Description, claims.Email)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal membuat konsultasi dan link pembayaran: "+err.Error())
	}

	// Kirim respons dengan detail konsultasi dan link pembayaran
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

func (c *ConsultationController) ApprovePayment(ctx echo.Context) error {
	consultationID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "ID konsultasi tidak valid")
	}

	// Verifikasi status pembayaran dari Midtrans
	paymentStatus, err := c.ConsultationUsecase.VerifyPayment(consultationID)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal memverifikasi pembayaran: "+err.Error())
	}

	if paymentStatus != "settlement" {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Pembayaran belum selesai")
	}

	// Setujui pembayaran setelah diverifikasi
	err = c.ConsultationUsecase.ApprovePayment(consultationID)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal menyetujui pembayaran")
	}

	return helper.JSONSuccessResponse(ctx, "Pembayaran berhasil disetujui.")
}

func (c *ConsultationController) GetUserConsultations(ctx echo.Context) error {
	claims, ok := ctx.Get("user").(*service.JwtCustomClaims)
	if !ok || claims == nil {
		return helper.JSONErrorResponse(ctx, http.StatusUnauthorized, "Pengguna tidak diizinkan")
	}

	consultations, err := c.ConsultationUsecase.GetUserConsultations(claims.UserID)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil daftar konsultasi")
	}

	var responses []ConsultationListResponse
	for _, consultation := range consultations {
		responses = append(responses, ConsultationListResponse{
			ID:        consultation.ID,
			Title:     consultation.Title,
			Status:    consultation.Status,
			Duration:  consultation.Duration,
			StartTime: consultation.StartTime.Format("2006-01-02 15:04:05"),
			CreatedAt: consultation.CreatedAt.Format("2006-01-02 15:04:05"),
			User: UserData{
				Name:      consultation.Doctor.Username,
				Avatar:    consultation.Doctor.Avatar,
				Pekerjaan: consultation.Doctor.Title.Name,
			},
		})
	}

	return helper.JSONSuccessResponse(ctx, responses)
}

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

	response := map[string]interface{}{
		"id":          consultation.ID,
		"title":       consultation.Title,
		"description": consultation.Description,
		"status":      consultation.Status,
		"doctor": map[string]interface{}{
			"name":  consultation.Doctor.Username,
			"price": consultation.Doctor.Price,
		},
		"created_at": consultation.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	return helper.JSONSuccessResponse(ctx, response)
}

func (c *ConsultationController) GetConsultationsForDoctor(ctx echo.Context) error {
	claims, ok := ctx.Get("doctor").(*service.JwtCustomClaims)
	if !ok || claims == nil {
		return helper.JSONErrorResponse(ctx, http.StatusUnauthorized, "Dokter tidak diizinkan")
	}

	consultations, err := c.ConsultationUsecase.GetConsultationsForDoctor(claims.UserID)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil daftar konsultasi")
	}

	var responses []ConsultationListResponse
	for _, consultation := range consultations {
		responses = append(responses, ConsultationListResponse{
			ID:        consultation.ID,
			Title:     consultation.Title,
			StartTime: consultation.StartTime.Format("2006-01-02 15:04:05"),
			CreatedAt: consultation.CreatedAt.Format("2006-01-02 15:04:05"),
			User: UserData{
				Name:      consultation.User.Username,
				Avatar:    consultation.User.Avatar,
				Pekerjaan: consultation.User.Pekerjaan,
			},
		})
	}

	return helper.JSONSuccessResponse(ctx, responses)
}

func (c *ConsultationController) GetPendingConsultations(ctx echo.Context) error {
	consultations, err := c.ConsultationUsecase.GetPendingConsultationsForAdmin()
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil daftar konsultasi yang belum disetujui")
	}

	var responses []ConsultationListResponse
	for _, consultation := range consultations {
		responses = append(responses, ConsultationListResponse{
			ID:        consultation.ID,
			Title:     consultation.Title,
			Status:    consultation.Status,
			Duration:  consultation.Duration,
			StartTime: consultation.StartTime.Format("2006-01-02 15:04:05"),
			CreatedAt: consultation.CreatedAt.Format("2006-01-02 15:04:05"),
			User: UserData{
				Name:      consultation.User.Username,
				Avatar:    consultation.User.Avatar,
				Pekerjaan: consultation.User.Pekerjaan,
			},
		})
	}

	return helper.JSONSuccessResponse(ctx, responses)
}

func (c *ConsultationController) ViewConsultationDetailsForAdmin(ctx echo.Context) error {
	consultationID, _ := strconv.Atoi(ctx.Param("id"))

	consultation, err := c.ConsultationUsecase.ViewConsultationDetailsForAdmin(consultationID)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil detail konsultasi")
	}

	response := ConsultationDetailResponse{
		ID:          consultation.ID,
		Title:       consultation.Title,
		Description: consultation.Description,
		Status:      consultation.Status,
		Duration:    consultation.Duration,
		StartTime:   consultation.StartTime.Format("2006-01-02 15:04:05"),
		CreatedAt:   consultation.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   consultation.UpdatedAt.Format("2006-01-02 15:04:05"),
		Doctor: DoctorData{
			Name:       consultation.Doctor.Username,
			Avatar:     consultation.Doctor.Avatar,
			Title:      consultation.Doctor.Title.Name,
			Experience: consultation.Doctor.Experience,
			About:      consultation.Doctor.About,
		},
	}

	return helper.JSONSuccessResponse(ctx, response)
}

func (c *ConsultationController) ViewConsultationDetails(ctx echo.Context) error {
	claims, ok := ctx.Get("doctor").(*service.JwtCustomClaims)
	if !ok || claims == nil {
		return helper.JSONErrorResponse(ctx, http.StatusUnauthorized, "Doctor is not authorized")
	}

	consultationID, _ := strconv.Atoi(ctx.Param("id"))
	consultation, err := c.ConsultationUsecase.ViewConsultationDetails(claims.UserID, consultationID)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Failed to fetch consultation details: "+err.Error())
	}

	// Mapping ke DTO
	var rekomendasiList []RecommendationData
	for _, rekomendasi := range consultation.Rekomendasi {
		rekomendasiList = append(rekomendasiList, RecommendationData{
			ID:          rekomendasi.ID,
			Rekomendasi: rekomendasi.Rekomendasi, // Field Content dari model.Rekomendasi
		})
	}

	response := ConsultationDetailResponse{
		ID:          consultation.ID,
		Title:       consultation.Title,
		Description: consultation.Description,
		Status:      consultation.Status,
		Duration:    consultation.Duration,
		User: UserData{
			Name:      consultation.User.Username,
			Avatar:    consultation.User.Avatar,
			Pekerjaan: consultation.User.Pekerjaan,
		},
		Doctor: DoctorData{
			Name:       consultation.Doctor.Username,
			Avatar:     consultation.Doctor.Avatar,
			Title:      consultation.Doctor.Title.Name,
			Experience: consultation.Doctor.Experience,
			About:      consultation.Doctor.About,
		},
		Rekomendasi: rekomendasiList, // Isi array rekomendasi yang telah dipetakan
		StartTime:   consultation.StartTime.Format("2006-01-02 15:04:05"),
		CreatedAt:   consultation.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   consultation.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	return helper.JSONSuccessResponse(ctx, response)
}

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
