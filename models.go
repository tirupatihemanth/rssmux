package main

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/tirupatihemanth/rssmux/internal/database"
)

// Since database/models.go is autogenerated by sqlc and it doesn't have json tags these models help format
// database objects the way we want in the json

// Keeps database and http APIs independent

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	ApiKey    string    `json:"api_key"`
}

type Feed struct {
	ID            uuid.UUID  `json:"id"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	Name          string     `json:"name"`
	Url           string     `json:"url"`
	UserID        uuid.UUID  `json:"user_id"`
	LastFetchedAt *time.Time `json:"last_fetched_at"`
}

type FeedFollow struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uuid.UUID `json:"user_id"`
	FeedID    uuid.UUID `json:"feed_id"`
}

type Post struct {
	ID          uuid.UUID  `json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Title       string     `json:"title"`
	Url         string     `json:"url"`
	Description *string    `json:"description"`
	PublishedAt *time.Time `json:"published_at"`
	FeedID      uuid.UUID  `json:"feed_id"`
}

func databaseUserToUser(dbUser database.User) User {
	return User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Name:      dbUser.Name,
		ApiKey:    dbUser.ApiKey,
	}
}

func databaseFeedToFeed(dbFeed database.Feed) Feed {
	return Feed{
		ID:            dbFeed.ID,
		CreatedAt:     dbFeed.CreatedAt,
		UpdatedAt:     dbFeed.UpdatedAt,
		UserID:        dbFeed.UserID,
		Name:          dbFeed.Name,
		Url:           dbFeed.Url,
		LastFetchedAt: nullTimeToTimePtr(dbFeed.LastFetchedAt),
	}
}

func databaseFeedsToFeeds(dbFeeds []database.Feed) []Feed {
	feeds := make([]Feed, len(dbFeeds))

	for i, dbFeed := range dbFeeds {
		feeds[i] = databaseFeedToFeed(dbFeed)
	}
	return feeds
}

func databaseFeedFollowToFeedFollow(dbFF database.FeedFollow) FeedFollow {
	return FeedFollow{
		ID:        dbFF.ID,
		CreatedAt: dbFF.CreatedAt,
		UpdatedAt: dbFF.UpdatedAt,
		UserID:    dbFF.UserID,
		FeedID:    dbFF.FeedID,
	}
}

func databaseFeedFollowsToFeedFollows(dbFFs []database.FeedFollow) []FeedFollow {
	feedFollows := make([]FeedFollow, len(dbFFs))

	for i, dbFF := range dbFFs {
		feedFollows[i] = databaseFeedFollowToFeedFollow(dbFF)
	}
	return feedFollows
}

func databasePostToPost(dbPost database.Post) Post {
	return Post{
		ID:          dbPost.ID,
		CreatedAt:   dbPost.CreatedAt,
		UpdatedAt:   dbPost.UpdatedAt,
		Title:       dbPost.Title,
		Url:         dbPost.Title,
		Description: nullStrToStrPtr(dbPost.Description),
		PublishedAt: nullTimeToTimePtr(dbPost.PublishedAt),
		FeedID:      dbPost.FeedID,
	}
}

func databasePostsToPosts(dbPosts []database.Post) []Post {
	posts := make([]Post, len(dbPosts))

	for i, dbPost := range dbPosts {
		posts[i] = databasePostToPost(dbPost)
	}
	return posts
}

func nullTimeToTimePtr(nullTime sql.NullTime) *time.Time {
	if nullTime.Valid {
		return &nullTime.Time
	}

	return nil
}

func nullStrToStrPtr(nullString sql.NullString) *string {
	if nullString.Valid {
		return &nullString.String
	}

	return nil
}
