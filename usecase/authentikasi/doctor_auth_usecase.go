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

func NewDoctorAuthUsecase(repo repository.DoctorRepository, jwtService service.JWTService, otpRepo repository.OtpRepository, otpService service.OtpService) DoctorUsecase {
	return &doctorUsecase{
		DoctorRepo: repo,
		JWTService: jwtService,
		OtpRepo:    otpRepo,
		OtpService: otpService,
	}
}

func (u *doctorUsecase) Register(doctor *model.Doctor) error {
	if doctor.Email == "" {
		return errors.New("Email wajib diisi. Mohon masukkan alamat email Anda.")
	}
	if doctor.Password == "" {
		return errors.New("Password wajib diisi. Mohon masukkan kata sandi Anda.")
	}
	if doctor.Username == "" {
		return errors.New("Username wajib diisi. Mohon masukkan nama pengguna Anda.")
	}

	if !helper.IsValidUsername(doctor.Username) {
		return errors.New("Username minimal 5 karakter dan hanya boleh mengandung huruf, angka, atau garis bawah (_).")
	}

	if !helper.IsValidPassword(doctor.Password) {
		return errors.New("Password harus minimal 8 karakter dan mengandung huruf besar, huruf kecil, angka, dan simbol.")
	}

	if doctor.TitleID == 0 {
		doctor.TitleID = 1
	}

	doctor.Role = "doctor"

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(doctor.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("Gagal mengenkripsi password. Silakan coba lagi.")
	}
	doctor.Password = string(hashPassword)

	err = u.DoctorRepo.CreateDoctor(doctor)
	if err != nil {
		log.Printf("Gagal menyimpan data pengguna: %s", err.Error())
		return errors.New(err.Error())
	}

	otpCode := u.OtpService.GenerateOtp()
	expiry := time.Now().Add(5 * time.Minute)

	err = u.OtpRepo.GenerateOtp(doctor.Email, otpCode, expiry)
	if err != nil {
		return errors.New("Gagal membuat kode OTP. Terjadi masalah saat menghasilkan kode OTP.")
	}

	err = helper.SendEmail(doctor.Email, otpCode)
	if err != nil {
		log.Printf("Gagal mengirim email ke %s: %v", doctor.Email, err)
		return errors.New("Gagal mengirim kode OTP. Pastikan koneksi internet Anda stabil dan coba lagi.")
	}

	log.Println("Proses registrasi selesai.")
	return nil
}

func (u *doctorUsecase) Login(email, password string) (string, error) {
	doctor, err := u.DoctorRepo.GetByEmail(email)
	if err != nil || doctor == nil {
		log.Printf("Login gagal: Email atau password salah untuk %s", email)
		return "", errors.New("Email atau password salah. Mohon periksa kembali informasi Anda.")
	}

	if !doctor.IsVerified {
		return "", errors.New("Akun Anda belum terverifikasi. Silakan verifikasi akun dengan kode OTP terlebih dahulu.")
	}

	err = bcrypt.CompareHashAndPassword([]byte(doctor.Password), []byte(password))
	if err != nil {
		log.Printf("Login gagal: Password tidak cocok untuk %s", email)
		return "", errors.New("Email atau password salah. Mohon periksa kembali informasi Anda.")
	}

	token, err := u.JWTService.GenerateJWT(doctor.Email, doctor.ID, doctor.Role, doctor.IsVerified)
	if err != nil {
		log.Printf("Gagal membuat token akses untuk %s: %s", email, err.Error())
		return "", errors.New("Gagal membuat token akses. Terjadi masalah dengan sistem kami. Coba lagi nanti.")
	}

	return token, nil
}

func (u *doctorUsecase) VerifyOtp(email string, code string) error {
	otp, err := u.OtpRepo.GetOtpByEmail(email)
	if err != nil {
		log.Printf("Gagal mengambil data OTP untuk %s: %s", email, err.Error())
		return errors.New("Terjadi kesalahan saat mengambil data OTP. Coba lagi nanti.")
	}

	if otp == nil {
		return errors.New("Kode OTP tidak ditemukan. Pastikan Anda telah menerima kode OTP.")
	}

	if u.OtpService.IsOtpExpired(otp.ExpiresAt) {
		return errors.New("Kode OTP sudah kedaluwarsa. Mohon minta ulang kode OTP.")
	}

	if otp.Code != code {
		return errors.New("Kode OTP yang Anda masukkan tidak valid. Mohon coba lagi.")
	}

	err = u.OtpRepo.DeleteOtpByEmail(email)
	if err != nil {
		log.Printf("Gagal menghapus OTP untuk %s: %s", email, err.Error())
		return errors.New("Gagal menghapus kode OTP. Silakan coba lagi.")
	}

	err = u.DoctorRepo.UpdateDokterVerificationStatus(email, true)
	if err != nil {
		log.Printf("Gagal memperbarui status verifikasi untuk %s: %s", email, err.Error())
		return errors.New("Gagal memperbarui status verifikasi akun Anda. Coba lagi nanti.")
	}

	return nil
}

func (u *doctorUsecase) ResendOtp(email string) error {
	doctor, err := u.DoctorRepo.GetByEmail(email)
	if err != nil || doctor == nil {
		log.Printf("Email tidak terdaftar: %s", email)
		return errors.New("Email tidak terdaftar. Pastikan Anda sudah melakukan registrasi.")
	}

	otpCode := u.OtpService.GenerateOtp()
	expiry := time.Now().Add(5 * time.Minute)

	err = u.OtpRepo.ResendOtp(email, otpCode, expiry)
	if err != nil {
		log.Printf("Gagal memperbarui kode OTP untuk %s: %s", email, err.Error())
		return errors.New("Gagal memperbarui kode OTP. Coba lagi nanti.")
	}

	err = helper.SendEmail(email, otpCode)
	if err != nil {
		log.Printf("Gagal mengirim OTP ke %s: %v", email, err)
		return errors.New("Gagal mengirim kode OTP. Mohon periksa koneksi internet Anda dan coba lagi.")
	}

	return nil
}
