package usecase

import (
	"calmind/helper"
	"calmind/model"
	"calmind/repository"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
)

type ConsultationUsecase interface {
	CreateConsultation(userID, doctorID int, title, description, email string) (string, *model.Consultation, error)
	GetConsultationByID(consultationID int) (*model.Consultation, error)
	MarkExpiredConsultations() error
	ApprovePayment(adminID, consultationID int) error
	ViewConsultationDetailsForAdmin(consultationID int) (*model.Consultation, error)
	GetConsultationsForDoctor(doctorID int) ([]model.Consultation, error)
	ViewConsultationDetails(doctorID, consultationID int) (*model.Consultation, error)
	AddRecommendation(doctorID, consultationID int, recommendation string) error
	GetUserConsultations(userID int) ([]model.Consultation, error)
	GetPendingConsultationsForAdmin() ([]model.Consultation, error)
	CreateMidtransPayment(consultationID int, amount float64, email string) (string, error)
	GetPendingPaymentsForAdmin() ([]model.Consultation, error)
	GetPaymentDetails(consultationID int) (map[string]interface{}, error)
	GetPendingConsultations() ([]model.Consultation, error)
	GetDoctorByID(doctorID int) (*model.Doctor, error)
	VerifyPaymentStatus(orderID string) (*coreapi.TransactionStatusResponse, error)
	UpdatePaymentStatus(orderID string, status string) error
}

type ConsultationUsecaseImpl struct {
	Repo *repository.ConsultationRepositoryImpl
}

func NewConsultationUsecaseImpl(repo *repository.ConsultationRepositoryImpl) *ConsultationUsecaseImpl {
	return &ConsultationUsecaseImpl{Repo: repo}
}

// Membuat konsultasi
func (uc *ConsultationUsecaseImpl) CreateConsultation(userID, doctorID int, title, description, email string) (string, *model.Consultation, error) {
	// Fetch doctor details from repository
	doctor, err := uc.Repo.GetDoctorByID(doctorID)
	if err != nil {
		return "", nil, fmt.Errorf("doctor not found: %w", err)
	}

	// Validate doctor's price
	if doctor.Price <= 0 {
		return "", nil, errors.New("doctor price is invalid")
	}

	// Create consultation object
	consultation := &model.Consultation{
		UserID:        userID,
		DoctorID:      doctorID,
		Title:         title,
		Description:   description,
		Status:        "pending",
		PaymentStatus: "pending", // Initialize payment status as pending
		IsApproved:    false,
		StartTime:     time.Now(),
	}

	// Save consultation to the database
	consultationID, err := uc.Repo.CreateConsultation(consultation)
	if err != nil {
		return "", nil, fmt.Errorf("failed to create consultation: %w", err)
	}

	// Generate payment URL using Midtrans
	paymentURL, err := uc.CreateMidtransPayment(consultationID, doctor.Price, email)
	if err != nil {
		return "", nil, fmt.Errorf("failed to create Midtrans payment: %w", err)
	}

	return paymentURL, consultation, nil
}

// Mendapatkan konsultasi berdasarkan ID
func (uc *ConsultationUsecaseImpl) GetConsultationByID(consultationID int) (*model.Consultation, error) {
	return uc.Repo.GetConsultationByID(consultationID)
}

// Menandai konsultasi yang sudah kedaluwarsa
func (uc *ConsultationUsecaseImpl) MarkExpiredConsultations() error {
	consultations, err := uc.Repo.GetActiveConsultations()
	if err != nil {
		return err
	}

	now := time.Now()
	for _, consultation := range consultations {
		endTime := consultation.StartTime.Add(time.Duration(consultation.Duration) * time.Minute)
		if now.After(endTime) {
			consultation.Status = "ended"
			if err := uc.Repo.UpdateConsultation(&consultation); err != nil {
				return errors.New("failed to update consultation status")
			}
		}
	}
	return nil
}
func (uc *ConsultationUsecaseImpl) ApprovePayment(adminID, consultationID int) error {
	consultation, err := uc.Repo.GetConsultationByID(consultationID)
	if err != nil {
		return fmt.Errorf("consultation not found: %w", err)
	}

	if consultation.PaymentStatus != "settlement" {
		return errors.New("payment not completed")
	}

	if consultation.IsApproved {
		return errors.New("payment already approved")
	}

	consultation.IsApproved = true
	consultation.Status = "active"
	consultation.StartTime = time.Now()

	return uc.Repo.UpdateConsultation(consultation)
}

// Mendapatkan detail konsultasi untuk admin
func (uc *ConsultationUsecaseImpl) ViewConsultationDetailsForAdmin(consultationID int) (*model.Consultation, error) {
	return uc.Repo.GetAdminViewConsultation(consultationID)
}

// Mendapatkan daftar konsultasi untuk dokter
func (uc *ConsultationUsecaseImpl) GetConsultationsForDoctor(doctorID int) ([]model.Consultation, error) {
	return uc.Repo.GetConsultationsForDoctor(doctorID)
}

// Mendapatkan detail konsultasi tertentu untuk dokter
func (uc *ConsultationUsecaseImpl) ViewConsultationDetails(doctorID, consultationID int) (*model.Consultation, error) {
	return uc.Repo.GetConsultationDetails(consultationID, doctorID)
}

// Menambahkan rekomendasi untuk konsultasi
func (uc *ConsultationUsecaseImpl) AddRecommendation(doctorID, consultationID int, recommendation string) error {
	recommendationObj := &model.Rekomendasi{
		ConsultationID: consultationID,
		DoctorID:       doctorID,
		Rekomendasi:    recommendation,
	}
	return uc.Repo.AddRecommendation(recommendationObj)
}

// Mendapatkan daftar konsultasi untuk user tertentu
func (uc *ConsultationUsecaseImpl) GetUserConsultations(userID int) ([]model.Consultation, error) {
	return uc.Repo.GetConsultationsWithDoctors(userID)
}

// Mendapatkan daftar konsultasi pending
func (uc *ConsultationUsecaseImpl) GetPendingConsultations() ([]model.Consultation, error) {
	return uc.Repo.GetPendingConsultations()
}

// Menambahkan income untuk dokter
func (uc *ConsultationUsecaseImpl) AddIncomeForDoctor(doctorID int, amount float64) error {
	return uc.Repo.AddIncomeForDoctor(doctorID, amount)
}

// Menambahkan income untuk admin
func (uc *ConsultationUsecaseImpl) AddIncomeForAdmin(amount float64) error {
	return uc.Repo.AddIncomeForAdmin(amount)
}

// Mendapatkan dokter berdasarkan ID
func (uc *ConsultationUsecaseImpl) GetDoctorByID(doctorID int) (*model.Doctor, error) {
	return uc.Repo.GetDoctorByID(doctorID)
}

// Mendapatkan pembayaran pending untuk admin
func (uc *ConsultationUsecaseImpl) GetPendingPaymentsForAdmin() ([]model.Consultation, error) {
	return uc.Repo.GetPendingPayments()
}

// Membuat pembayaran menggunakan Midtrans
func (uc *ConsultationUsecaseImpl) CreateMidtransPayment(consultationID int, amount float64, email string) (string, error) {
	client := snap.Client{}
	client.New(os.Getenv("MIDTRANS_SERVER_KEY"), midtrans.Sandbox)

	snapReq := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  fmt.Sprintf("consultation-%d-%d", consultationID, time.Now().Unix()),
			GrossAmt: int64(amount),
		},
		CustomerDetail: &midtrans.CustomerDetails{
			Email: email,
		},
	}

	snapTokenResp, err := client.CreateTransaction(snapReq)
	if err != nil {
		return "", fmt.Errorf("failed to create midtrans payment: %w", err)
	}

	return snapTokenResp.RedirectURL, nil
}

// Mendapatkan detail pembayaran
func (uc *ConsultationUsecaseImpl) GetPaymentDetails(consultationID int) (map[string]interface{}, error) {
	consultation, err := uc.Repo.GetConsultationByID(consultationID)
	if err != nil {
		return nil, fmt.Errorf("consultation not found: %w", err)
	}

	return map[string]interface{}{
		"consultation": map[string]interface{}{
			"id":          consultation.ID,
			"title":       consultation.Title,
			"description": consultation.Description,
			"status":      consultation.Status,
		},
		"payment": map[string]interface{}{
			"payment_status": consultation.PaymentStatus,
		},
	}, nil
}
func (uc *ConsultationUsecaseImpl) UpdatePaymentStatus(orderID string, status string) error {
	consultationID, err := helper.ParseOrderID(orderID)
	if err != nil {
		return fmt.Errorf("failed to parse order ID: %w", err)
	}

	consultation, err := uc.Repo.GetConsultationByID(consultationID)
	if err != nil {
		return fmt.Errorf("consultation not found: %w", err)
	}

	consultation.PaymentStatus = status
	if status == "settlement" {
		consultation.Status = "active"
		consultation.IsApproved = true
	} else if status == "cancel" || status == "expire" {
		consultation.Status = "failed"
	} else {
		consultation.Status = "pending"
	}

	return uc.Repo.UpdateConsultation(consultation)
}

func (uc *ConsultationUsecaseImpl) VerifyPaymentStatus(orderID string) (*coreapi.TransactionStatusResponse, error) {
	client := coreapi.Client{}
	client.New(os.Getenv("MIDTRANS_SERVER_KEY"), midtrans.Sandbox)

	transactionStatusResp, err := client.CheckTransaction(orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to verify payment status: %w", err)
	}

	return transactionStatusResp, nil
}
