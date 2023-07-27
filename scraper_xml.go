package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/tirupatihemanth/rssmux/internal/database"
)

const (
	SCRAPING_INTERVAL    = time.Hour * 24
	SCRAPING_CONCURRENCY = 10
)

type RSSFeed struct {
	Channel struct {
		Title string         `xml:"title"`
		Link  string         `xml:"link"`
		Items []RSSFeedItems `xml:"item"`
	} `xml:"channel"`
}

type RSSFeedItems struct {
	Title       string `xml:"title"`
	Url         string `xml:"link"`
	PublishedAt string `xml:"pubDate"`
	Description string `xml:"description"`
}

func scheduleScraping() {
	log.Printf("Starting scraping every :%s with a concurrency of %v\n", SCRAPING_INTERVAL, SCRAPING_CONCURRENCY)
	ticker := time.NewTicker(SCRAPING_INTERVAL)

	for ; ; <-ticker.C {
		dbFeeds, err := apiCfg.DB.GetNextNFeeds(context.Background(), SCRAPING_CONCURRENCY)
		if err != nil {
			log.Println("Couldn't GetNextNFeeds")
			return
		}
		log.Printf("obtained %v Feeds to Scrape\n", len(dbFeeds))
		wg := &sync.WaitGroup{}
		for _, dbFeed := range dbFeeds {
			wg.Add(1)
			go scrapeFeed(dbFeed, wg)
		}
		log.Println("waiting for scraping interval to complete")
	}
}

func scrapeFeed(dbFeed database.Feed, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Println("started scraping Url Id:", dbFeed.Url, dbFeed.ID)
	rssFeed, err := scrapeXMLURL(dbFeed.Url)

	if err != nil {
		log.Println("Couldn't scrape feed:", dbFeed.Url)
		return
	}

	log.Printf("Found %d posts\n", len(rssFeed.Channel.Items))
	for _, feedItem := range rssFeed.Channel.Items {
		description := sql.NullString{}
		if feedItem.Description != "" {
			description.String = feedItem.Description
			description.Valid = true
		}

		pubAt := sql.NullTime{}
		time, err := time.Parse(time.RFC1123Z, feedItem.PublishedAt)
		if err == nil {
			pubAt.Time = time
			pubAt.Valid = true
		}
		_, err = apiCfg.DB.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			Title:       feedItem.Title,
			Url:         feedItem.Url,
			Description: description,
			PublishedAt: pubAt,
			FeedID:      dbFeed.ID,
		})

		if err != nil {
			// If you insert a post with same url again
			if strings.Contains(err.Error(), "duplicate") {
				continue
			}
			log.Println("Couldn't create post:", err)
		}
	}

	log.Println("Scraping done for Url ID:", dbFeed.Url, dbFeed.ID)
}

func scrapeXMLURL(url string) (*RSSFeed, error) {
	httpClient := http.Client{
		Timeout: 30 * time.Second,
	}

	// using http.Get("url") directly uses the Defaultclient.Get()
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var rssFeed RSSFeed
	// Alt: err = xml.NewDecoder(resp.Body).Decode(&rssFeed)
	bytes, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	xml.Unmarshal(bytes, &rssFeed)

	if err != nil {
		return nil, err
	}

	return &rssFeed, nil
}
