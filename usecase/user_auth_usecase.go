package usecase

import (
	"calmind/model"
	"calmind/repository"
	"calmind/service"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type UserUsecase interface {
	Register(*model.User) error
	Login(email string, password string) (string, error)
}

type AuthUsecase struct {
	UserRepo   repository.UserRepository
	JWTService service.JWTService
}

func NewAuthUsecase(repo repository.UserRepository, jwtService service.JWTService) UserUsecase {
	return &AuthUsecase{UserRepo: repo, JWTService: jwtService}
}

func (u *AuthUsecase) Register(user *model.User) error {
	if user.Email == "" {
		return errors.New("email is required")
	}
	if user.Password == "" {
		return errors.New("password is required")
	}
	if user.NoHp == "" {
		return errors.New("no handphone is required")
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

	return u.UserRepo.CreateUser(user)
}

func (u *AuthUsecase) Login(email, password string) (string, error) {
	if email == "" {
		return "", errors.New("email is required")
	}

	user, err := u.UserRepo.GetByUsername(email)
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
