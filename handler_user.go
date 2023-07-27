package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/tirupatihemanth/rssmux/internal/database"
)

func getUserHandler(w http.ResponseWriter, r *http.Request, user database.User) {
	respondWithJSON(w, http.StatusOK, databaseUserToUser(user))
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	var userParams struct {
		Name string `json:"name"`
	}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&userParams)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode request body json")
		return
	}

	newUser, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      userParams.Name,
	})

	if err != nil {
		log.Println("Unable to create user:", err)
		respondWithError(w, http.StatusInternalServerError, "Unable to create user")
		return
	}

	respondWithJSON(w, http.StatusOK, newUser)
}
