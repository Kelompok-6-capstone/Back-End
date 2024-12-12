package repository

import (
	"calmind/model"
	"time"

	"gorm.io/gorm"
)

type CustServiceRepository interface {
	SaveCustService(custService *model.CustService) error
	GetUnansweredMessages() ([]model.CustService, error)
	AnswerMessage(id int, answer string) error
}

type CustServiceRepositoryImpl struct {
	DB *gorm.DB
}

func NewCustServiceRepository(db *gorm.DB) CustServiceRepository {
	return &CustServiceRepositoryImpl{DB: db}
}

func (r *CustServiceRepositoryImpl) SaveCustService(custService *model.CustService) error {
	return r.DB.Create(custService).Error
}

func (r *CustServiceRepositoryImpl) GetUnansweredMessages() ([]model.CustService, error) {
	var messages []model.CustService
	err := r.DB.Where("is_answered = ?", false).Find(&messages).Error
	return messages, err
}

func (r *CustServiceRepositoryImpl) AnswerMessage(id int, answer string) error {
	return r.DB.Model(&model.CustService{}).Where("id = ?", id).Updates(map[string]interface{}{
		"answer":      answer,
		"is_answered": true,
		"answered_at": time.Now(),
	}).Error
}
