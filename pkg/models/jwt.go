package models

import (
	"errors"

	"github.com/golang-jwt/jwt/v4"
)

type Role int

func (r Role) String() string {
	return [...]string{"admin", "user"}[r]
}
func (r Role) IsValid() bool {
	return r == Admin || r == User
}

const (
	Admin Role = iota
	User  Role = iota
)

type JwtClaim struct {
	Role  Role `json:"role"`
	Claim jwt.StandardClaims
}

func (c JwtClaim) Valid() error {
	if !c.Role.IsValid() {
		return errors.New("role is not valid")
	}
	return c.Claim.Valid()
}
