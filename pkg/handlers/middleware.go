package handlers

import (
	"log"
	"net/http"

	"github.com/floge77/go-authenticator/pkg/jwt"
	"github.com/floge77/go-authenticator/pkg/models"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

func Chain(f http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, m := range middlewares {
		f = m(f)
	}
	return f
}

func JwtVerify(jwt jwt.Jwt) Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("token")
			if err != nil {
				if err == http.ErrNoCookie {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			signedToken := cookie.Value
			_, err = jwt.ParseAndVerifyToken(signedToken)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}
	}

}

func EnforceAuthenticatedAndRole(role models.Role, jwt jwt.Jwt) Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("token")
			if err != nil {
				if err == http.ErrNoCookie {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			signedToken := cookie.Value
			claim, err := jwt.ParseAndVerifyToken(signedToken)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			if claim.Role != role {
				log.Printf("url %q needs role %q but user %q has only role %q", r.URL.Path, role.String(), claim.Username, claim.Role)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			f(w, r)
		}
	}
}
