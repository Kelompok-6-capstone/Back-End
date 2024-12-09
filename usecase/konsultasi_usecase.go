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
	CreateConsultation(userID, doctorID int, title, description, email string) (string, error)
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
	VerifyPayment(consultationID int) (string, error)
	GetPendingPaymentsForAdmin() ([]model.Consultation, error)
	GetPaymentDetails(consultationID int) (map[string]interface{}, error)
}

type ConsultationUsecaseImpl struct {
	Repo *repository.ConsultationRepositoryImpl
}

func NewConsultationUsecaseImpl(repo *repository.ConsultationRepositoryImpl) *ConsultationUsecaseImpl {
	return &ConsultationUsecaseImpl{Repo: repo}
}

// Membuat konsultasi
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
		IsPaid:      false,
		IsApproved:  false,
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
	consultation, err := uc.Repo.GetConsultationByID(consultationID)
	if err != nil {
		return nil, err
	}

	// Periksa apakah konsultasi telah kedaluwarsa
	now := time.Now()
	endTime := consultation.StartTime.Add(time.Duration(consultation.Duration) * time.Minute)
	if now.After(endTime) && consultation.Status != "ended" {
		consultation.Status = "ended"
		err = uc.Repo.UpdateConsultation(consultation)
		if err != nil {
			return nil, errors.New("failed to update consultation status")
		}
	}

	return consultation, nil
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
func (uc *ConsultationUsecaseImpl) ApprovePayment(consultationID int) error {
	consultation, err := uc.Repo.GetConsultationByID(consultationID)
	if err != nil {
		return fmt.Errorf("consultation not found: %w", err)
	}

	// Validasi apakah pembayaran telah selesai
	if !consultation.IsPaid {
		return errors.New("payment not completed")
	}

	// Cegah persetujuan duplikat
	if consultation.IsApproved {
		return errors.New("payment already approved")
	}

	// Perbarui status konsultasi
	consultation.IsApproved = true
	consultation.Status = "active"
	consultation.StartTime = time.Now()

	return uc.Repo.UpdateConsultation(consultation)
}

// Mendapatkan detail konsultasi untuk admin
func (uc *ConsultationUsecaseImpl) ViewConsultationDetailsForAdmin(consultationID int) (*model.Consultation, error) {
	return uc.Repo.GetConsultationByID(consultationID)
}

// Mendapatkan daftar konsultasi untuk dokter
func (uc *ConsultationUsecaseImpl) GetConsultationsForDoctor(doctorID int) ([]model.Consultation, error) {
	return uc.Repo.GetConsultationsForDoctor(doctorID)
}

// Mendapatkan detail konsultasi tertentu untuk dokter
func (uc *ConsultationUsecaseImpl) ViewConsultationDetails(doctorID, consultationID int) (*model.Consultation, error) {
	consultation, err := uc.Repo.GetConsultationDetails(consultationID, doctorID)
	if err != nil {
		return nil, err
	}

	if consultation.DoctorID != doctorID {
		return nil, errors.New("unauthorized access to consultation")
	}

	return consultation, nil
}

// Menambahkan rekomendasi untuk konsultasi
func (uc *ConsultationUsecaseImpl) AddRecommendation(doctorID, consultationID int, recommendation string) error {
	consultation, err := uc.Repo.GetConsultationDetails(consultationID, doctorID)
	if err != nil {
		return err
	}

	if consultation.DoctorID != doctorID {
		return errors.New("unauthorized access to consultation")
	}

	recommendationObj := &model.Rekomendasi{
		ConsultationID: consultation.ID,
		DoctorID:       doctorID,
		Rekomendasi:    recommendation,
	}
	return uc.Repo.AddRecommendation(recommendationObj)
}

// Mendapatkan daftar konsultasi dengan dokter untuk user tertentu
func (uc *ConsultationUsecaseImpl) GetUserConsultations(userID int) ([]model.Consultation, error) {
	return uc.Repo.GetConsultationsWithDoctors(userID)
}

func (uc *ConsultationUsecaseImpl) CreateMidtransPayment(consultationID int, amount float64, email string) (string, error) {
	// Konfigurasi Midtrans
	client := snap.Client{}
	client.New(os.Getenv("MIDTRANS_SERVER_KEY"), midtrans.Sandbox)

	// Buat detail transaksi
	snapReq := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  fmt.Sprintf("consultation-%d", consultationID),
			GrossAmt: int64(amount), // Gunakan harga dokter
		},
		CustomerDetail: &midtrans.CustomerDetails{
			Email: email,
		},
	}

	// Generate Snap Token
	snapTokenResp, err := client.CreateTransaction(snapReq)
	if err != nil {
		return "", fmt.Errorf("failed to create midtrans payment: %w", err)
	}

	return snapTokenResp.RedirectURL, nil
}

func (uc *ConsultationUsecaseImpl) VerifyPayment(consultationID int) (string, error) {
	client := coreapi.Client{}
	client.New(os.Getenv("MIDTRANS_SERVER_KEY"), midtrans.Sandbox)

	// Gunakan orderID untuk mencari status transaksi
	orderID := fmt.Sprintf("consultation-%d", consultationID)

	// Panggil API untuk mendapatkan status transaksi
	transactionStatusResp, err := client.CheckTransaction(orderID)
	if err != nil {
		return "", fmt.Errorf("failed to verify payment: %w", err)
	}

	// Jika pembayaran berhasil, update is_paid
	if transactionStatusResp.TransactionStatus == "settlement" {
		consultation, err := uc.Repo.GetConsultationByID(consultationID)
		if err != nil {
			return "", fmt.Errorf("failed to find consultation: %w", err)
		}

		// Update is_paid tanpa mengubah status lainnya
		consultation.IsPaid = true
		err = uc.Repo.UpdateConsultation(consultation)
		if err != nil {
			return "", fmt.Errorf("failed to update consultation payment status: %w", err)
		}
	}

	return transactionStatusResp.TransactionStatus, nil
}

func (uc *ConsultationUsecaseImpl) GetPendingPaymentsForAdmin() ([]model.Consultation, error) {
	return uc.Repo.GetPendingPayments()
}
func (uc *ConsultationUsecaseImpl) GetPaymentDetails(consultationID int) (map[string]interface{}, error) {
	// Ambil detail konsultasi dari database
	consultation, err := uc.Repo.GetConsultationByID(consultationID)
	if err != nil {
		return nil, fmt.Errorf("consultation not found: %w", err)
	}

	// Ambil detail pembayaran dari Midtrans
	client := coreapi.Client{}
	client.New(os.Getenv("MIDTRANS_SERVER_KEY"), midtrans.Sandbox)

	orderID := fmt.Sprintf("consultation-%d", consultationID)
	transactionStatusResp, err := client.CheckTransaction(orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch payment details: %w", err)
	}

	// Kombinasikan data konsultasi dan detail pembayaran
	return map[string]interface{}{
		"consultation": consultation,
		"payment": map[string]interface{}{
			"order_id":           transactionStatusResp.OrderID,
			"payment_type":       transactionStatusResp.PaymentType,
			"transaction_time":   transactionStatusResp.TransactionTime,
			"transaction_status": transactionStatusResp.TransactionStatus,
			"gross_amount":       transactionStatusResp.GrossAmount,
		},
	}, nil
}
