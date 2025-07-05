package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/kevvarlar/gator/internal/database"
)

func scrapeFeeds(s* state) error {
	nextFeed, err := s.gatorDB.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get next feed: %w", err)
	}
	if _, err := s.gatorDB.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
		LastFetchedAt: sql.NullTime{
			Time: time.Now(),
			Valid: true,
		},
		ID: nextFeed.ID,
	}); err != nil {
		return fmt.Errorf("failed to update feed as fetched: %w", err)
	}
	feed, err := fetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		return fmt.Errorf("failed to fetch created feed: %w", err)
	}
	for _, item := range feed.Channel.Item {
		fmt.Println(" -", item.Title)
	}
	return nil
}