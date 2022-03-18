package models

import (
	"golang.org/x/crypto/bcrypt"
)

const (
	Admin  Role = iota
	Viewer Role = iota
)

type Role int

func (r Role) String() string {
	switch r {
	case 0:
		return "admin"
	case 1:
		return "viewer"
	default:
		return "undefined"
	}
}
func (r Role) IsValid() bool {
	return r == Admin || r == Viewer
}

func ToRole(s string) Role {
	switch s {
	case "admin":
		return 0
	case "viewer":
		return 1
	default:
		return -1
	}
}

type User struct {
	Username string `json:"username" gorm:"primaryKey"`
	Password string `json:"password"`
	Role     Role   `json:"role"`
	Updated  int64  `gorm:"autoUpdateTime"`
	Created  int64  `gorm:"autoCreateTime"`
}

func (user *User) ValidatePassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return err
	}

	return nil
}

func (u *User) HashPassword() error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(u.Password), 14)
	if err != nil {
		return err
	}
	u.Password = string(bytes)

	return nil

}
