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

type UserUsecase interface {
	Register(*model.User) error
	Login(email string, password string) (string, error)
	VerifyOtp(email string, code string) error
	ResendOtp(email string) error
}

type AuthUsecase struct {
	UserRepo   repository.UserRepository
	JWTService service.JWTService
	OtpRepo    repository.OtpRepository
	OtpService service.OtpService
}

func NewAuthUsecase(repo repository.UserRepository, jwtService service.JWTService, otpRepo repository.OtpRepository, otpService service.OtpService) UserUsecase {
	return &AuthUsecase{UserRepo: repo, JWTService: jwtService, OtpRepo: otpRepo, OtpService: otpService}
}

func (u *AuthUsecase) Register(user *model.User) error {
	if user.Email == "" {
		return errors.New("email is required")
	}
	if user.Password == "" {
		return errors.New("password is required")
	}
	if user.Username == "" {
		return errors.New("username is required")
	}

	user.Role = "user"

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashPassword)

	u.UserRepo.CreateUser(user)
	// Generate OTP
	otpCode := u.OtpService.GenerateOtp()
	expiry := time.Now().Add(5 * time.Minute)

	err = u.OtpRepo.GenerateOtp(user.Email, otpCode, expiry)
	if err != nil {
		return err
	}

	err = helper.SendEmail(user.Email, otpCode)
	if err != nil {
		log.Printf("Gagal mengirim email ke %s: %v", user.Email, err)
	}
	log.Println("Proses selesai")

	return nil
}

func (u *AuthUsecase) Login(email, password string) (string, error) {
	user, err := u.UserRepo.GetByUsername(email)
	if err != nil || user == nil {
		return "", errors.New("invalid credentials")
	}

	if !user.IsVerified {
		return "", errors.New("account not verified. Please verify your OTP")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	token, err := u.JWTService.GenerateJWT(user.Email, user.ID, user.Role, user.IsVerified)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (a *AuthUsecase) VerifyOtp(email string, code string) error {
	otp, err := a.OtpRepo.GetOtpByEmail(email)
	if err != nil {
		return err
	}

	if otp == nil {
		return errors.New("otp not found")
	}

	if a.OtpService.IsOtpExpired(otp.ExpiresAt) {
		return errors.New("otp expired")
	}

	if otp.Code != code {
		return errors.New("invalid otp")
	}

	// OTP valid, hapus OTP
	err = a.OtpRepo.DeleteOtpByEmail(email)
	if err != nil {
		return err
	}

	// Perbarui status pengguna menjadi verified
	err = a.UserRepo.UpdateUserVerificationStatus(email, true)
	if err != nil {
		return err
	}

	return nil
}

func (u *AuthUsecase) ResendOtp(email string) error {
	// Validasi apakah email terdaftar
	doctor, err := u.UserRepo.GetByUsername(email)
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
