package main

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/PraneGIT/rssagg/internal/database"
	"github.com/google/uuid"
)

func startScraping(db *database.Queries, concurrency int, timeBetweenRequests time.Duration) {
	log.Printf("scraping on %v goroutines every %s duration", concurrency, timeBetweenRequests)

	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		log.Println("scraping...")
		feeds, err := db.GetNextFeedsToFetch(context.Background(), int32(concurrency))
		if err != nil {
			log.Println(err)
			continue
		}

		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)

			go scrapFeed(wg, db, feed)
		}
		wg.Wait()
	}
}

func scrapFeed(wg *sync.WaitGroup, db *database.Queries, feed database.Feed) {
	defer wg.Done()

	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Println(err)
		return
	}

	rssFeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Println(err)
		return
	}

	for _, item := range rssFeed.Channel.Item {

		description := sql.NullString{}
		if item.Description != "" {
			description.Valid = true
			description.String = item.Description
		}

		t, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			log.Println(err)
			continue
		}

		_, err = db.CreatePosts(context.Background(), database.CreatePostsParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       item.Title,
			Description: description,
			PublishedAt: t,
			Url:         item.Link,
			FeedID:      feed.ID,
		})

		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				continue
			}
			log.Println(err)
			continue
		}
	}
	log.Printf("feed %s collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))

}
