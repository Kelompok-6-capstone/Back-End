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
		return errors.New("consultation not found")
	}

	if consultation.IsApproved {
		return errors.New("payment already approved")
	}

	if consultation.Doctor.Price <= 0 {
		return errors.New("invalid doctor price")
	}

	// Tentukan pembagian biaya
	adminFee := 0.1 * float64(consultation.Doctor.Price) // 10% biaya admin
	doctorFee := float64(consultation.Doctor.Price) - adminFee

	// Update konsultasi
	consultation.IsApproved = true
	consultation.Status = "active"
	consultation.StartTime = time.Now()

	err = uc.Repo.UpdateConsultation(consultation)
	if err != nil {
		return errors.New("failed to update consultation status")
	}

	// Tambahkan pendapatan dokter dan admin
	err = uc.Repo.AddIncomeForDoctor(consultation.DoctorID, doctorFee)
	if err != nil {
		return errors.New("failed to update doctor income")
	}

	err = uc.Repo.AddIncomeForAdmin(adminFee)
	if err != nil {
		return errors.New("failed to update admin income")
	}

	return nil
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

// Mendapatkan daftar konsultasi yang belum disetujui pembayarannya
func (uc *ConsultationUsecaseImpl) GetPendingConsultationsForAdmin() ([]model.Consultation, error) {
	return uc.Repo.GetPendingConsultations()
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
	client.New(os.Getenv("MIDTRANS_SERVER_KEY"), midtrans.Sandbox) // Gunakan midtrans.Production jika live

	// Gunakan orderID untuk mencari status transaksi
	orderID := fmt.Sprintf("consultation-%d", consultationID)

	// Panggil API untuk mendapatkan status transaksi
	transactionStatusResp, err := client.CheckTransaction(orderID)
	if err != nil {
		return "", fmt.Errorf("failed to verify payment: %w", err)
	}

	return transactionStatusResp.TransactionStatus, nil
}
