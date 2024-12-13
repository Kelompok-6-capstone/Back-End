package repository

import (
	"calmind/model"

	"gorm.io/gorm"
)

// ChatRepository defines the interface for chat-related database operations.
type ChatRepository interface {
	GetRoomByUserAndDoctor(userID int, doctorID int) (*model.Chat, error)
	CreateRoom(roomID string, userID int, doctorID int) error
	SaveChat(chat *model.Chat) error
	GetChatsByRoomID(roomID string) ([]model.Chat, error)
}

// ChatRepositoryImpl is the implementation of the ChatRepository interface.
type ChatRepositoryImpl struct {
	DB *gorm.DB
}

// NewChatRepositoryImpl creates a new instance of ChatRepositoryImpl.
func NewChatRepositoryImpl(db *gorm.DB) *ChatRepositoryImpl {
	return &ChatRepositoryImpl{DB: db}
}

// GetRoomByUserAndDoctor retrieves an existing chat room for a specific user and doctor.
// If no room exists, it returns nil and an error.
func (r *ChatRepositoryImpl) GetRoomByUserAndDoctor(userID int, doctorID int) (*model.Chat, error) {
	var chat model.Chat
	err := r.DB.Where("user_id = ? AND doctor_id = ?", userID, doctorID).First(&chat).Error
	if err != nil {
		return nil, err
	}
	return &chat, nil
}

// CreateRoom creates a new chat room with the given RoomID, UserID, and DoctorID.
func (r *ChatRepositoryImpl) CreateRoom(roomID string, userID int, doctorID int) error {
	room := model.Chat{
		RoomID:   roomID,
		UserID:   userID,
		DoctorID: doctorID,
	}
	return r.DB.Create(&room).Error
}

// SaveChat saves a new chat message to the database.
func (r *ChatRepositoryImpl) SaveChat(chat *model.Chat) error {
	return r.DB.Create(chat).Error
}

// GetChatsByRoomID retrieves all chat messages associated with a specific RoomID.
// The messages are ordered by their creation time.
func (r *ChatRepositoryImpl) GetChatsByRoomID(roomID string) ([]model.Chat, error) {
	var chats []model.Chat
	err := r.DB.Where("room_id = ?", roomID).Order("created_at").Find(&chats).Error
	if err != nil {
		return nil, err
	}
	return chats, nil
}
