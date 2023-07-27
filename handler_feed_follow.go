package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/tirupatihemanth/rssmux/internal/database"
)

func feedFollowHandler(w http.ResponseWriter, r *http.Request, user database.User) {

	var feedFollowParams struct {
		FeedId uuid.UUID `json:"feed_id"`
	}

	err := json.NewDecoder(r.Body).Decode(&feedFollowParams)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode feed follow params")
		return
	}

	feedFollow, err := apiCfg.DB.FollowFeed(r.Context(), database.FollowFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feedFollowParams.FeedId,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't follow feed")
		return
	}

	respondWithJSON(w, http.StatusOK, databaseFeedFollowToFeedFollow(feedFollow))
}

func unfollowFeedHandler(w http.ResponseWriter, r *http.Request, user database.User) {
	feedIdStr := chi.URLParam(r, "feedId")
	feedId, err := uuid.Parse(feedIdStr)

	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusBadRequest, "Couldn't get valid feed id from URL params")
		return
	}

	err = apiCfg.DB.UnfollowFeed(r.Context(), database.UnfollowFeedParams{
		UserID: user.ID,
		FeedID: feedId,
	})

	if err != nil {
		log.Println("Couldn't unfollow feed:", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't unfollow feed")
		return
	}
}

func getUserFeedFollowsHandler(w http.ResponseWriter, r *http.Request, user database.User) {
	userFeeds, err := apiCfg.DB.GetUserFeeds(r.Context(), user.ID)

	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Error Getting feeds for the user")
		return
	}

	respondWithJSON(w, http.StatusOK, databaseFeedFollowsToFeedFollows(userFeeds))
}
