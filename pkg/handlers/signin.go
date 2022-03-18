package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/floge77/go-authenticator/pkg/db"
	jsonMessage "github.com/floge77/go-authenticator/pkg/json"
	"github.com/floge77/go-authenticator/pkg/models"
	"github.com/golang-jwt/jwt/v4"
)

var (
	SigningKey      string
	registeredUsers = []models.Credentials{
		{Username: "typ1", Password: "abcd"},
		{Username: "typ2", Password: "efgh"},
	}
)

type SigninHandler interface {
	HandleSignIn(w http.ResponseWriter, r *http.Request)
}

type signinHandler struct {
	db db.DatabaseConnection
}

func NewSigninHandler(db db.DatabaseConnection) SigninHandler {
	return &signinHandler{db: db}
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
	expireIn := time.Now().Add(time.Minute * 5)
	claims := &models.JwtClaim{Role: models.User, Claim: jwt.StandardClaims{ExpiresAt: expireIn.Unix()}}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedJwt, err := token.SignedString([]byte(SigningKey))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	jwtResponse := &models.JwtResponse{Jwt: signedJwt}

	if err = jsonMessage.WriteJsonResponse(jwtResponse, w); err != nil {
		log.Println(err)
	}
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
