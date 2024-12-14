package usecase

import (
	"calmind/model"
	"calmind/repository"
	"errors"
	"fmt"
	"log"
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
	UpdatePaymentStatus(orderID, transactionStatus string) error
}

type ConsultationUsecaseImpl struct {
	Repo *repository.ConsultationRepositoryImpl
}

func NewConsultationUsecaseImpl(repo *repository.ConsultationRepositoryImpl) *ConsultationUsecaseImpl {
	return &ConsultationUsecaseImpl{Repo: repo}
}

func (uc *ConsultationUsecaseImpl) MarkExpiredConsultations() error {
	consultations, err := uc.Repo.GetActiveConsultations() // Ambil semua konsultasi yang masih aktif
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
				log.Printf("Failed to update consultation %d to expired: %v", consultation.ID, err)
			} else {
				log.Printf("Consultation %d marked as expired", consultation.ID)
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
	consultation.IsApproved = true
	consultation.StartTime = time.Now()
	// Simpan perubahan ke database
	err = uc.Repo.UpdateConsultation(consultation)
	if err != nil {
		return fmt.Errorf("failed to update consultation: %w", err)
	}

	go uc.StartTimerForConsultation(consultation)
	return nil
}

func (uc *ConsultationUsecaseImpl) StartTimerForConsultation(consultation *model.Consultation) {
	duration := time.Duration(consultation.Duration) * time.Minute
	time.Sleep(duration) // Tunggu hingga durasi selesai

	// Setelah timer selesai, perbarui status konsultasi menjadi expired
	consultation.Status = "expired"
	err := uc.Repo.UpdateConsultation(consultation)
	if err != nil {
		log.Printf("Failed to update consultation %d to expired: %v", consultation.ID, err)
	} else {
		log.Printf("Consultation %d has expired", consultation.ID)
	}
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

	// Buat konsultasi tanpa order_id
	consultation := &model.Consultation{
		UserID:      userID,
		DoctorID:    doctorID,
		Title:       title,
		Description: description,
		Status:      "pending",
		StartTime:   time.Now(),
	}

	// Simpan ke database untuk mendapatkan consultationID
	consultationID, err := uc.Repo.CreateConsultation(consultation)
	if err != nil {
		return "", nil, fmt.Errorf("failed to create consultation: %w", err)
	}

	// Buat order_id berdasarkan consultationID
	consultation.OrderID = fmt.Sprintf("consultation-%d", consultationID)

	// Perbarui database dengan order_id
	err = uc.Repo.UpdateConsultation(consultation)
	if err != nil {
		return "", nil, fmt.Errorf("failed to update order_id: %w", err)
	}

	// Buat URL pembayaran menggunakan Midtrans
	paymentURL, err := uc.CreateMidtransPayment(consultation.OrderID, doctor.Price, email)
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
func (uc *ConsultationUsecaseImpl) UpdatePaymentStatus(orderID, transactionStatus string) error {
	// Ambil data konsultasi berdasarkan orderID
	consultation, err := uc.Repo.GetConsultationByOrderID(orderID)
	if err != nil {
		log.Printf("Consultation not found for order_id=%s: %v", orderID, err)
		return fmt.Errorf("consultation not found: %w", err)
	}

	// Validasi status lama agar tidak mengubah status paid menjadi pending atau failed
	if consultation.PaymentStatus == "paid" {
		log.Printf("Payment already completed for order_id=%s", orderID)
		return nil
	}

	// Update payment status sesuai status transaksi
	switch transactionStatus {
	case "settlement":
		consultation.PaymentStatus = "paid"
	case "pending":
		consultation.PaymentStatus = "pending"
	case "cancel", "deny", "expire":
		consultation.PaymentStatus = "failed"
	default:
		log.Printf("Unknown transaction status: %s", transactionStatus)
		return fmt.Errorf("unknown transaction status: %s", transactionStatus)
	}

	// Simpan perubahan ke database
	if err := uc.Repo.UpdateConsultation(consultation); err != nil {
		log.Printf("Failed to update consultation for order_id=%s: %v", orderID, err)
		return fmt.Errorf("failed to update consultation: %w", err)
	}

	log.Printf("Payment status updated successfully for order_id=%s", orderID)
	return nil
}

// Membuat pembayaran menggunakan Midtrans
func (uc *ConsultationUsecaseImpl) CreateMidtransPayment(orderID string, amount float64, email string) (string, error) {
	client := snap.Client{}
	client.New(os.Getenv("MIDTRANS_SERVER_KEY"), midtrans.Sandbox)

	snapReq := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderID, // Gunakan orderID langsung
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
