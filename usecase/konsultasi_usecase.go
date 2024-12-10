package usecase

import (
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
	GetConsultationsForDoctor(doctorID int) ([]model.Consultation, error)
	ViewConsultationDetails(doctorID, consultationID int) (*model.Consultation, error)
	AddRecommendation(doctorID, consultationID int, recommendation string) error
	GetUserConsultations(userID int) ([]model.Consultation, error)
	GetPendingConsultations() ([]model.Consultation, error)
	ApproveConsultation(consultationID int) error
	CreateMidtransPayment(consultationID int, amount float64, email string) (string, error)
	VerifyPayment(consultationID int) (string, error)
	MarkExpiredConsultations() error
}

type ConsultationUsecaseImpl struct {
	Repo *repository.ConsultationRepositoryImpl
}

func NewConsultationUsecaseImpl(repo *repository.ConsultationRepositoryImpl) *ConsultationUsecaseImpl {
	return &ConsultationUsecaseImpl{Repo: repo}
}

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

// Membuat konsultasi baru
func (uc *ConsultationUsecaseImpl) CreateConsultation(userID, doctorID int, title, description, email string) (string, *model.Consultation, error) {
	doctor, err := uc.Repo.GetDoctorByID(doctorID)
	if err != nil {
		return "", nil, fmt.Errorf("doctor not found: %w", err)
	}

	if doctor.Price <= 0 {
		return "", nil, errors.New("doctor price is invalid")
	}

	consultation := &model.Consultation{
		UserID:      userID,
		DoctorID:    doctorID,
		Title:       title,
		Description: description,
		Status:      "pending",
		StartTime:   time.Now(),
	}

	consultationID, err := uc.Repo.CreateConsultation(consultation)
	if err != nil {
		return "", nil, errors.New("failed to create consultation")
	}

	paymentURL, err := uc.CreateMidtransPayment(consultationID, doctor.Price, email)
	if err != nil {
		return "", nil, fmt.Errorf("failed to create midtrans payment: %w", err)
	}

	return paymentURL, consultation, nil
}

// Mendapatkan konsultasi berdasarkan ID
func (uc *ConsultationUsecaseImpl) GetConsultationByID(consultationID int) (*model.Consultation, error) {
	return uc.Repo.GetConsultationByID(consultationID)
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

// Mendapatkan daftar konsultasi yang menunggu persetujuan admin
func (uc *ConsultationUsecaseImpl) GetPendingConsultations() ([]model.Consultation, error) {
	return uc.Repo.GetPendingConsultations()
}

// Menyetujui konsultasi
func (uc *ConsultationUsecaseImpl) ApproveConsultation(consultationID int) error {
	consultation, err := uc.Repo.GetConsultationByID(consultationID)
	if err != nil {
		return fmt.Errorf("consultation not found: %w", err)
	}

	if consultation.Status != "paid" {
		return errors.New("consultation is not paid yet")
	}

	consultation.Status = "approved"
	return uc.Repo.UpdateConsultation(consultation)
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

// Verifikasi pembayaran menggunakan Midtrans
func (uc *ConsultationUsecaseImpl) VerifyPayment(consultationID int) (string, error) {
	client := coreapi.Client{}
	client.New(os.Getenv("MIDTRANS_SERVER_KEY"), midtrans.Sandbox)

	orderID := fmt.Sprintf("consultation-%d", consultationID)

	transactionStatusResp, err := client.CheckTransaction(orderID)
	if err != nil {
		return "", fmt.Errorf("failed to verify payment: %w", err)
	}

	if transactionStatusResp.TransactionStatus == "settlement" {
		consultation, err := uc.Repo.GetConsultationByID(consultationID)
		if err != nil {
			return "", fmt.Errorf("consultation not found: %w", err)
		}

		if consultation.Status != "paid" {
			consultation.Status = "paid"
			err = uc.Repo.UpdateConsultation(consultation)
			if err != nil {
				return "", fmt.Errorf("failed to update consultation status: %w", err)
			}
		}
	}

	return transactionStatusResp.TransactionStatus, nil
}
