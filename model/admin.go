package model

type Admin struct {
	ID       int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Username string    `json:"username" gorm:"not null"`
	Email    string    `gorm:"unique;not null" json:"email"`
	Password string    `gorm:"not null" json:"password"`
	Role     string    `gorm:"not null" json:"role"`
	Artikel  []Artikel `json:"artikel" gorm:"foreignKey:AdminID"`
}
