package usecase

import (
	"calmind/model"
	"calmind/repository"
	"calmind/service"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type DoctorUsecase interface {
	Register(*model.Doctor) error
	Login(email string, password string) (string, error)
}

type doctorUsecase struct {
	DoctorRepo repository.DoctorRepository
	JWTService service.JWTService
}

// NewDoctorUsecase creates a new instance of DoctorUsecase
func NewDoctorAuthUsecase(repo repository.DoctorRepository, jwtService service.JWTService) DoctorUsecase {
	return &doctorUsecase{
		DoctorRepo: repo,
		JWTService: jwtService,
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
	if doctor.NoHp == "" {
		return errors.New("no handphone is required")
	}
	if doctor.Username == "" {
		return errors.New("username is required")
	}

	// Hash password
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(doctor.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	doctor.Password = string(hashPassword)

	// Set default role for doctors
	doctor.Role = "doctor"

	// Save to repository
	return u.DoctorRepo.CreateDoctor(doctor)
}

// Login handles doctor authentication
func (u *doctorUsecase) Login(email, password string) (string, error) {
	if email == "" {
		return "", errors.New("email is required")
	}

	// Retrieve doctor by email
	doctor, err := u.DoctorRepo.GetByEmail(email)
	if err != nil || doctor == nil {
		return "", errors.New("invalid credentials")
	}

	// Compare password hash
	err = bcrypt.CompareHashAndPassword([]byte(doctor.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// Generate JWT token
	token, err := u.JWTService.GenerateJWT(doctor.Email, doctor.ID, doctor.Role)
	if err != nil {
		return "", err
	}

	return token, nil
}
