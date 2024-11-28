package usecase

import (
	"calmind/repository"
	"calmind/service"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type AdminAuthUsecase interface {
	LoginAdmin(email, password string) (string, error)
}

type AdminAuthUsecaseImpl struct {
	AdminRepo  repository.AdminAuthRepository
	JWTService service.JWTService
}

func NewAdminAuthUsecase(repo repository.AdminAuthRepository, jwt service.JWTService) AdminAuthUsecase {
	return &AdminAuthUsecaseImpl{AdminRepo: repo, JWTService: jwt}
}

func (u *AdminAuthUsecaseImpl) LoginAdmin(email, password string) (string, error) {
	if email == "" {
		return "", errors.New("email is required")
	}

	user, err := u.AdminRepo.GetByEmail(email)
	if err != nil || user == nil {
		return "", errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	token, err := u.JWTService.GenerateJWT(user.Email, user.ID, user.Role)
	if err != nil {
		return "", err
	}

	return token, nil
}
