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
	SearchConsultationsByName(doctorID int, searchName string) ([]model.Consultation, error)
	GetConsultationByID(consultationID int) (*model.Consultation, error)
	GetConsultationsForDoctor(doctorID int) ([]model.Consultation, error)
	ViewConsultationDetails(doctorID, consultationID int) (*model.Consultation, error)
	AddRecommendation(doctorID, consultationID int, recommendation string) error
	GetUserConsultations(userID int) ([]model.Consultation, error)
	GetPendingConsultations() ([]model.Consultation, error)
	CreateMidtransPayment(consultationID int, amount float64, email string) (string, error)
	VerifyPayment(consultationID int) (string, error)
	MarkExpiredConsultations() error
	ApprovePaymentAndConsultation(consultationID int, paymentStatus string) error
	GetApprovedConsultations() ([]model.Consultation, error)
	GetAllStatusConsultations() ([]model.Consultation, error)
}

type ConsultationUsecaseImpl struct {
	Repo repository.ConsultationRepository
}

func NewConsultationUsecaseImpl(repo repository.ConsultationRepository) *ConsultationUsecaseImpl {
	return &ConsultationUsecaseImpl{Repo: repo}
}

// Menandai konsultasi yang sudah expired
func (uc *ConsultationUsecaseImpl) MarkExpiredConsultations() error {
	// Mengambil semua konsultasi aktif tanpa memfilter berdasarkan userID dan doctorID
	consultations, err := uc.Repo.GetAllActiveConsultations()
	if err != nil {
		return fmt.Errorf("failed to fetch active consultations: %v", err)
	}

	now := time.Now()
	for _, consultation := range consultations {
		endTime := consultation.StartTime.Add(time.Duration(consultation.Duration) * time.Minute)
		if now.After(endTime) {
			consultation.Status = "expired"
			// Perbarui status konsultasi ke database
			if err := uc.Repo.UpdateConsultation(&consultation); err != nil {
				return fmt.Errorf("failed to update consultation %d: %v", consultation.ID, err)
			}
		}
	}
	return nil
}

// Menyetujui pembayaran dan konsultasi
func (uc *ConsultationUsecaseImpl) ApprovePaymentAndConsultation(consultationID int, paymentStatus string) error {
	if paymentStatus != "paid" {
		return errors.New("invalid payment status, only 'paid' is allowed")
	}

	consultation, err := uc.Repo.GetConsultationByID(consultationID)
	if err != nil {
		return fmt.Errorf("consultation not found: %w", err)
	}

	consultation.Status = "approved"
	consultation.PaymentStatus = "paid"

	if err := uc.Repo.UpdateConsultation(consultation); err != nil {
		return fmt.Errorf("failed to update consultation: %w", err)
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

	orderID := fmt.Sprintf("CONSULT-%d-%d", userID, time.Now().Unix())

	consultation := &model.Consultation{
		UserID:      userID,
		DoctorID:    doctorID,
		Title:       title,
		Description: description,
		Status:      "pending",
		OrderID:     orderID,
		StartTime:   time.Now(),
	}

	consultationID, err := uc.Repo.CreateConsultation(consultation)
	if err != nil {
		return "", nil, fmt.Errorf("failed to create consultation: %w", err)
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

// Mendapatkan detail konsultasi tertentu
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

// Mendapatkan konsultasi untuk user tertentu
func (uc *ConsultationUsecaseImpl) GetUserConsultations(userID int) ([]model.Consultation, error) {
	return uc.Repo.GetConsultationsWithDoctors(userID)
}

// Mendapatkan konsultasi yang menunggu persetujuan
func (uc *ConsultationUsecaseImpl) GetPendingConsultations() ([]model.Consultation, error) {
	return uc.Repo.GetPendingConsultations()
}

// Mendapatkan konsultasi yang sudah disetujui
func (uc *ConsultationUsecaseImpl) GetApprovedConsultations() ([]model.Consultation, error) {
	return uc.Repo.GetApprovedConsultations()
}

// Mendapatkan konsultasi dengan semua status
func (uc *ConsultationUsecaseImpl) GetAllStatusConsultations() ([]model.Consultation, error) {
	return uc.Repo.GetAllStatusConsultations()
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

	return transactionStatusResp.TransactionStatus, nil
}
