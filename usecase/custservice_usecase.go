package usecase

import (
	"calmind/model"
	"calmind/repository"
)

type CustServiceUsecase interface {
	SaveCustService(userID int, message string) error
	GetUnansweredMessages() ([]model.CustService, error)
	AnswerMessage(id int, answer string) error
}

type CustServiceUsecaseImpl struct {
	CustServiceRepo repository.CustServiceRepository
}

func NewCustServiceUsecase(repo repository.CustServiceRepository) CustServiceUsecase {
	return &CustServiceUsecaseImpl{CustServiceRepo: repo}
}

func (u *CustServiceUsecaseImpl) SaveCustService(userID int, message string) error {
	custService := &model.CustService{
		UserID:     userID,
		Message:    message,
		IsAnswered: false,
	}
	return u.CustServiceRepo.SaveCustService(custService)
}

func (u *CustServiceUsecaseImpl) GetUnansweredMessages() ([]model.CustService, error) {
	return u.CustServiceRepo.GetUnansweredMessages()
}

func (u *CustServiceUsecaseImpl) AnswerMessage(id int, answer string) error {
	return u.CustServiceRepo.AnswerMessage(id, answer)
}
