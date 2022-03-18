package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/floge77/go-authenticator/pkg/db"
	"github.com/floge77/go-authenticator/pkg/jwt"
	"github.com/floge77/go-authenticator/pkg/models"
)

var (
	SigningKey      string
	registeredUsers = []models.Credentials{
		{Username: "typ1", Password: "abcd", Role: "admin"},
		{Username: "typ2", Password: "efgh", Role: "admin"},
	}
)

type SigninHandler interface {
	HandleSignIn(w http.ResponseWriter, r *http.Request)
}

type signinHandler struct {
	db  db.Database
	jwt jwt.Jwt
}

func NewSigninHandler(db db.Database, jwt jwt.Jwt) SigninHandler {
	return &signinHandler{db: db, jwt: jwt}
}

func Healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("I am alive!")
}

func (sh *signinHandler) HandleSignIn(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var creds models.Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !hasValidCredentials(&creds) {
		responseUnauthorized(w)
		return
	}
	signedToken, err := sh.jwt.GenerateToken(creds)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   signedToken,
		Expires: time.Now().Add(jwt.DefaultExpireIn),
	})
	json.NewEncoder(w).Encode(signedToken)
}

func hasValidCredentials(creds *models.Credentials) bool {
	for _, user := range registeredUsers {
		if user.Username == creds.Username {
			if user.Password == creds.Password {
				return true
			}
		}
	}
	return false
}

func responseUnauthorized(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode("unauthorized")
}
