package model

// Model untuk Doctor
type Doctor struct {
	ID       int    `json:"id" gorm:"primaryKey"`
	Usename  string `json:"username"`
	No_Hp    string `json:"no_hp"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}
