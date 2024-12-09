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
	StartTime   string               `json:"start_time"`
	CreatedAt   string               `json:"created_at"`
	UpdatedAt   string               `json:"updated_at"`
	User        UserData             `json:"user"`
	Doctor      DoctorData           `json:"doctor"`
	Rekomendasi []RecommendationData `json:"rekomendasi,omitempty"`
}

type RecommendationData struct {
	ID          int    `json:"id"`
	Rekomendasi string `json:"rekomendasi"`
}

// DTO untuk respons admin
type AdminConsultationResponse struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Duration    int    `json:"duration"`
	StartTime   string `json:"start_time,omitempty"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// DTO untuk data user
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

// **1. Membuat Konsultasi**
func (c *ConsultationController) CreateConsultation(ctx echo.Context) error {
	claims, ok := ctx.Get("user").(*service.JwtCustomClaims)
	if !ok || claims == nil {
		return helper.JSONErrorResponse(ctx, http.StatusUnauthorized, "User is not authorized")
	}

	var request struct {
		DoctorID    int    `json:"doctor_id"`
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	if err := ctx.Bind(&request); err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Invalid input")
	}

	err := c.ConsultationUsecase.CreateConsultation(claims.UserID, request.DoctorID, request.Title, request.Description)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Failed to create consultation: "+err.Error())
	}

	return helper.JSONSuccessResponse(ctx, map[string]interface{}{
		"message": "Consultation created successfully",
	})
}

// **2. Membayar Konsultasi**
func (c *ConsultationController) PayConsultation(ctx echo.Context) error {
	claims, ok := ctx.Get("user").(*service.JwtCustomClaims)
	if !ok || claims == nil {
		return helper.JSONErrorResponse(ctx, http.StatusUnauthorized, "User is not authorized")
	}

	consultationID, _ := strconv.Atoi(ctx.Param("id"))
	var request struct {
		Amount float64 `json:"amount"`
	}

	if err := ctx.Bind(&request); err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Invalid input")
	}

	err := c.ConsultationUsecase.PayConsultation(claims.UserID, consultationID, request.Amount)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Failed to process payment: "+err.Error())
	}

	return helper.JSONSuccessResponse(ctx, map[string]interface{}{
		"message": "Payment successful. Waiting for admin approval.",
	})
}

// **3. Menyetujui Pembayaran (Admin)**
func (c *ConsultationController) ApprovePayment(ctx echo.Context) error {
	consultationID, _ := strconv.Atoi(ctx.Param("id"))

	err := c.ConsultationUsecase.ApprovePayment(0, consultationID)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Failed to approve payment: "+err.Error())
	}

	return helper.JSONSuccessResponse(ctx, map[string]interface{}{
		"message": "Payment approved successfully",
	})
}

// **4. Mendapatkan Daftar Konsultasi untuk Dokter**
func (c *ConsultationController) GetConsultationsForDoctor(ctx echo.Context) error {
	claims, ok := ctx.Get("doctor").(*service.JwtCustomClaims)
	if !ok || claims == nil {
		return helper.JSONErrorResponse(ctx, http.StatusUnauthorized, "Doctor is not authorized")
	}

	consultations, err := c.ConsultationUsecase.GetConsultationsForDoctor(claims.UserID)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Failed to fetch consultations: "+err.Error())
	}

	// Mapping ke DTO
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

	return helper.JSONSuccessResponse(ctx, map[string]interface{}{
		"message": "Consultations fetched successfully",
		"data":    responses,
	})
}

// **5. Mendapatkan Detail Konsultasi**
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
		StartTime:   consultation.StartTime.Format("2006-01-02 15:04:05"),
		CreatedAt:   consultation.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   consultation.UpdatedAt.Format("2006-01-02 15:04:05"),
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
	}

	return helper.JSONSuccessResponse(ctx, map[string]interface{}{
		"message": "Consultation details fetched successfully",
		"data":    response,
	})
}

// **6. Memberikan Rekomendasi**
func (c *ConsultationController) AddRecommendation(ctx echo.Context) error {
	claims, ok := ctx.Get("doctor").(*service.JwtCustomClaims)
	if !ok || claims == nil {
		return helper.JSONErrorResponse(ctx, http.StatusUnauthorized, "Doctor is not authorized")
	}

	consultationID, _ := strconv.Atoi(ctx.Param("id"))
	var request struct {
		Recommendation string `json:"recommendation"`
	}

	if err := ctx.Bind(&request); err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Invalid input")
	}

	err := c.ConsultationUsecase.AddRecommendation(claims.UserID, consultationID, request.Recommendation)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Failed to add recommendation: "+err.Error())
	}

	return helper.JSONSuccessResponse(ctx, map[string]interface{}{
		"message": "Recommendation added successfully",
	})
}

// **7. Menandai Konsultasi yang Kadaluwarsa**
func (c *ConsultationController) MarkExpiredConsultations(ctx echo.Context) error {
	err := c.ConsultationUsecase.MarkExpiredConsultations()
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Failed to mark expired consultations: "+err.Error())
	}

	return helper.JSONSuccessResponse(ctx, map[string]interface{}{
		"message": "Expired consultations marked successfully",
	})
}

// **8. Mendapatkan Detail untuk Admin**
func (c *ConsultationController) ViewConsultationDetailsForAdmin(ctx echo.Context) error {
	consultationID, _ := strconv.Atoi(ctx.Param("id"))

	consultation, err := c.ConsultationUsecase.ViewConsultationDetailsForAdmin(consultationID)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Failed to fetch consultation details: "+err.Error())
	}

	// Respons hanya menampilkan data penting
	response := AdminConsultationResponse{
		ID:          consultation.ID,
		Title:       consultation.Title,
		Description: consultation.Description,
		Status:      consultation.Status,
		Duration:    consultation.Duration,
		StartTime:   consultation.StartTime.Format("2006-01-02 15:04:05"),
		CreatedAt:   consultation.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   consultation.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	return helper.JSONSuccessResponse(ctx, map[string]interface{}{
		"message": "Consultation details fetched successfully",
		"data":    response,
	})
}
