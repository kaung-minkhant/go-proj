package main

import (
	"context"
	"database/sql"
	"go-proj/6-rss-agg/internal/database"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

func StartScrapping(db *database.Queries, concurrency int, timeBetweenRequest time.Duration) {
	log.Printf("Scrapping on %v goroutine %s duration \n", concurrency, timeBetweenRequest)
	ticker := time.NewTicker(timeBetweenRequest)
	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(context.Background(), int32(concurrency))
		if err != nil {
			log.Println("Error fetching feeds", err)
			continue
		}
		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)
			go scrapeFeed(wg, db, feed)
		}
		wg.Wait()
	}

}

func scrapeFeed(wg *sync.WaitGroup, db *database.Queries, feed database.Feed) {
	defer wg.Done()
	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("Updating feed %v failed %v", feed.ID, err)
		return
	}
	rssFeed, err := URLToFeed(feed.Url)
	if err != nil {
		log.Printf("Scraping %v failed %v", feed.Url, err)
		return
	}
	for _, feedItem := range rssFeed.Channel.Item {
		description := sql.NullString{}
		if feedItem.Description != "" {
			description.String = feedItem.Description
			description.Valid = true
		}
		publishedAt, err := time.Parse(time.RFC1123Z, feedItem.PubDate)
		if err != nil {
			log.Println("Cannot parse publish date of Post", feedItem.Title, "on feed", feed.Name, "with error", err)
			continue
		}
		// log.Println("Post found:", feedItem.Title, "on feed", feed.Name)
		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			Title:       feedItem.Title,
			Description: description,
			PublishedAt: publishedAt,
			Url:         feedItem.Link,
			FeedID:      feed.ID,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint \"posts_url_key\"") {
				continue
			}
			log.Printf("Failed to create post: %v of feed %v with error %v", feedItem.Title, feed.Name, err)
		}
	}
	log.Printf("Feed %v collected, %v posts found. \n", feed.Name, len(rssFeed.Channel.Item))
}
