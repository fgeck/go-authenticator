package models

import (
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username string `gorm:"primaryKey"`
	Role     string
	Password string
	Updated  int64 `gorm:"autoUpdateTime"`
	Created  int64 `gorm:"autoCreateTime"`
}

func (user *User) ValidaePassword(password string) error {
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
