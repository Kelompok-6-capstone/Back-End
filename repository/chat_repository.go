package repository

import (
	"calmind/model"

	"gorm.io/gorm"
)

type ChatRepository interface {
	SaveMessage(chat *model.ChatMessage) error
	GetMessages(userID, doctorID int) ([]model.ChatMessage, error)
}

type ChatRepositoryImpl struct {
	DB *gorm.DB
}

func (r *ChatRepositoryImpl) SaveMessage(chat *model.ChatMessage) error {
	return r.DB.Create(chat).Error
}

func (r *ChatRepositoryImpl) GetMessages(userID, doctorID int) ([]model.ChatMessage, error) {
	var messages []model.ChatMessage
	err := r.DB.Where("user_id = ? AND doctor_id = ?", userID, doctorID).
		Order("created_at ASC").
		Find(&messages).Error
	return messages, err
}
