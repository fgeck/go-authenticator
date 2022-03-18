package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/floge77/go-authenticator/pkg/db"
	"github.com/floge77/go-authenticator/pkg/jwt"
	"github.com/floge77/go-authenticator/pkg/models"
)

type RegisterHandler interface {
	RegisterUser(w http.ResponseWriter, r *http.Request)
	GetUsers(w http.ResponseWriter, r *http.Request)
}

type registerHandler struct {
	database db.Database
}

func NewRegisterHandler(connection db.Database, jwt jwt.Jwt) RegisterHandler {
	return &registerHandler{database: connection}
}

func (rh *registerHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var credential models.Credentials
	err = json.Unmarshal(body, &credential)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = rh.database.AddCredential(&credential)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("error")
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode("Created")
}

func (rh *registerHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	credentials, err := rh.database.GetAllCredentials()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(credentials)
}
