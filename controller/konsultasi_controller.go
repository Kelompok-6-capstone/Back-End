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
	PaymentUsecase      usecase.PaymentUsecase
}

func NewConsultationController(consultationUsecase usecase.ConsultationUsecase, paymentUsecase usecase.PaymentUsecase) *ConsultationController {
	return &ConsultationController{
		ConsultationUsecase: consultationUsecase,
		PaymentUsecase:      paymentUsecase,
	}
}

// SelectSchedule - User selects a schedule for consultation
func (c *ConsultationController) SelectSchedule(ctx echo.Context) error {
	claims, _ := ctx.Get("user").(*service.JwtCustomClaims)

	var request struct {
		DoctorID int    `json:"doctor_id"`
		Date     string `json:"date"`
		Time     string `json:"time"`
	}

	if err := ctx.Bind(&request); err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Invalid input")
	}

	startTime, err := time.Parse("2006-01-02 15:04", request.Date+" "+request.Time)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Invalid date/time format")
	}
	endTime := startTime.Add(1 * time.Hour)

	consultation := model.Consultation{
		UserID:        claims.UserID,
		DoctorID:      request.DoctorID,
		StartTime:     startTime,
		EndTime:       endTime,
		Status:        "scheduled",
		PaymentStatus: "unpaid",
	}

	err = c.ConsultationUsecase.CreateConsultation(&consultation)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Failed to schedule consultation")
	}

	return helper.JSONSuccessResponse(ctx, map[string]interface{}{
		"message": "Schedule successfully selected",
		"data":    consultation,
	})
}

// AddMessageToConsultation - User adds a message to their consultation
func (c *ConsultationController) AddMessageToConsultation(ctx echo.Context) error {
	var request struct {
		ConsultationID int    `json:"consultation_id"`
		Message        string `json:"message"`
	}

	if err := ctx.Bind(&request); err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Invalid input")
	}

	err := c.ConsultationUsecase.UpdateMessage(request.ConsultationID, request.Message)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Failed to save complaint")
	}

	return helper.JSONSuccessResponse(ctx, "Complaint successfully added")
}

// CreatePayment - User makes a payment
func (c *ConsultationController) CreatePayment(ctx echo.Context) error {
	var request struct {
		ConsultationID int `json:"consultation_id"`
		Amount         int `json:"amount"`
	}

	if err := ctx.Bind(&request); err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Invalid input")
	}

	err := c.PaymentUsecase.CreatePayment(request.ConsultationID, request.Amount)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Failed to create payment")
	}

	err = c.ConsultationUsecase.UpdatePaymentStatus(request.ConsultationID, "pending")
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Failed to update payment status")
	}

	return helper.JSONSuccessResponse(ctx, "Payment successfully created")
}

// VerifyPayment - Admin verifies the payment
func (c *ConsultationController) VerifyPayment(ctx echo.Context) error {
	paymentID, _ := strconv.Atoi(ctx.Param("payment_id"))

	err := c.PaymentUsecase.VerifyPayment(paymentID)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Payment verification failed")
	}

	err = c.ConsultationUsecase.UpdateConsultationStatus(paymentID, "paid")
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Failed to update consultation status")
	}

	return helper.JSONSuccessResponse(ctx, "Payment verified successfully")
}

// GetConsultationsAllDoctor - Doctor views all consultations
func (c *ConsultationController) GetConsultationsAllDoctor(ctx echo.Context) error {
	claims, ok := ctx.Get("doctor").(*service.JwtCustomClaims)
	if !ok || claims == nil {
		return helper.JSONErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized")
	}

	consultations, err := c.ConsultationUsecase.GetConsultationsAllDoctor(claims.UserID)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve consultations")
	}

	var response []map[string]interface{}
	for _, consultation := range consultations {
		response = append(response, map[string]interface{}{
			"id":         consultation.ID,
			"start_time": consultation.StartTime,
			"end_time":   consultation.EndTime,
			"user": map[string]interface{}{
				"id":       consultation.User.ID,
				"username": consultation.User.Username,
			},
			"status": consultation.Status,
		})
	}

	return helper.JSONSuccessResponse(ctx, response)
}

// GetConsultationDetail - Doctor views consultation details
func (c *ConsultationController) GetConsultationDetail(ctx echo.Context) error {
	id := ctx.Param("id")

	consultation, err := c.ConsultationUsecase.GetConsultationByID(id)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusNotFound, "Consultation not found")
	}

	response := map[string]interface{}{
		"user": map[string]interface{}{
			"username": consultation.User.Username,
			"email":    consultation.User.Email,
			"message":  consultation.Message,
		},
		"start_time":     consultation.StartTime,
		"end_time":       consultation.EndTime,
		"recommendation": consultation.Recommendation,
	}

	return helper.JSONSuccessResponse(ctx, response)
}

// GiveRecommendation - Doctor provides recommendation
func (c *ConsultationController) GiveRecommendation(ctx echo.Context) error {
	claims, _ := ctx.Get("doctor").(*service.JwtCustomClaims)
	doctorID := claims.UserID

	var request struct {
		Recommendation string `json:"recommendation"`
	}

	id := ctx.Param("id")
	if err := ctx.Bind(&request); err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Invalid input")
	}

	consultation, err := c.ConsultationUsecase.GetConsultationByID(id)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusNotFound, "Consultation not found")
	}

	if consultation.DoctorID != doctorID {
		return helper.JSONErrorResponse(ctx, http.StatusForbidden, "Access to this consultation is denied")
	}

	err = c.ConsultationUsecase.UpdateRecommendation(consultation.ID, request.Recommendation)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Failed to save recommendation")
	}

	return helper.JSONSuccessResponse(ctx, "Recommendation successfully provided")
}

// CheckConsultationStatus - Check consultation and payment status
func (c *ConsultationController) CheckConsultationStatus(ctx echo.Context) error {
	consultationID := ctx.Param("id")

	consultation, err := c.ConsultationUsecase.GetConsultationByID(consultationID)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusNotFound, "Konsultasi tidak ditemukan")
	}

	response := map[string]interface{}{
		"consultation_status": consultation.Status,
		"payment_status":      consultation.PaymentStatus,
	}

	return helper.JSONSuccessResponse(ctx, response)
}
