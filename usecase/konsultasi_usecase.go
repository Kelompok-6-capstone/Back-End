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
	Repo *repository.ConsultationRepositoryImpl
}

func NewConsultationUsecaseImpl(repo *repository.ConsultationRepositoryImpl) *ConsultationUsecaseImpl {
	return &ConsultationUsecaseImpl{Repo: repo}
}

func (uc *ConsultationUsecaseImpl) MarkExpiredConsultations() error {
	consultations, err := uc.Repo.GetActiveConsultations()
	if err != nil {
		return fmt.Errorf("failed to fetch active consultations: %v", err)
	}

	now := time.Now()
	for _, consultation := range consultations {
		endTime := consultation.StartTime.Add(time.Duration(consultation.Duration) * time.Minute)
		if now.After(endTime) {
			consultation.Status = "expired"
			err := uc.Repo.UpdateConsultation(&consultation)
			if err != nil {
				return fmt.Errorf("failed to update consultation %d: %v", consultation.ID, err)
			}
		}
	}
	return nil
}

func (uc *ConsultationUsecaseImpl) ApprovePaymentAndConsultation(consultationID int, paymentStatus string) error {
	// Validasi pembayaran
	if paymentStatus != "paid" {
		return errors.New("invalid payment status, only 'paid' is allowed")
	}

	// Dapatkan data konsultasi berdasarkan ID
	consultation, err := uc.Repo.GetConsultationByID(consultationID)
	if err != nil {
		return fmt.Errorf("consultation not found: %w", err)
	}

	// Perbarui status konsultasi dan pembayaran
	consultation.Status = "approved"
	consultation.PaymentStatus = "paid" // Pastikan ini juga diperbarui

	// Simpan perubahan ke database
	err = uc.Repo.UpdateConsultation(consultation)
	if err != nil {
		return fmt.Errorf("failed to update consultation: %w", err)
	}

	return nil
}

// Membuat konsultasi baru
func (uc *ConsultationUsecaseImpl) CreateConsultation(userID, doctorID int, title, description, email string) (string, *model.Consultation, error) {
	// Validasi user dan dokter
	err := uc.Repo.ValidateUserAndDoctor(userID, doctorID)
	if err != nil {
		return "", nil, err
	}

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

// Mendapatkan konsultasi berdasarkan nama user
func (u *ConsultationUsecaseImpl) SearchConsultationsByName(doctorID int, searchName string) (*[]model.Consultation, error) {
	var consultations *[]model.Consultation
	err := u.Repo.FindConsultationsByDoctorAndName(doctorID, searchName, consultations)
	if err != nil {
		return nil, err
	}
	return consultations, nil
}

// Mendapatkan konsultasi berdasarkan ID
func (uc *ConsultationUsecaseImpl) GetConsultationByID(consultationID int) (*model.Consultation, error) {
	return uc.Repo.GetConsultationByID(consultationID)
}

// Mendapatkan daftar konsultasi untuk dokter
func (uc *ConsultationUsecaseImpl) GetConsultationsForDoctor(doctorID int) ([]model.Consultation, error) {
	consultations, err := uc.Repo.GetConsultationsForDoctor(doctorID)
	if err != nil {
		return nil, err
	}
	if len(consultations) == 0 {
		fmt.Printf("No consultations found for doctor ID %d\n", doctorID) // Debugging tambahan
	}
	return consultations, nil
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
func (uc *ConsultationUsecaseImpl) GetApprovedConsultations() ([]model.Consultation, error) {
	return uc.Repo.GetApprovedConsultations()
}
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
