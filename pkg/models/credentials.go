package models

type Credentials struct {
	Username string `json:"username" gorm:"primaryKey"`
	Password string `json:"password"`
}
