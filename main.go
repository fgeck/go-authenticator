package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/floge77/go-authenticator/pkg/config"
	"github.com/floge77/go-authenticator/pkg/db"
	"github.com/floge77/go-authenticator/pkg/handlers"
	"github.com/gorilla/mux"
)

const dbName = "authenticator"

func main() {
	fmt.Println("Starting go-authenticator")
	// for local testing:
	os.Setenv(config.JwtSigningKeyEnvVar, "supersecret")
	os.Setenv(config.DbAddressEnvVar, "postgres")
	os.Setenv(config.DbPortEnvVar, "5432")
	os.Setenv(config.DbUserEnvVar, "pgAdmin")
	os.Setenv(config.DbPasswordEnvVar, "crazyPass123")

	config, err := config.ConfigFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	db, err := db.NewDatabaseConnection(config.DbAddress, config.DbPort, config.DbUser, config.DbPassword, dbName)
	if err != nil {
		log.Fatal(err)
	}
	router := mux.NewRouter()
	registerHandler := handlers.NewRegisterHandler(db)
	signinHandler := handlers.NewSigninHandler(db)

	router.HandleFunc("/healthz", handlers.Healthz)
	router.HandleFunc("/sigin", signinHandler.HandleSignIn).Methods(http.MethodPost)
	router.HandleFunc("/register", registerHandler.RegisterUser).Methods(http.MethodPost)
	router.HandleFunc("/register", registerHandler.GetUsers).Methods(http.MethodGet)
	log.Fatal(http.ListenAndServe(":9123", router))
}
