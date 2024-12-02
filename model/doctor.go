package model

type Doctor struct {
	ID          int     `json:"id" gorm:"primaryKey;autoIncrement"`
	Username    string  `gorm:"not null" json:"username"`
	NoHp        string  `gorm:"not null" json:"no_hp"`
	Email       string  `gorm:"unique;not null" json:"email"`
	Password    string  `gorm:"not null" json:"password"`
	Role        string  `gorm:"not null" json:"role"`
	Avatar      string  `json:"avatar"`
	DateOfBirth string  `json:"date_of_birth"`
	Address     string  `json:"address"`
	Schedule    string  `json:"schedule"`
	IsVerified  bool    `json:"is_verified" gorm:"default:false"`
	IsActive    bool    `json:"is_active" gorm:"default:true"`
	Title       string  `json:"title"`
	Price       float64 `json:"price"`
	Experience  int     `json:"experience"`
	STRNumber   string  `json:"str_number"`
	About       string  `json:"about"`
}
