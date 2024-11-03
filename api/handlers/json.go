package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

type HttpErrorReponse struct {
	Message string `json:"message"`
}

func RespondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}

	errResp := HttpErrorReponse{
		Message: msg,
	}
	RespondWithJson(w, code, errResp)
}

func RespondWithJson(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("content-type", "application/json")

	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		http.Error(w, "error while marshelling json", 500)
		return
	}
	w.WriteHeader(code)
	w.Write(dat)
}
