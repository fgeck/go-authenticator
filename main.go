package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/floge77/go-authoriza/pkg/handlers"
)

func main() {
	fmt.Println("Starting go-authoriza")
	err := handlers.HandleAuthrizeRequests()
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(http.ListenAndServe(":9123", nil))
}
