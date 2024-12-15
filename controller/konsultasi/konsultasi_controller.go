package controller

import (
	"calmind/helper"
	"calmind/model"
	"calmind/service"
	usecase "calmind/usecase/konsultasi"

	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type ConsultationController struct {
	ConsultationUsecase *usecase.ConsultationUsecaseImpl
}

func NewConsultationController(consultationUsecase *usecase.ConsultationUsecaseImpl) *ConsultationController {
	return &ConsultationController{ConsultationUsecase: consultationUsecase}
}

//  -- user --

// Melihat semua konsultasi user
func (c *ConsultationController) GetUserConsultations(ctx echo.Context) error {
	claims, ok := ctx.Get("user").(*service.JwtCustomClaims)
	if !ok || claims == nil {
		return helper.JSONErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized access.")
	}

	consultations, err := c.ConsultationUsecase.GetUserConsultations(claims.UserID)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve consultations.")
	}

	var response []model.ConsultationDTO
	for _, cons := range consultations {
		response = append(response, mapConsultationToDTO(cons))
	}

	return helper.JSONSuccessResponse(ctx, response)
}

// Melihat detail konsultasi user
func (c *ConsultationController) GetUserConsultationDetails(ctx echo.Context) error {
	claims, ok := ctx.Get("user").(*service.JwtCustomClaims)
	if !ok || claims == nil {
		return helper.JSONErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized access.")
	}

	consultationID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Invalid consultation ID.")
	}

	consultation, err := c.ConsultationUsecase.GetConsultationByID(consultationID)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusNotFound, "Consultation not found.")
	}

	if consultation.UserID != claims.UserID {
		return helper.JSONErrorResponse(ctx, http.StatusForbidden, "You do not have access to this consultation.")
	}

	return helper.JSONSuccessResponse(ctx, mapConsultationToDTO(*consultation))
}

// Membuat konsultasi baru
func (c *ConsultationController) CreateConsultation(ctx echo.Context) error {
	claims, ok := ctx.Get("user").(*service.JwtCustomClaims)
	if !ok || claims == nil {
		return helper.JSONErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized access.")
	}

	var request struct {
		DoctorID    int    `json:"doctor_id"`
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	if err := ctx.Bind(&request); err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Invalid input.")
	}

	paymentURL, consultation, err := c.ConsultationUsecase.CreateConsultation(claims.UserID, request.DoctorID, request.Title, request.Description, claims.Email)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Failed to create consultation: "+err.Error())
	}

	// Mapping response ke SimpleConsultationDTO
	response := model.SimpleConsultationDTO{
		ID:            consultation.ID,
		Title:         consultation.Title,
		Description:   consultation.Description,
		Duration:      consultation.Duration,
		Status:        consultation.Status,
		PaymentStatus: consultation.PaymentStatus,
		StartTime:     consultation.StartTime.Format(time.RFC3339),
		OrderID:       consultation.OrderID,
		PaymentURL:    paymentURL,
	}

	return helper.JSONSuccessResponse(ctx, response)
}

// Melihat semua konsultasi pasien
func (c *ConsultationController) GetConsultationsForDoctor(ctx echo.Context) error {
	claims, ok := ctx.Get("doctor").(*service.JwtCustomClaims)
	if !ok || claims == nil {
		return helper.JSONErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized access.")
	}

	consultations, err := c.ConsultationUsecase.GetConsultationsForDoctor(claims.UserID)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve consultations.")
	}

	var response []model.ConsultationDTO
	for _, cons := range consultations {
		if cons.Status == "paid" || cons.Status == "approved" { // Tambahkan "approved"
			response = append(response, mapConsultationToDTO(cons))
		}
	}

	return helper.JSONSuccessResponse(ctx, response)
}

// Mendapatkan daftar konsultasi untuk dokter (search by name)
func (c *ConsultationController) SearchConsultationsByName(ctx echo.Context) error {
	claims, ok := ctx.Get("doctor").(*service.JwtCustomClaims)
	if !ok || claims == nil {
		return helper.JSONErrorResponse(ctx, http.StatusUnauthorized, "Doctor is not authorized.")
	}

	doctorID := claims.UserID
	searchName := ctx.QueryParam("nama") // Query param untuk nama user

	// Validasi parameter kosong
	if searchName == "" {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Query parameter 'nama' is required.")
	}

	consultations, err := c.ConsultationUsecase.SearchConsultationsByName(doctorID, searchName)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve consultations.")
	}

	// Konversi ke DTO sebelum dikembalikan
	var response []model.ConsultationDTO
	for _, cons := range consultations {
		response = append(response, mapConsultationToDTO(cons))
	}

	return helper.JSONSuccessResponse(ctx, response)
}

// Melihat detail konsultasi pasien
func (c *ConsultationController) ViewConsultationDetails(ctx echo.Context) error {
	claims, ok := ctx.Get("doctor").(*service.JwtCustomClaims)
	if !ok || claims == nil {
		return helper.JSONErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized access.")
	}

	consultationID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Invalid consultation ID.")
	}

	consultation, err := c.ConsultationUsecase.ViewConsultationDetails(claims.UserID, consultationID)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusNotFound, "Consultation not found.")
	}

	return helper.JSONSuccessResponse(ctx, mapConsultationToDTO(*consultation))
}

// Menambahkan rekomendasi untuk konsultasi
func (c *ConsultationController) AddRecommendation(ctx echo.Context) error {
	claims, ok := ctx.Get("doctor").(*service.JwtCustomClaims)
	if !ok || claims == nil {
		return helper.JSONErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized access.")
	}

	consultationID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Invalid consultation ID.")
	}

	var request struct {
		Recommendation string `json:"recommendation"`
	}

	if err := ctx.Bind(&request); err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Invalid input.")
	}

	err = c.ConsultationUsecase.AddRecommendation(claims.UserID, consultationID, request.Recommendation)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Failed to add recommendation: "+err.Error())
	}

	return helper.JSONSuccessResponse(ctx, map[string]interface{}{
		"message": "Recommendation added successfully.",
	})
}

// Melihat semua konsultasi yang menunggu persetujuan
func (c *ConsultationController) GetPendingConsultations(ctx echo.Context) error {
	claims, ok := ctx.Get("admin").(*service.JwtCustomClaims)
	if !ok || claims == nil {
		return helper.JSONErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized access.")
	}

	consultations, err := c.ConsultationUsecase.GetPendingConsultations()
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve pending consultations.")
	}

	var response []model.ConsultationDTO
	for _, cons := range consultations {
		response = append(response, mapConsultationToDTO(cons))
	}

	return helper.JSONSuccessResponse(ctx, response)
}
func (c *ConsultationController) GetAproveConsultations(ctx echo.Context) error {
	claims, ok := ctx.Get("admin").(*service.JwtCustomClaims)
	if !ok || claims == nil {
		return helper.JSONErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized access.")
	}

	consultations, err := c.ConsultationUsecase.GetApprovedConsultations()
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve pending consultations.")
	}

	var response []model.ConsultationDTO
	for _, cons := range consultations {
		response = append(response, mapConsultationToDTO(cons))
	}

	return helper.JSONSuccessResponse(ctx, response)
}
func (c *ConsultationController) GetAllStatusConsultations(ctx echo.Context) error {
	claims, ok := ctx.Get("admin").(*service.JwtCustomClaims)
	if !ok || claims == nil {
		return helper.JSONErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized access.")
	}

	consultations, err := c.ConsultationUsecase.GetAllStatusConsultations()
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve pending consultations.")
	}

	var response []model.ConsultationDTO
	for _, cons := range consultations {
		response = append(response, mapConsultationToDTO(cons))
	}

	return helper.JSONSuccessResponse(ctx, response)
}

// Melihat detail konsultasi untuk persetujuan
func (c *ConsultationController) ViewPendingConsultation(ctx echo.Context) error {
	claims, ok := ctx.Get("admin").(*service.JwtCustomClaims)
	if !ok || claims == nil {
		return helper.JSONErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized access.")
	}

	consultationID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Invalid consultation ID.")
	}

	consultation, err := c.ConsultationUsecase.GetConsultationByID(consultationID)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusNotFound, "Consultation not found.")
	}

	return helper.JSONSuccessResponse(ctx, mapConsultationToDTO(*consultation))
}
func (c *ConsultationController) ApprovePaymentAndConsultation(ctx echo.Context) error {
	claims, ok := ctx.Get("admin").(*service.JwtCustomClaims)
	if !ok || claims == nil {
		return helper.JSONErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized access.")
	}

	consultationID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Invalid consultation ID.")
	}

	var request struct {
		Status string `json:"status"` // Status pembayaran
	}

	if err := ctx.Bind(&request); err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusBadRequest, "Invalid input.")
	}

	// Panggil usecase untuk menyetujui pembayaran dan konsultasi
	err = c.ConsultationUsecase.ApprovePaymentAndConsultation(consultationID, request.Status)
	if err != nil {
		return helper.JSONErrorResponse(ctx, http.StatusInternalServerError, "Failed to approve consultation: "+err.Error())
	}

	return helper.JSONSuccessResponse(ctx, map[string]interface{}{
		"message": "Consultation and payment approved successfully.",
	})
}

// -- dto respon --

func mapConsultationToDTO(consultation model.Consultation) model.ConsultationDTO {
	return model.ConsultationDTO{
		ID:            consultation.ID,
		Title:         consultation.Title,
		Description:   consultation.Description,
		Duration:      consultation.Duration,
		Status:        consultation.Status,
		PaymentStatus: consultation.PaymentStatus,
		StartTime:     consultation.StartTime.Format(time.RFC3339),
		OrderID:       consultation.OrderID, // Include OrderID
		CreatedAt:     consultation.CreatedAt,
		UpdatedAt:     consultation.UpdatedAt,
		User: &model.UserDTO{
			Avatar:       consultation.User.Avatar,
			Username:     consultation.User.Username,
			Email:        consultation.User.Email,
			Pekerjaan:    consultation.User.Pekerjaan, // Pastikan ini benar
			JenisKelamin: consultation.User.JenisKelamin,
			NoHp:         consultation.User.NoHp,
			TglLahir:     consultation.User.TglLahir,
		},
		Doctor: &model.DoctorDTO{
			Avatar:   consultation.User.Avatar,
			Username: consultation.Doctor.Username,
			Email:    consultation.Doctor.Email,
			Title:    consultation.Doctor.Title,
			About:    consultation.Doctor.About, // Pastikan ini benar
		},
		Rekomendasi: mapRecommendationsToDTO(consultation.Rekomendasi),
	}
}

func mapRecommendationsToDTO(recommendations []model.Rekomendasi) []model.RecommendationDTO {
	var result []model.RecommendationDTO
	for _, r := range recommendations {
		result = append(result, model.RecommendationDTO{
			ID:             r.ID,
			ConsultationID: r.ConsultationID,
			DoctorID:       r.DoctorID,
			Recommendation: r.Rekomendasi,
		})
	}
	return result
}
