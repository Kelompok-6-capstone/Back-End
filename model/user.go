package model

type User struct {
	ID       int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Username string `gorm:"not null" json:"username"`
	NoHp     string `gorm:"not null" json:"no_hp"`
	Email    string `gorm:"unique;not null" json:"email"`
	Password string `gorm:"not null" json:"password"`
	Role     string `gorm:"not null" json:"role"`
	Avatar   string `gorm:"" json:"avatar"`
	Bio      string `gorm:"" json:"bio"`
}
