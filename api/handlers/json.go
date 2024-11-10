package handlers

import (
	"encoding/json"
	"net/http"
)

type HttpErrorResponse struct {
	Message string `json:"message"`
}

func RespondWithError(w http.ResponseWriter, code int, msg string) {
	errResp := HttpErrorResponse{
		Message: msg,
	}
	RespondWithJson(w, code, errResp)
}

func RespondWithJson(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("content-type", "application/json")

	dat, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "error while marshelling json", 500)
		return
	}
	w.WriteHeader(code)
	w.Write(dat)
}
