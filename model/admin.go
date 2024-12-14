package model

type Admin struct {
	ID       int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Username string    `gorm:"not null" json:"username"`
	Avatar   string    `gorm:"not null" json:"avatar"`
	Email    string    `gorm:"unique;not null" json:"email"`
	Password string    `gorm:"not null" json:"password"`
	Role     string    `gorm:"not null" json:"role"`
	Artikels []Artikel `gorm:"foreignKey:AdminID" json:"artikels"`
}
