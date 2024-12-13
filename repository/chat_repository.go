package repository

import (
	"calmind/model"

	"gorm.io/gorm"
)

type ChatRepository interface {
	SaveChat(chat *model.Chat) error
	GetChatHistoryByUser(userID int) ([]model.Chat, error)
	GetUserByID(userID int) (*model.User, error)
	GetDoctorByID(doctorID int) (*model.Doctor, error)
}
type ChatRepositoryImpl struct {
	DB *gorm.DB
}

func NewChatRepositoryImpl(db *gorm.DB) *ChatRepositoryImpl {
	return &ChatRepositoryImpl{DB: db}
}

func (r *ChatRepositoryImpl) SaveChat(chat *model.Chat) error {
	return r.DB.Create(chat).Error
}

func (r *ChatRepositoryImpl) GetChatHistoryByUser(userID int) ([]model.Chat, error) {
	var chats []model.Chat
	err := r.DB.Preload("User").Preload("Doctor").
		Where("user_id = ?", userID).
		Order("created_at").
		Find(&chats).Error
	return chats, err
}

func (r *ChatRepositoryImpl) GetUserByID(userID int) (*model.User, error) {
	var user model.User
	err := r.DB.First(&user, userID).Error
	return &user, err
}

func (r *ChatRepositoryImpl) GetDoctorByID(doctorID int) (*model.Doctor, error) {
	var doctor model.Doctor
	err := r.DB.First(&doctor, doctorID).Error
	return &doctor, err
}
