package usecase

import (
	"calmind/helper"
	"calmind/model"
	repository "calmind/repository/authentikasi"
	"calmind/service"
	"errors"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type DoctorUsecase interface {
	Register(*model.Doctor) error
	Login(email string, password string) (string, error)
	VerifyOtp(email string, code string) error
	ResendOtp(email string) error
}

type doctorUsecase struct {
	DoctorRepo repository.DoctorRepository
	JWTService service.JWTService
	OtpRepo    repository.OtpRepository
	OtpService service.OtpService
}

// NewDoctorAuthUsecase creates a new instance of DoctorUsecase
func NewDoctorAuthUsecase(repo repository.DoctorRepository, jwtService service.JWTService, otpRepo repository.OtpRepository, otpService service.OtpService) DoctorUsecase {
	return &doctorUsecase{
		DoctorRepo: repo,
		JWTService: jwtService,
		OtpRepo:    otpRepo,
		OtpService: otpService,
	}
}

// Register handles doctor registration
func (u *doctorUsecase) Register(doctor *model.Doctor) error {
	if doctor.Email == "" {
		return errors.New("email is required")
	}
	if doctor.Password == "" {
		return errors.New("password is required")
	}
	if doctor.Username == "" {
		return errors.New("username is required")
	}

	// Tetapkan default title_id jika tidak diberikan
	if doctor.TitleID == 0 {
		doctor.TitleID = 1 // ID default dari title
	}

	// Set default role for doctors
	doctor.Role = "doctor"

	// Hash password
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(doctor.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	doctor.Password = string(hashPassword)

	// Save to repository
	err = u.DoctorRepo.CreateDoctor(doctor)
	if err != nil {
		return err
	}

	// Generate OTP
	otpCode := u.OtpService.GenerateOtp()
	expiry := time.Now().Add(5 * time.Minute)

	err = u.OtpRepo.GenerateOtp(doctor.Email, otpCode, expiry)
	if err != nil {
		return err
	}

	// Send OTP via email
	err = helper.SendEmail(doctor.Email, otpCode)
	if err != nil {
		log.Printf("Gagal mengirim email ke %s: %v", doctor.Email, err)
	}

	return nil
}

// Login handles doctor authentication
func (u *doctorUsecase) Login(email, password string) (string, error) {
	doctor, err := u.DoctorRepo.GetByEmail(email)
	if err != nil || doctor == nil {
		return "", errors.New("invalid credentials")
	}

	if !doctor.IsVerified {
		return "", errors.New("account not verified. Please verify your OTP")
	}

	err = bcrypt.CompareHashAndPassword([]byte(doctor.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	token, err := u.JWTService.GenerateJWT(doctor.Email, doctor.ID, doctor.Role, doctor.IsVerified)
	if err != nil {
		return "", err
	}

	return token, nil
}

// VerifyOtp handles OTP verification for doctors
func (u *doctorUsecase) VerifyOtp(email string, code string) error {
	otp, err := u.OtpRepo.GetOtpByEmail(email)
	if err != nil {
		return err
	}

	if otp == nil {
		return errors.New("otp not found")
	}

	if u.OtpService.IsOtpExpired(otp.ExpiresAt) {
		return errors.New("otp expired")
	}

	if otp.Code != code {
		return errors.New("invalid otp")
	}

	// OTP valid, delete OTP
	err = u.OtpRepo.DeleteOtpByEmail(email)
	if err != nil {
		return err
	}

	// Update doctor's verification status
	err = u.DoctorRepo.UpdateDokterVerificationStatus(email, true)
	if err != nil {
		return err
	}

	return nil
}

func (u *doctorUsecase) ResendOtp(email string) error {
	// Validasi apakah email terdaftar
	doctor, err := u.DoctorRepo.GetByEmail(email)
	if err != nil || doctor == nil {
		return errors.New("email not registered")
	}

	// Generate OTP baru
	otpCode := u.OtpService.GenerateOtp()
	expiry := time.Now().Add(5 * time.Minute)

	// Simpan atau perbarui OTP
	err = u.OtpRepo.ResendOtp(email, otpCode, expiry)
	if err != nil {
		return errors.New("failed to resend otp")
	}

	// Kirim OTP melalui email
	err = helper.SendEmail(email, otpCode)
	if err != nil {
		return errors.New("failed to send otp")
	}

	return nil
}
