package handlers

import "github.com/floge77/go-authenticator/pkg/jwt"

type JwtMiddleWare interface {

}

type JwtMiddleWare struct {
	jwt jwt.Jwt
}


