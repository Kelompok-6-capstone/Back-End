package repository

import (
	"calmind/model"

	"gorm.io/gorm"
)

type PaymentRepository interface {
	CreatePayment(payment *model.Payment) error
	VerifyPayment(paymentID int) error
}

type PaymentRepositoryImpl struct {
	DB *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &PaymentRepositoryImpl{DB: db}
}

func (r *PaymentRepositoryImpl) CreatePayment(payment *model.Payment) error {
	return r.DB.Create(payment).Error
}

func (r *PaymentRepositoryImpl) VerifyPayment(paymentID int) error {
	return r.DB.Model(&model.Payment{}).Where("id = ?", paymentID).Update("status", "paid").Error
}
