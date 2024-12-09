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

// Membuat konsultasi
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
		"data": map[string]interface{}{
			"user_id":     claims.UserID,
			"doctor_id":   request.DoctorID,
			"title":       request.Title,
			"description": request.Description,
			"status":      "pending",
			"is_paid":     false,
			"is_approved": false,
		},
	})
}

// Membayar konsultasi
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
		"data": map[string]interface{}{
			"consultation_id": consultationID,
			"paid_amount":     request.Amount,
			"status":          "waiting approval",
		},
	})
}

// Menyetujui pembayaran (admin)
func (c *ConsultationController) ApprovePayment(ctx echo.Context) error {
	consultationID, _ := strconv.Atoi(ctx.Param("id"))

	err := c.ConsultationUsecase.ApprovePayment(0, consultationID) // Admin ID can be retrieved dynamically
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Failed to approve payment: "+err.Error())
	}

	return helper.JSONSuccessResponse(ctx, map[string]interface{}{
		"message": "Payment approved successfully",
		"data": map[string]interface{}{
			"consultation_id": consultationID,
			"status":          "approved",
		},
	})
}

// Mendapatkan detail konsultasi untuk admin
func (c *ConsultationController) ViewConsultationDetailsForAdmin(ctx echo.Context) error {
	consultationID, _ := strconv.Atoi(ctx.Param("id"))

	consultation, err := c.ConsultationUsecase.ViewConsultationDetailsForAdmin(consultationID)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Failed to fetch consultation details: "+err.Error())
	}

	return helper.JSONSuccessResponse(ctx, map[string]interface{}{
		"message": "Consultation details fetched successfully",
		"data":    consultation,
	})
}

// Mendapatkan daftar konsultasi untuk dokter
func (c *ConsultationController) GetConsultationsForDoctor(ctx echo.Context) error {
	claims, ok := ctx.Get("doctor").(*service.JwtCustomClaims)
	if !ok || claims == nil {
		return helper.JSONErrorResponse(ctx, http.StatusUnauthorized, "Doctor is not authorized")
	}

	consultations, err := c.ConsultationUsecase.GetConsultationsForDoctor(claims.UserID)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Failed to fetch consultations: "+err.Error())
	}

	return helper.JSONSuccessResponse(ctx, map[string]interface{}{
		"message": "Consultations fetched successfully",
		"data":    consultations,
	})
}

// Mendapatkan detail konsultasi untuk dokter
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

	return helper.JSONSuccessResponse(ctx, map[string]interface{}{
		"message": "Consultation details fetched successfully",
		"data":    consultation,
	})
}

// Memberikan rekomendasi
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
		"data": map[string]interface{}{
			"consultation_id": consultationID,
			"recommendation":  request.Recommendation,
		},
	})
}

// Menandai konsultasi yang telah kedaluwarsa (otomatis di background job)
func (c *ConsultationController) MarkExpiredConsultations(ctx echo.Context) error {
	err := c.ConsultationUsecase.MarkExpiredConsultations()
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Failed to mark expired consultations: "+err.Error())
	}

	return helper.JSONSuccessResponse(ctx, map[string]interface{}{
		"message": "Expired consultations marked successfully",
	})
}
