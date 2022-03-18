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
	registeredUsers = []models.User{
		{Username: "typ1", Password: "abcd", Role: models.Admin},
		{Username: "typ2", Password: "efgh"},
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
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	foundCredentials, err := sh.FindUser(&user)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if foundCredentials.Password != user.Password {
		responseUnauthorized(w)
		return
	}
	signedToken, err := sh.jwt.GenerateToken(&user)
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

func (sh *signinHandler) FindUser(creds *models.User) (*models.User, error) {
	return sh.db.UserByName(creds.Username)
}

func responseUnauthorized(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode("unauthorized")
}
