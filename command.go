package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kevvarlar/gator/internal/database"
)

type command struct {
	name string
	arguments []string
}

type commands struct {
	commandNames map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	f, ok := c.commandNames[cmd.name]
	if !ok {
		return fmt.Errorf("command does not exist")
	}
	err := f(s, cmd)
	if err != nil {
		return err
	}
	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.commandNames[name] = f
}

func handlerLogin(s *state, cmd command, _ database.User) error {
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("the login handler expects a single argument, the username")
	}
	err := s.gatorConfig.SetUser(cmd.arguments[0])
	if err != nil {
		return err
	}
	fmt.Println("The user has been set")
	return nil
}

func handlerRegister(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("the register handler expects a single argument, the username")
	}
	_, err := s.gatorDB.CreateUser(context.Background(), database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: sql.NullTime{
			Time: time.Now(),
			Valid: true,
		},
		UpdatedAt: sql.NullTime{
			Time: time.Now(),
			Valid: true,
		},
		Name: cmd.arguments[0],
	})
	if err != nil {
		return err
	}
	err = handlerLogin(s, cmd, user)
	if err != nil {
		return err
	}
	fmt.Println("user was created")
	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.gatorDB.Reset(context.Background())
	if err != nil {
		return fmt.Errorf("failed to reset users database: %w", err)
	}
	fmt.Println("Successfully reseted users database")
	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.gatorDB.GetUsers(context.Background())
	if err != nil {
		return err
	}
	for _, user := range users {
		if user.Name == s.gatorConfig.CurrentUserName {
			fmt.Println(" *", user.Name, "(current)")
		} else {
			fmt.Println(" *", user.Name)
		}
	}
	return nil
}

func handlerAgg(s *state, cmd command) error {
	// if len(cmd.arguments) == 0 {
	// 	return fmt.Errorf("the agg handler expects a single argument, the feed url")
	// }
	rss, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}
	fmt.Print(rss)
	return nil
}

func handlerAddfeed(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) < 2 {
		return fmt.Errorf("the add feed handler expects two arguments, the name of the feed and the url of the feed")
	}

	feed, err := s.gatorDB.CreateFeed(context.Background(), database.CreateFeedParams{
		ID: uuid.New(),
		CreatedAt: sql.NullTime{
			Time: time.Now(),
			Valid: true,
		},
		UpdatedAt: sql.NullTime{
			Time: time.Now(),
			Valid: true,
		},
		Name: cmd.arguments[0],
		Url: cmd.arguments[1],
		UserID: uuid.NullUUID{
			UUID: user.ID,
			Valid: true,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create feed: %w", err)
	}
	feed_follow, err := s.gatorDB.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: sql.NullTime{
			Time: time.Now(),
			Valid: true,
		},
		UpdatedAt: sql.NullTime{
			Time: time.Now(),
			Valid: true,
		},
		UserID: uuid.NullUUID{
			UUID: user.ID,
			Valid: true,
		},
		FeedID: uuid.NullUUID{
			UUID: feed.ID,
			Valid: true,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create feed follow: %w", err)
	}
	fmt.Println(feed_follow.FeedName, "was created")
	fmt.Println(feed)
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.gatorDB.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get feeds: %w", err)
	}
	for _, feed := range feeds {
		fmt.Println(feed)
	}
	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("the follow handler expects a single argument, the feed url")
	}

	feed, err := s.gatorDB.GetFeed(context.Background(), cmd.arguments[0])
	if err != nil {
		return fmt.Errorf("failed to get the feed: %w", err)
	}
	feed_follow, err := s.gatorDB.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: sql.NullTime{
			Time: time.Now(),
			Valid: true,
		},
		UpdatedAt: sql.NullTime{
			Time: time.Now(),
			Valid: true,
		},
		UserID: uuid.NullUUID{
			UUID: user.ID,
			Valid: true,
		},
		FeedID: uuid.NullUUID{
			UUID: feed.ID,
			Valid: true,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create feed follow: %w", err)
	}
	fmt.Println("Successfully followed feed!")
	fmt.Println("Feed name:", feed_follow.FeedName, "|", "Feed user:", feed_follow.UserName)
	return nil
}

func handleFollowing(s *state, cmd command) error {
	feed_follows, err := s.gatorDB.GetFeedFollowsForUser(context.Background(), s.gatorConfig.CurrentUserName)
	if err != nil {
		return fmt.Errorf("failed to get feed follows for current user: %w", err)
	}
	fmt.Println("All of the feeds", s.gatorConfig.CurrentUserName, "is following: ")
	for _, follow := range feed_follows {
		fmt.Println(" -", follow.FeedName)
	}
	return nil
}