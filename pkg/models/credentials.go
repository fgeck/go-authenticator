package models

type Credentials struct {
	Username string `json:"username" gorm:"primaryKey"`
	Password string `json:"password"`
	Role     string `json:"role"`
}
