package usecase

import (
	"calmind/model"
	"calmind/repository"
)

type PaymentUsecase interface {
	CreatePayment(consultationID int, amount int) error
	VerifyPayment(paymentID int) error
}

type PaymentUsecaseImpl struct {
	PaymentRepo repository.PaymentRepository
}

func NewPaymentUsecase(repo repository.PaymentRepository) PaymentUsecase {
	return &PaymentUsecaseImpl{PaymentRepo: repo}
}

func (u *PaymentUsecaseImpl) CreatePayment(consultationID int, amount int) error {
	payment := model.Payment{
		ConsultationID: consultationID,
		Amount:         amount,
		Status:         "pending",
	}
	return u.PaymentRepo.CreatePayment(&payment)
}

func (u *PaymentUsecaseImpl) VerifyPayment(paymentID int) error {
	return u.PaymentRepo.VerifyPayment(paymentID)
}
