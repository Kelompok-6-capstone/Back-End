package usecase

import (
	"calmind/helper"
	"calmind/model"
	repository "calmind/repository/authentikasi"
	"calmind/service"
	"errors"
	"log"
	"regexp"
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

func isValidUsername(username string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9_]{5,}$`)
	return re.MatchString(username)
}

func isValidPassword(password string) bool {
	re := regexp.MustCompile(`^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[\W_]).{8,}$`)
	return re.MatchString(password)
}

func (u *AuthUsecase) Register(user *model.User) error {
	if user.Email == "" {
		return errors.New("email wajib diisi")
	}
	if user.Password == "" {
		return errors.New("password wajib diisi")
	}
	if user.Username == "" {
		return errors.New("username wajib diisi")
	}

	if !isValidUsername(user.Username) {
		return errors.New("username minimal 5 karakter dan hanya boleh mengandung huruf, angka, atau garis bawah (_)")
	}

	if !isValidPassword(user.Password) {
		return errors.New("password minimal 8 karakter dan harus mengandung huruf besar, huruf kecil, angka, dan simbol")
	}

	user.Role = "user"

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("gagal mengenkripsi password")
	}
	user.Password = string(hashPassword)

	err = u.UserRepo.CreateUser(user)
	if err != nil {
		return errors.New("gagal menyimpan data pengguna")
	}

	otpCode := u.OtpService.GenerateOtp()
	expiry := time.Now().Add(5 * time.Minute)

	err = u.OtpRepo.GenerateOtp(user.Email, otpCode, expiry)
	if err != nil {
		return errors.New("gagal membuat kode OTP")
	}

	err = helper.SendEmail(user.Email, otpCode)
	if err != nil {
		log.Printf("Gagal mengirim email ke %s: %v", user.Email, err)
	}

	log.Println("Proses registrasi selesai")
	return nil
}

func (u *AuthUsecase) Login(email, password string) (string, error) {
	user, err := u.UserRepo.GetByUsername(email)
	if err != nil || user == nil {
		return "", errors.New("email atau password salah")
	}

	if !user.IsVerified {
		return "", errors.New("akun belum terverifikasi. Silakan verifikasi OTP terlebih dahulu")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("email atau password salah")
	}

	token, err := u.JWTService.GenerateJWT(user.Email, user.ID, user.Role, user.IsVerified)
	if err != nil {
		return "", errors.New("gagal membuat token akses")
	}

	return token, nil
}

func (a *AuthUsecase) VerifyOtp(email string, code string) error {
	otp, err := a.OtpRepo.GetOtpByEmail(email)
	if err != nil {
		return errors.New("gagal mengambil data OTP")
	}

	if otp == nil {
		return errors.New("kode OTP tidak ditemukan")
	}

	if a.OtpService.IsOtpExpired(otp.ExpiresAt) {
		return errors.New("kode OTP sudah kedaluwarsa")
	}

	if otp.Code != code {
		return errors.New("kode OTP tidak valid")
	}

	err = a.OtpRepo.DeleteOtpByEmail(email)
	if err != nil {
		return errors.New("gagal menghapus kode OTP")
	}

	err = a.UserRepo.UpdateUserVerificationStatus(email, true)
	if err != nil {
		return errors.New("gagal memperbarui status verifikasi pengguna")
	}

	return nil
}

func (u *AuthUsecase) ResendOtp(email string) error {
	user, err := u.UserRepo.GetByUsername(email)
	if err != nil || user == nil {
		return errors.New("email tidak terdaftar")
	}

	otpCode := u.OtpService.GenerateOtp()
	expiry := time.Now().Add(5 * time.Minute)

	err = u.OtpRepo.ResendOtp(email, otpCode, expiry)
	if err != nil {
		return errors.New("gagal memperbarui kode OTP")
	}

	err = helper.SendEmail(email, otpCode)
	if err != nil {
		return errors.New("gagal mengirim kode OTP")
	}

	return nil
}
