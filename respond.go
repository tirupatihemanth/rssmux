package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, errMsg string) {
	if code >= 500 {
		log.Println("Responding with 5xx error:", errMsg)
	}

	respondWithJSON(w, code, struct {
		Error string `json:"error"`
	}{Error: errMsg})
}

func respondWithJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")

	data, err := json.Marshal(payload)

	if err != nil {
		log.Println("Error marshalling JSON:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(code)
	w.Write(data)
}
