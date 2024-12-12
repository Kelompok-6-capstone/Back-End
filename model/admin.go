package model

type Admin struct {
	ID       int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Income   float64   `gorm:"default:0"`
	Username string    `gorm:"not null" json:"username"`
	Email    string    `gorm:"unique;not null" json:"email"`
	Password string    `gorm:"not null" json:"password"`
	Role     string    `gorm:"not null" json:"role"`
	Artikels []Artikel `gorm:"foreignKey:AdminID" json:"artikels"`
}