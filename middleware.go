package main

import (
	"context"
	"fmt"

	"github.com/kevvarlar/gator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.gatorDB.GetUser(context.Background(), s.gatorConfig.CurrentUserName)
		if err != nil {
			return fmt.Errorf("failed to get current user: %w", err)
		}
		err = handler(s, cmd, user)
		if err != nil {
			return err
		}
		return nil
	}
}