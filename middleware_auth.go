package main

import (
	"log"
	"net/http"

	"github.com/tirupatihemanth/rssmux/internal/auth"
	"github.com/tirupatihemanth/rssmux/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func middleware_auth(ah authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetApiKey(r.Header)
		if err != nil {
			log.Println("Authorization ApiKey not found:", err)
			respondWithError(w, http.StatusUnauthorized, "Add valid Authorization: ApiKey <Value> to Headers")
			return
		}
		user, err := apiCfg.DB.GetUser(r.Context(), apiKey)

		if err != nil {
			log.Println("Authorization Failed. Unabled to get user: ", err)
			respondWithError(w, http.StatusNotFound, "Unable to get user")
			return
		}
		ah(w, r, user)
	}
}
