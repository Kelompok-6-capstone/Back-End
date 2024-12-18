// repository/chatbot_ai_doctor/chat_log_repository.go
package repository

import (
	"calmind/model"
	"gorm.io/gorm"
)

type DoctorChatbotRepository interface {
	GetLogsByDoctorID(doctorID int) ([]model.Chatbot, error)
	SaveLog(log *model.Chatbot) error
}

type DoctorChatbotRepositoryImpl struct {
	DB *gorm.DB
}

func NewDoctorChatbotRepository(db *gorm.DB) DoctorChatbotRepository {
	return &DoctorChatbotRepositoryImpl{DB: db}
}

// GetLogsByDoctorID mengambil semua log percakapan untuk dokter tertentu berdasarkan ID.
func (repo *DoctorChatbotRepositoryImpl) GetLogsByDoctorID(doctorID int) ([]model.Chatbot, error) {
	var logs []model.Chatbot
	err := repo.DB.Where("user_id = ?", doctorID).Find(&logs).Error
	if err != nil {
		return nil, err
	}
	return logs, nil
}

// SaveLog menyimpan log percakapan baru ke dalam database.
func (repo *DoctorChatbotRepositoryImpl) SaveLog(log *model.Chatbot) error {
	return repo.DB.Create(log).Error
}
