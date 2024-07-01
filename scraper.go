package main

import (
	"context"
	"github.com/unrealandychan/rssagg/internal/database"
	"log"
	"sync"
	"time"
)

func startScraping(
	db *database.Queries,
	concurrency int,
	timeBetweenFetches time.Duration,
) {
	log.Printf("Starting scraping with %d workers\n, every %s second", concurrency, timeBetweenFetches)

	ticker := time.NewTicker(timeBetweenFetches)
	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(context.Background(), int32(concurrency))
		if err != nil {
			log.Printf("Error getting feeds to fetch: %v", err)
			continue
		}
		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)
			go scrapFeed(db, wg, feed)
		}
		wg.Wait()
	}
}

func scrapFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()

	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("Error marking feed as fetching: %v", err)
	}

	urlToFeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Printf("Error fetching feed: %v", err)
		return
	}
	for _, item := range urlToFeed.Channel.Item {
		log.Println("Found post: ", item.Title, "on feed", feed.Name)
	}
	log.Printf("Feed %s collected, %v posts found", feed.Name, len(urlToFeed.Channel.Item))
}
