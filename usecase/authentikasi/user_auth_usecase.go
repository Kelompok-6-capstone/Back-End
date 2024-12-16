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

// Memperbaiki fungsi isValidPassword dengan ekspresi reguler yang lebih kompatibel
func isValidPassword(password string) bool {
	// Memastikan password memiliki minimal 8 karakter
	if len(password) < 8 {
		return false
	}
	// Memeriksa apakah ada huruf besar
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	// Memeriksa apakah ada huruf kecil
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	// Memeriksa apakah ada angka
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
	// Memeriksa apakah ada simbol atau karakter khusus
	hasSpecial := regexp.MustCompile(`[\W_]`).MatchString(password)

	// Semua kondisi harus terpenuhi
	return hasUpper && hasLower && hasDigit && hasSpecial
}

func (u *AuthUsecase) Register(user *model.User) error {
	// Validasi input
	if user.Email == "" {
		return errors.New("Email wajib diisi. Mohon masukkan alamat email Anda.")
	}
	if user.Password == "" {
		return errors.New("Password wajib diisi. Mohon masukkan kata sandi Anda.")
	}
	if user.Username == "" {
		return errors.New("Username wajib diisi. Mohon masukkan nama pengguna Anda.")
	}

	// Validasi username
	if !isValidUsername(user.Username) {
		return errors.New("Username minimal 5 karakter dan hanya boleh mengandung huruf, angka, atau garis bawah (_).")
	}

	// Validasi password
	if !isValidPassword(user.Password) {
		return errors.New("Password harus minimal 8 karakter dan mengandung huruf besar, huruf kecil, angka, dan simbol.")
	}

	// Menetapkan peran pengguna
	user.Role = "user"

	// Enkripsi password
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("Gagal mengenkripsi password. Silakan coba lagi.")
	}
	user.Password = string(hashPassword)

	// Simpan user ke database
	err = u.UserRepo.CreateUser(user)
	if err != nil {
		log.Printf("Gagal menyimpan data pengguna: %s", err.Error())
		return errors.New(err.Error())
	}

	// Generate OTP
	otpCode := u.OtpService.GenerateOtp()
	expiry := time.Now().Add(5 * time.Minute)

	// Simpan OTP ke database
	err = u.OtpRepo.GenerateOtp(user.Email, otpCode, expiry)
	if err != nil {
		return errors.New("Gagal membuat kode OTP. Terjadi masalah saat menghasilkan kode OTP.")
	}

	// Kirim email dengan OTP
	err = helper.SendEmail(user.Email, otpCode)
	if err != nil {
		log.Printf("Gagal mengirim email ke %s: %v", user.Email, err)
		return errors.New("Gagal mengirim kode OTP. Pastikan koneksi internet Anda stabil dan coba lagi.")
	}

	log.Println("Proses registrasi selesai.")
	return nil
}

func (u *AuthUsecase) Login(email, password string) (string, error) {
	// Ambil data user berdasarkan email
	user, err := u.UserRepo.GetByUsername(email)
	if err != nil || user == nil {
		log.Printf("Login gagal: Email atau password salah untuk %s", email)
		return "", errors.New("Email atau password salah. Mohon periksa kembali informasi Anda.")
	}

	// Periksa apakah akun sudah terverifikasi
	if !user.IsVerified {
		return "", errors.New("Akun Anda belum terverifikasi. Silakan verifikasi akun dengan kode OTP terlebih dahulu.")
	}

	// Periksa password yang dimasukkan
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Printf("Login gagal: Password tidak cocok untuk %s", email)
		return "", errors.New("Email atau password salah. Mohon periksa kembali informasi Anda.")
	}

	// Generate token JWT
	token, err := u.JWTService.GenerateJWT(user.Email, user.ID, user.Role, user.IsVerified)
	if err != nil {
		log.Printf("Gagal membuat token akses untuk %s: %s", email, err.Error())
		return "", errors.New("Gagal membuat token akses. Terjadi masalah dengan sistem kami. Coba lagi nanti.")
	}

	return token, nil
}

func (a *AuthUsecase) VerifyOtp(email string, code string) error {
	// Ambil data OTP berdasarkan email
	otp, err := a.OtpRepo.GetOtpByEmail(email)
	if err != nil {
		log.Printf("Gagal mengambil data OTP untuk %s: %s", email, err.Error())
		return errors.New("Terjadi kesalahan saat mengambil data OTP. Coba lagi nanti.")
	}

	// Pastikan OTP ditemukan
	if otp == nil {
		return errors.New("Kode OTP tidak ditemukan. Pastikan Anda telah menerima kode OTP.")
	}

	// Validasi apakah OTP sudah kadaluarsa
	if a.OtpService.IsOtpExpired(otp.ExpiresAt) {
		return errors.New("Kode OTP sudah kedaluwarsa. Mohon minta ulang kode OTP.")
	}

	// Cek apakah kode OTP yang dimasukkan valid
	if otp.Code != code {
		return errors.New("Kode OTP yang Anda masukkan tidak valid. Mohon coba lagi.")
	}

	// Hapus OTP setelah berhasil diverifikasi
	err = a.OtpRepo.DeleteOtpByEmail(email)
	if err != nil {
		log.Printf("Gagal menghapus OTP untuk %s: %s", email, err.Error())
		return errors.New("Gagal menghapus kode OTP. Silakan coba lagi.")
	}

	// Update status verifikasi pengguna
	err = a.UserRepo.UpdateUserVerificationStatus(email, true)
	if err != nil {
		log.Printf("Gagal memperbarui status verifikasi untuk %s: %s", email, err.Error())
		return errors.New("Gagal memperbarui status verifikasi akun Anda. Coba lagi nanti.")
	}

	return nil
}

func (u *AuthUsecase) ResendOtp(email string) error {
	// Ambil user berdasarkan email
	user, err := u.UserRepo.GetByUsername(email)
	if err != nil || user == nil {
		log.Printf("Email tidak terdaftar: %s", email)
		return errors.New("Email tidak terdaftar. Pastikan Anda sudah melakukan registrasi.")
	}

	// Generate ulang OTP
	otpCode := u.OtpService.GenerateOtp()
	expiry := time.Now().Add(5 * time.Minute)

	// Perbarui OTP
	err = u.OtpRepo.ResendOtp(email, otpCode, expiry)
	if err != nil {
		log.Printf("Gagal memperbarui kode OTP untuk %s: %s", email, err.Error())
		return errors.New("Gagal memperbarui kode OTP. Coba lagi nanti.")
	}

	// Kirim email dengan OTP yang baru
	err = helper.SendEmail(email, otpCode)
	if err != nil {
		log.Printf("Gagal mengirim OTP ke %s: %v", email, err)
		return errors.New("Gagal mengirim kode OTP. Mohon periksa koneksi internet Anda dan coba lagi.")
	}

	return nil
}
