package model

type Admin struct {
	ID       int    `json:"id" gorm:"primaryKey"`
	Email    string `gorm:"unique" json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}
