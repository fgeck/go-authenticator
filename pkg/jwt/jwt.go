package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/floge77/go-authenticator/pkg/models"
	jwtV4 "github.com/golang-jwt/jwt/v4"
)

type Role int

func (r Role) String() string {
	return [...]string{"admin", "creator", "user"}[r]
}
func (r Role) IsValid() bool {
	return r == Admin || r == User || r == Creator
}

func ToRole(s string) Role {
	switch s {
	case "admin":
		return 0
	case "creator":
		return 1
	case "user":
		return 2
	default:
		return -1
	}
}

const (
	Admin   Role = iota
	Creator Role = iota
	User    Role = iota

	DefaultExpireIn = time.Minute * 5
)

type JwtResponse struct {
	Jwt string `json:"jwt"`
}

type JwtClaim struct {
	Role   Role                 `json:"role"`
	User   string               `json:"user"`
	Claims jwtV4.StandardClaims `json:"claims"`
}

func (c JwtClaim) Valid() error {
	if !c.Role.IsValid() {
		return errors.New("role is not valid")
	}
	return c.Claims.Valid()
}

type Jwt interface {
	GenerateToken(credentials models.Credentials) (string, error)
	ValidateToken(string) error
}

type jwt struct {
	signingKey string
}

func NewJwt(signingKey string) Jwt {
	return &jwt{signingKey: signingKey}
}

func (j *jwt) GenerateToken(credentials models.Credentials) (string, error) {
	expireIn := time.Now().Add(DefaultExpireIn)
	role := ToRole(credentials.Role)
	claims := &JwtClaim{Role: role, Claims: jwtV4.StandardClaims{ExpiresAt: expireIn.Unix()}}
	token := jwtV4.NewWithClaims(jwtV4.SigningMethodHS256, claims)
	signedJwt, err := token.SignedString([]byte(j.signingKey))
	if err != nil {
		return "", err
	}
	return signedJwt, nil
}

func (j *jwt) ValidateToken(signedToken string) error {
	token, err := jwtV4.ParseWithClaims(signedToken, &JwtClaim{}, func(t *jwtV4.Token) (interface{}, error) {
		return []byte(j.signingKey), nil
	})
	if err != nil {
		return err
	}
	claims, ok := token.Claims.(*JwtClaim)
	if !ok {
		return fmt.Errorf("cannot parse claims: %v", token.Claims)
	}

	if claims.Claims.ExpiresAt < time.Now().Local().Unix() {
		return errors.New("token is expired")
	}
	return claims.Valid()
}
