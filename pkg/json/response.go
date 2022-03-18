package json

import (
	"encoding/json"
	"net/http"
)

func WriteJsonResponse(response interface{}, w http.ResponseWriter) error {
	data, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}
	w.Write(data)
	return nil
}

func WriteJsonMessage(msg string, w http.ResponseWriter) error {
	type message struct {
		Message string `json:"message"`
	}
	return WriteJsonResponse(&message{Message: msg}, w)
}
