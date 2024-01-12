package model

type User struct {
	Iduser   int    `json:"iduser" gorm:"primary_key;index"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     int    `json:"role"` // 1 = Admin 2 = Petugas
	BaseModel
}