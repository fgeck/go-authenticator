package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/floge77/go-authoriza/pkg/models"
	"github.com/golang-jwt/jwt/v4"
)

var (
	signingKey      string
	registeredUsers = []models.Credentials{
		{Username: "typ1", Password: "abcd"},
		{Username: "typ2", Password: "efgh"},
	}
)

func Healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "I am alive")
}

func Register(w http.ResponseWriter, r *http.Request) {
	// check here if the user is admin!
	fmt.Fprint(w, "Registering")
}

func SignIn(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
	}
	var creds models.Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !hasValidCredentials(&creds) {
		log.Println(creds.Username)
		log.Println(creds.Password)
		responseUnauthorized(w)
		return
	}
	expireIn := time.Now().Add(time.Minute * 5)
	claims := &models.JwtClaim{Role: models.User, Claim: jwt.StandardClaims{ExpiresAt: expireIn.Unix()}}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedJwt, err := token.SignedString([]byte(signingKey))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, signedJwt)
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

func HandleAuthrizeRequests() error {
	key := os.Getenv("JWT_SIGNING_KEY")
	if key == "" {
		signingKey = "mysupersecret!"
		// return errors.New("no signingKey in env found")
	} else {
		signingKey = signingKey
	}
	http.HandleFunc("/register", Register)
	http.HandleFunc("/signin", SignIn)
	http.HandleFunc("/healthz", Healthz)

	return nil
}

func responseUnauthorized(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	fmt.Fprintf(w, "unauthorized")
}
