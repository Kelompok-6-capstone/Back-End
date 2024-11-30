package model

type User struct {
	ID           int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Username     string `gorm:"not null" json:"username"`
	NoHp         string `gorm:"not null" json:"no_hp"`
	Email        string `gorm:"unique;not null" json:"email"`
	Password     string `gorm:"not null" json:"password"`
	Role         string `gorm:"not null" json:"role"`
	Avatar       string `gorm:"" json:"avatar"`
	Alamat       string `gorm:"" json:"alamat"`
	Tgl_lahir    string `gorm:"" json:"tgl_lahir"`
	JenisKelamin string `gorm:"type:enum('Laki-laki', 'Perempuan');" json:"jenis_kelamin"`
	IsVerified   bool   `json:"is_verified" gorm:"default:false"`
}
