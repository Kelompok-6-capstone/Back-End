package repository

import (
	"calmind/model"

	"gorm.io/gorm"
)

type ChatLogRepository interface {
	SaveLog(chatLog *model.ChatLog) error
	GetLogsByUserID(userID int) ([]model.ChatLog, error)
}

type ChatLogRepositoryImpl struct {
	DB *gorm.DB
}

func NewChatLogRepository(db *gorm.DB) ChatLogRepository {
	return &ChatLogRepositoryImpl{DB: db}
}

func (r *ChatLogRepositoryImpl) SaveLog(chatLog *model.ChatLog) error {
	return r.DB.Create(chatLog).Error
}

func (r *ChatLogRepositoryImpl) GetLogsByUserID(userID int) ([]model.ChatLog, error) {
	var logs []model.ChatLog
	err := r.DB.Where("user_id = ?", userID).Order("created_at ASC").Find(&logs).Error
	return logs, err
}
