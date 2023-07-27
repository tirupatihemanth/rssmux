package main

import (
	"net/http"
	"strconv"

	"github.com/tirupatihemanth/rssmux/internal/database"
)

func getPostsForUserHandler(w http.ResponseWriter, r *http.Request, user database.User) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))

	if err != nil {
		limit = 10
	}

	dbPosts, err := apiCfg.DB.GetPostsForUser(r.Context(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to get posts for the user")
		return
	}
	respondWithJSON(w, http.StatusOK, databasePostsToPosts(dbPosts))
}
