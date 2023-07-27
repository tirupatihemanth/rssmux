package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/tirupatihemanth/rssmux/internal/database"
)

func createFeedHandler(w http.ResponseWriter, r *http.Request, user database.User) {
	var feedParams struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}

	err := json.NewDecoder(r.Body).Decode(&feedParams)
	if err != nil {
		log.Println("Couldn't decode create feed request body json:", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode create feed request body json")
		return
	}

	feed, err := apiCfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      feedParams.Name,
		Url:       feedParams.URL,
		UserID:    user.ID,
	})

	if err != nil {
		log.Println("Couldn't create feed for the user:", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't create feed for the user")
		return
	}

	feedFollow, err := apiCfg.DB.FollowFeed(r.Context(), database.FollowFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})

	if err != nil {
		log.Println("Unable to follow user's own feed:", err)
		respondWithError(w, http.StatusInternalServerError, "Unable to follow owner's own feed")
		return
	}

	respondWithJSON(w, http.StatusOK, struct {
		Feed       Feed       `json:"feed"`
		FeedFollow FeedFollow `json:"feed_follow"`
	}{
		Feed:       databaseFeedToFeed(feed),
		FeedFollow: databaseFeedFollowToFeedFollow(feedFollow),
	})
}

func getAllFeedsHandler(w http.ResponseWriter, r *http.Request) {
	feeds, err := apiCfg.DB.GetAllFeeds(r.Context())
	if err != nil {
		log.Println("error GetAllFeeds:", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't get feeds")
		return
	}

	respondWithJSON(w, http.StatusOK, databaseFeedsToFeeds(feeds))
}
