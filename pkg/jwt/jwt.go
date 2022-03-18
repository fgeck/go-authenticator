package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/floge77/go-authenticator/pkg/models"
	jwtV4 "github.com/golang-jwt/jwt/v4"
)

const (
	DefaultExpireIn = time.Minute * 5
)

type JwtResponse struct {
	Jwt string `json:"jwt"`
}

type JwtClaim struct {
	Role     models.Role          `json:"role"`
	Username string               `json:"username"`
	Claims   jwtV4.StandardClaims `json:"claims"`
}

func (c JwtClaim) Valid() error {
	if !c.Role.IsValid() {
		return errors.New("role is not valid")
	}
	return c.Claims.Valid()
}

type Jwt interface {
	GenerateToken(user *models.User) (string, error)
	ParseAndVerifyToken(signedToken string) (*JwtClaim, error)
}

type jwt struct {
	signingKey string
}

func NewJwt(signingKey string) Jwt {
	return &jwt{signingKey: signingKey}
}

func (j *jwt) GenerateToken(user *models.User) (string, error) {
	expireIn := time.Now().Add(DefaultExpireIn)
	claims := &JwtClaim{Role: user.Role, Username: user.Username, Claims: jwtV4.StandardClaims{ExpiresAt: expireIn.Unix()}}
	token := jwtV4.NewWithClaims(jwtV4.SigningMethodHS256, claims)
	signedJwt, err := token.SignedString([]byte(j.signingKey))
	if err != nil {
		return "", err
	}
	return signedJwt, nil
}

func (j *jwt) ParseAndVerifyToken(signedToken string) (*JwtClaim, error) {
	token, err := jwtV4.ParseWithClaims(signedToken, &JwtClaim{}, func(t *jwtV4.Token) (interface{}, error) {
		return []byte(j.signingKey), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*JwtClaim)
	if !ok {
		return nil, fmt.Errorf("cannot parse claims: %v", token.Claims)
	}

	if claims.Claims.ExpiresAt < time.Now().Local().Unix() {
		return nil, errors.New("token is expired")
	}
	if err = claims.Valid(); err != nil {
		return nil, err
	}
	return claims, nil
}
