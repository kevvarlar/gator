package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kevvarlar/gator/internal/database"
)

func scrapeFeeds(s* state) error {
	nextFeed, err := s.gatorDB.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get next feed: %w", err)
	}
	fetchedFeed, err := s.gatorDB.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
		LastFetchedAt: sql.NullTime{
			Time: time.Now(),
			Valid: true,
		},
		ID: nextFeed.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to update feed as fetched: %w", err)
	}
	feed, err := fetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		return fmt.Errorf("failed to fetch created feed: %w", err)
	}
	for _, post := range feed.Channel.Item {
		layout := "Mon, 02 Jan 2006 15:04:05 MST"
		pubDate, err := time.Parse(layout, post.PubDate)
		if err != nil {
			return err
		}
		_, err = s.gatorDB.CreatePost(context.Background(), database.CreatePostParams{
			ID: uuid.New(),
			CreatedAt: sql.NullTime{
				Time: time.Now(),
				Valid: true,
			},
			UpdatedAt: sql.NullTime{
				Time: time.Now(),
				Valid: true,
			},
			Title: post.Title,
			Url: post.Link,
			Description: post.Description,
			PublishedAt: pubDate,
			FeedID: fetchedFeed.ID,
		})
		if err != nil {
			fmt.Println(err)
		}
	}
	return nil
}